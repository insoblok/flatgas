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
