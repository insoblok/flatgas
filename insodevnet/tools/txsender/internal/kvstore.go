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

type Buckets struct {
	Current *bbolt.Bucket
	Journal *bbolt.Bucket
	Audit   *bbolt.Bucket
}

func getBuckets(tx *bbolt.Tx, schema StoreBuckets, allowCreate bool) (*Buckets, error) {
	var get func(name Bucket) (*bbolt.Bucket, error)

	if allowCreate {
		get = func(name Bucket) (*bbolt.Bucket, error) {
			return tx.CreateBucketIfNotExists([]byte(name))
		}
	} else {
		get = func(name Bucket) (*bbolt.Bucket, error) {
			bkt := tx.Bucket([]byte(name))
			if bkt == nil {
				return nil, fmt.Errorf("bucket '%s' not found", name)
			}
			return bkt, nil
		}
	}

	current, err := get(schema.Current)
	if err != nil {
		return nil, err
	}
	journal, err := get(schema.Journal)
	if err != nil {
		return nil, err
	}
	audit, err := get(schema.Audit)
	if err != nil {
		return nil, err
	}

	return &Buckets{Current: current, Journal: journal, Audit: audit}, nil
}

func CreateRecord[V any](
	tx *bbolt.Tx,
	key string,
	value V,
	schema StoreBuckets,
) error {

	buckets, err := getBuckets(tx, schema, true)
	if err != nil {
		return fmt.Errorf("Buckets failed: %w", err)
	}

	if buckets.Current.Get([]byte(key)) != nil {
		return fmt.Errorf("key '%s' already exists", key)
	}

	record := Record[V]{
		Action:    ActionCreate,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal record: %w", err)
	}

	if err = buckets.Current.Put([]byte(key), recordBytes); err != nil {
		return fmt.Errorf("put in current bucket: %w", err)
	}

	keyBytesTS := []byte(record.Timestamp.Format(time.RFC3339Nano))
	if err = buckets.Journal.Put(keyBytesTS, recordBytes); err != nil {
		return fmt.Errorf("put in journal bucket: %w", err)
	}

	if err = buckets.Audit.Put(keyBytesTS, recordBytes); err != nil {
		return fmt.Errorf("put in audit bucket: %w", err)
	}

	return nil
}

func UpdateRecord[V any](
	tx *bbolt.Tx,
	key string,
	value V,
	schema StoreBuckets,
) error {
	buckets, err := getBuckets(tx, schema, false)
	if err != nil {
		return fmt.Errorf("Buckets failed: %w", err)
	}

	original := buckets.Current.Get([]byte(key))
	if original == nil {
		return fmt.Errorf("key '%s' does not exist", key)
	}

	// 🧠 Capture pre-update record
	var prev Record[V]
	if err := json.Unmarshal(original, &prev); err != nil {
		return fmt.Errorf("unmarshal current value: %w", err)
	}

	timestamp := time.Now()

	// 🔁 Save the old value into journal (so rollback can restore it)
	journalEntry := Record[V]{
		Action:    ActionUpdate,
		Key:       key,
		Value:     prev.Value,
		Timestamp: timestamp,
	}
	journalBytes, err := json.Marshal(journalEntry)
	if err != nil {
		return fmt.Errorf("marshal journal entry: %w", err)
	}
	keyBytesTS := []byte(timestamp.Format(time.RFC3339Nano))
	if err = buckets.Journal.Put(keyBytesTS, journalBytes); err != nil {
		return fmt.Errorf("put in journal: %w", err)
	}

	// ✅ Write updated value into current
	newRecord := Record[V]{
		Action:    ActionUpdate,
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
	}
	currentBytes, err := json.Marshal(newRecord)
	if err != nil {
		return fmt.Errorf("marshal updated record: %w", err)
	}
	if err = buckets.Current.Put([]byte(key), currentBytes); err != nil {
		return fmt.Errorf("put in current bucket: %w", err)
	}

	// 📝 Always log new state to audit
	if err = buckets.Audit.Put(keyBytesTS, currentBytes); err != nil {
		return fmt.Errorf("put in audit: %w", err)
	}

	return nil
}

