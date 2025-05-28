package internal

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"testing"
)

type testValue struct {
	Message string `json:"message"`
}

func TestPutAndEnumerateRecord(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	buckets := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "auditlog",
	}

	tval := testValue{Message: "Hello"}
	err = db.Update(func(tx *bbolt.Tx) error {
		return PutRecord[testValue](tx, "test-key", tval, buckets, ActionCreate)
	})
	if err != nil {
		t.Fatalf("Failed to put record: %v", err)
	}

	t.Run("CheckJournalEntry", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(buckets.Journal))
			if bucket == nil {
				t.Fatal("Journal bucket not found")
			}
			var found bool
			err := bucket.ForEach(func(k, v []byte) error {
				var rec Record[testValue]
				if err := json.Unmarshal(v, &rec); err != nil {
					t.Errorf("Invalid record JSON: %v", err)
				}
				if rec.Action == ActionCreate && rec.Key == "test-key" && rec.Value.Message == "Hello" {
					found = true
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !found {
				t.Error("Journal entry not found")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to check journal: %v", err)
		}
	})

	t.Run("EnumerateCurrent", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			records, err := EnumerateBucket[testValue](tx, buckets.Current)
			if err != nil {
				return fmt.Errorf("enumerate current failed: %w", err)
			}

			found := false
			for _, rec := range records {
				if rec.Key == "test-key" && rec.Value.Message == "Hello" {
					found = true
				}
			}
			if !found {
				t.Errorf("Expected record not found in current bucket")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Transaction failed: %v", err)
		}
	})

	t.Run("CheckAuditLogEntry", func(t *testing.T) {
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(buckets.Audit))
			if bucket == nil {
				t.Fatal("Audit bucket not found")
			}
			var found bool
			err := bucket.ForEach(func(k, v []byte) error {
				var rec Record[testValue]
				if err := json.Unmarshal(v, &rec); err != nil {
					t.Errorf("Invalid record JSON: %v", err)
				}
				if rec.Action == ActionCreate && rec.Key == "test-key" && rec.Value.Message == "Hello" {
					found = true
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !found {
				t.Error("Audit log entry not found")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to check audit log: %v", err)
		}
	})
}
