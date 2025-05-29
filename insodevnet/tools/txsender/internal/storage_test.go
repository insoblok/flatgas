package internal

import (
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
	"os"
	"testing"
)

func TestFlatgasStorageLifecycle(t *testing.T) {
	db, err := bbolt.Open("/tmp/test_flatgas.db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test_flatgas.db")

	type MyValue struct {
		Data string
	}

	buckets := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	key := "user:123"
	val1 := MyValue{Data: "initial"}
	val2 := MyValue{Data: "updated"}

	// Create
	err = db.Update(func(tx *bbolt.Tx) error {
		return CreateRecord(tx, key, val1, buckets)
	})
	require.NoError(t, err)

	// Read and check
	err = db.View(func(tx *bbolt.Tx) error {
		rec, err := GetRecord[MyValue](tx, key, buckets)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, ActionCreate, rec.Action)
		require.Equal(t, val1.Data, rec.Value.Data)
		return nil
	})

	// Update
	err = db.Update(func(tx *bbolt.Tx) error {
		return UpdateRecord(tx, key, val2, buckets)
	})
	require.NoError(t, err)

	// Delete
	err = db.Update(func(tx *bbolt.Tx) error {
		return DeleteRecord[MyValue](tx, key, buckets)
	})
	require.NoError(t, err)

	// Rollback 3x (Delete → Update → Create)
	for i := 0; i < 3; i++ {
		err = db.Update(func(tx *bbolt.Tx) error {
			return RollbackRecord[MyValue](tx, buckets)
		})
		require.NoError(t, err)
	}

	// Final check: should be back to original value
	err = db.View(func(tx *bbolt.Tx) error {
		rec, err := GetRecord[MyValue](tx, key, buckets)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, ActionCreate, rec.Action)
		require.Equal(t, val1.Data, rec.Value.Data)

		curr, _ := GetCurrentRecords[MyValue](tx, buckets)
		journal, _ := GetJournalRecords[MyValue](tx, buckets)
		audit, _ := GetAuditRecords[MyValue](tx, buckets)

		require.Len(t, curr, 1)
		require.Len(t, journal, 0)
		require.Len(t, audit, 6) // Create, Update, Delete, Rollback x3

		return nil
	})
}