func DeleteRecord[V any](tx *bbolt.Tx, key string, schema StoreBuckets) error {

	buckets, err := getBuckets(tx, schema, false)
	if err != nil {
		return fmt.Errorf("Buckets failed: %w", err)
	}

	original := buckets.Current.Get([]byte(key))
	if original == nil {
		return fmt.Errorf("key '%s' not found for delete", key)
	}

	var value V
	if err = json.Unmarshal(original, &value); err != nil {
		return fmt.Errorf("unmarshal current value: %w", err)
	}

	timestamp := time.Now()
	record := Record[V]{
		Action:    ActionDelete,
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal delete record: %w", err)
	}

	if err = buckets.Current.Delete([]byte(key)); err != nil {
		return fmt.Errorf("delete from current: %w", err)
	}
	keyTs := []byte(timestamp.Format(time.RFC3339Nano))
	if err = buckets.Journal.Put(keyTs, recordBytes); err != nil {
		return fmt.Errorf("write journal entry: %w", err)
	}
	if err = buckets.Audit.Put(keyTs, recordBytes); err != nil {
		return fmt.Errorf("write audit entry: %w", err)
	}

	return nil
}

func GetRecord[V any](tx *bbolt.Tx, key string, schema StoreBuckets) (*Record[V], error) {
	buckets, err := getBuckets(tx, schema, false)
	if err != nil {
		return nil, fmt.Errorf("Buckets failed: %w", err)
	}

	original := buckets.Current.Get([]byte(key))
	if original == nil {
		return nil, nil
	}

	var value Record[V]
	if err = json.Unmarshal(original, &value); err != nil {
		return nil, fmt.Errorf("unmarshal current value: %w", err)
	}
	return &value, nil
}

func RollbackRecord[V any](tx *bbolt.Tx, schema StoreBuckets) error {
	buckets, err := getBuckets(tx, schema, false)
	if err != nil {
		return fmt.Errorf("Buckets failed: %w", err)
	}

	c := buckets.Journal.Cursor()
	k, v := c.Last()
	if k == nil {
		return fmt.Errorf("no journal entries to rollback")
	}

	var entry Record[V]
	if err = json.Unmarshal(v, &entry); err != nil {
		return fmt.Errorf("unmarshal journal entry: %w", err)
	}

	switch entry.Action {
	case ActionCreate:
		// Created earlier, so now we delete
		if err = buckets.Current.Delete([]byte(entry.Key)); err != nil {
			return fmt.Errorf("rollback delete failed: %w", err)
		}

	case ActionDelete, ActionUpdate:
		// Deleted or updated earlier, so now we restore
		valueBytes, err := json.Marshal(entry.Value)
		if err != nil {
			return fmt.Errorf("marshal restored value: %w", err)
		}
		if err := buckets.Current.Put([]byte(entry.Key), valueBytes); err != nil {
			return fmt.Errorf("rollback put failed: %w", err)
		}

	case ActionRollback:
		return fmt.Errorf("cannot rollback a rollback entry")

	default:
		return fmt.Errorf("unknown action: %s", entry.Action)
	}

	if err = buckets.Journal.Delete(k); err != nil {
		return fmt.Errorf("delete journal entry: %w", err)
	}

	rollbackEntry := Record[V]{
		Action:    ActionRollback,
		Key:       entry.Key,
		Value:     entry.Value,
		Timestamp: time.Now(),
	}

	auditKey := []byte(rollbackEntry.Timestamp.Format(time.RFC3339Nano))
	data, err := json.Marshal(rollbackEntry)
	if err != nil {
		return fmt.Errorf("marshal rollback audit entry: %w", err)
	}
	if err := buckets.Audit.Put(auditKey, data); err != nil {
		return fmt.Errorf("write rollback audit entry: %w", err)
	}

	return nil
}

func GetCurrentRecords[V any](tx *bbolt.Tx, schema StoreBuckets) ([]Record[V], error) {
	return GetBucketRecords[V](tx, tx.Bucket([]byte(schema.Current)))
}
func GetJournalRecords[V any](tx *bbolt.Tx, schema StoreBuckets) ([]Record[V], error) {
	return GetBucketRecords[V](tx, tx.Bucket([]byte(schema.Journal)))
}
func GetAuditRecords[V any](tx *bbolt.Tx, schema StoreBuckets) ([]Record[V], error) {
	return GetBucketRecords[V](tx, tx.Bucket([]byte(schema.Audit)))
}

func GetBucketRecords[V any](tx *bbolt.Tx, bucket *bbolt.Bucket) ([]Record[V], error) {

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

	if action == ActionCreate && current.Get([]byte(key)) != nil {
		return fmt.Errorf("key '%s' already exists", key)
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
