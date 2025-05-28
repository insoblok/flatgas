package internal

import (
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"testing"
)

type testValue struct {
	Message string `json:"message"`
}

func TestPutAndEnumerate(t *testing.T) {
	tempDir := os.TempDir()
	dbPath := filepath.Join(tempDir, "test_accounts.db")
	os.Remove(dbPath) // cleanup before

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	buckets := StoreBuckets{
		Current: "aliases",
		Journal: "journal",
		Audit:   "auditlog",
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		return PutRecord[testValue](tx, "test-key", testValue{Message: "Hello"}, buckets, ActionCreate)
	})
	if err != nil {
		t.Fatalf("failed to put record: %v", err)
	}

	err = db.View(func(tx *bbolt.Tx) error {
		records, err := EnumerateBucket[testValue](tx, buckets.Current)
		if err != nil {
			t.Fatalf("failed to enumerate current bucket: %v", err)
		}
		if len(records) != 1 || records[0].Value.Message != "Hello" {
			t.Fatalf("unexpected record data: %+v", records)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("view transaction failed: %v", err)
	}
}
