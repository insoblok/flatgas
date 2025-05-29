package internal

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"path/filepath"
	"testing"
)

type testValue struct {
	Message string `json:"message"`
}

func TestPutAndEnumerateRecord(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	buckets := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	t.Run("PutRecord", func(t *testing.T) {
		err = db.Update(func(tx *bbolt.Tx) error {
			return PutRecord[testValue](tx, "test-key", testValue{Message: "Hello"}, buckets, ActionCreate)
		})
		if err != nil {
			t.Fatalf("Failed to put record: %v", err)
		}
	})

	t.Run("CheckJournalEntry", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			journal := tx.Bucket([]byte(buckets.Journal))
			if journal == nil {
				t.Fatal("Journal bucket not found")
			}
			cursor := journal.Cursor()
			k, v := cursor.First()
			if k == nil {
				t.Fatal("No journal entry found")
			}
			var rec Record[testValue]
			if err := json.Unmarshal(v, &rec); err != nil {
				t.Fatalf("Failed to unmarshal journal entry: %v", err)
			}
			if rec.Action != ActionCreate || rec.Key != "test-key" || rec.Value.Message != "Hello" {
				t.Errorf("Unexpected journal record: %+v", rec)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to verify journal: %v", err)
		}
	})

	t.Run("EnumerateCurrent", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			records, err := EnumerateBucket[testValue](tx, buckets.Current)
			if err != nil {
				t.Fatalf("Enumeration failed: %v", err)
			}
			if len(records) != 1 || records[0].Key != "test-key" || records[0].Value.Message != "Hello" {
				t.Errorf("Unexpected records: %+v", records)
			}
			return nil
		})
	})

	t.Run("CheckAuditLogEntry", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			audit := tx.Bucket([]byte(buckets.Audit))
			if audit == nil {
				t.Fatal("Audit bucket not found")
			}
			cursor := audit.Cursor()
			k, v := cursor.First()
			if k == nil {
				t.Fatal("No audit entry found")
			}
			var rec Record[testValue]
			if err := json.Unmarshal(v, &rec); err != nil {
				t.Fatalf("Failed to unmarshal audit entry: %v", err)
			}
			if rec.Action != ActionCreate || rec.Key != "test-key" || rec.Value.Message != "Hello" {
				t.Errorf("Unexpected audit record: %+v", rec)
			}
			return nil
		})
	})

	t.Run("DeleteRecord", func(t *testing.T) {
		err = db.Update(func(tx *bbolt.Tx) error {
			return DeleteRecord[testValue](tx, "test-key", buckets)
		})
		if err != nil {
			t.Fatalf("DeleteRecord failed: %v", err)
		}

		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(buckets.Current))
			if bucket == nil {
				t.Fatal("current bucket not found after delete")
			}
			val := bucket.Get([]byte("test-key"))
			if val != nil {
				t.Errorf("Expected key to be deleted, but found value: %s", val)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error verifying delete: %v", err)
		}
	})
}

func TestRollbackCreate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	buckets := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	// Put an entry first
	err = db.Update(func(tx *bbolt.Tx) error {
		return PutRecord[testValue](tx, "key1", testValue{Message: "to-delete"}, buckets, ActionCreate)
	})
	if err != nil {
		t.Fatalf("PutRecord failed: %v", err)
	}

	// Rollback it
	err = db.Update(func(tx *bbolt.Tx) error {
		return RollbackRecord[testValue](tx, buckets)
	})
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Confirm that the entry is no longer in the current bucket
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(buckets.Current))
		if b == nil {
			t.Fatal("Current bucket not found")
		}
		v := b.Get([]byte("key1"))
		if v != nil {
			t.Errorf("Expected key1 to be deleted, found: %s", string(v))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}
}
