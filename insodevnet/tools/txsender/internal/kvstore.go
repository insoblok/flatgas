package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Action string

type Bucket string

const (
	ActionCreate     Action = "create"
	ActionDelete     Action = "delete"
	ActionUpdate     Action = "update"
	ActionRollback   Action = "rollback"
	ActionUpdateMeta Action = "update-meta"
)

type StoreBuckets struct {
	Current Bucket
	Journal Bucket
	Audit   Bucket
}

type Record[V any] struct {
	Action    Action    `json:"action"`
	Key       string    `json:"key"`
	Value     V         `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func PutRecord[V any](
	tx *bbolt.Tx,
	key string,
	value V,
	schema StoreBuckets,
	action Action,
) error {

	record := Record[V]{
		Action:    action,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal record: %w", err)
	}

	// Current bucket
	current, err := tx.CreateBucketIfNotExists([]byte(schema.Current))
	if err != nil {
		return fmt.Errorf("current bucket creation failed: %w", err)
	}
	if err := current.Put([]byte(key), recordBytes); err != nil {
		return fmt.Errorf("put in current bucket: %w", err)
	}

	// Journal (history) bucket
	journal, err := tx.CreateBucketIfNotExists([]byte(schema.Journal))
	if err != nil {
		return fmt.Errorf("journal bucket creation failed: %w", err)
	}
	keyBytesTS := []byte(record.Timestamp.Format(time.RFC3339Nano))
	if err := journal.Put(keyBytesTS, recordBytes); err != nil {
		return fmt.Errorf("put in journal bucket: %w", err)
	}

	// Audit log bucket
	audit, err := tx.CreateBucketIfNotExists([]byte(schema.Audit))
	if err != nil {
		return fmt.Errorf("audit bucket creation failed: %w", err)
	}
	if err := audit.Put(keyBytesTS, recordBytes); err != nil {
		return fmt.Errorf("put in audit bucket: %w", err)
	}

	return nil
}

func EnumerateBucket[V any](tx *bbolt.Tx, bucketName Bucket) ([]Record[V], error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, fmt.Errorf("bucket %s not found", bucketName)
	}

	var records []Record[V]
	err := bucket.ForEach(func(_, v []byte) error {
		var record Record[V]
		if err := json.Unmarshal(v, &record); err != nil {
			return fmt.Errorf("unmarshal record: %w", err)
		}
		records = append(records, record)
		return nil
	})
	return records, err
}

func RollbackRecord[V any](tx *bbolt.Tx, buckets StoreBuckets) error {
	journal := tx.Bucket([]byte(buckets.Journal))
	if journal == nil {
		return fmt.Errorf("journal bucket not found: %s", buckets.Journal)
	}

	// Find the latest record by iterating in reverse
	c := journal.Cursor()
	k, v := c.Last()
	if k == nil {
		return fmt.Errorf("no journal entries to rollback")
	}

	var entry Record[V]
	if err := json.Unmarshal(v, &entry); err != nil {
		return fmt.Errorf("unmarshal journal entry: %w", err)
	}

	current := tx.Bucket([]byte(buckets.Current))
	if current == nil {
		return fmt.Errorf("current bucket not found: %s", buckets.Current)
	}

	// Apply rollback based on original action
	switch entry.Action {
	case ActionCreate:
		// Created earlier, so now we delete
		if err := current.Delete([]byte(entry.Key)); err != nil {
			return fmt.Errorf("rollback delete failed: %w", err)
		}

	case ActionDelete, ActionUpdate:
		// Deleted or updated earlier, so now we restore
		valueBytes, err := json.Marshal(entry.Value)
		if err != nil {
			return fmt.Errorf("marshal restored value: %w", err)
		}
		if err := current.Put([]byte(entry.Key), valueBytes); err != nil {
			return fmt.Errorf("rollback put failed: %w", err)
		}

	case ActionRollback:
		return fmt.Errorf("cannot rollback a rollback entry")

	default:
		return fmt.Errorf("unknown action: %s", entry.Action)
	}

	// Remove the journal entry
	if err := journal.Delete(k); err != nil {
		return fmt.Errorf("delete journal entry: %w", err)
	}

	// Add audit entry
	rollbackEntry := Record[V]{
		Action:    ActionRollback,
		Key:       entry.Key,
		Value:     entry.Value,
		Timestamp: time.Now(),
	}
	audit := tx.Bucket([]byte(buckets.Audit))
	if audit == nil {
		return fmt.Errorf("audit bucket not found: %s", buckets.Audit)
	}
	auditKey := []byte(rollbackEntry.Timestamp.Format(time.RFC3339Nano))
	data, err := json.Marshal(rollbackEntry)
	if err != nil {
		return fmt.Errorf("marshal rollback audit entry: %w", err)
	}
	if err := audit.Put(auditKey, data); err != nil {
		return fmt.Errorf("write rollback audit entry: %w", err)
	}

	return nil
}
