package internal

import (
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
	"os"
	"testing"
	"time"
)

func TestFlatgasStorageLifecycle(t *testing.T) {
	db, err := bbolt.Open("/tmp/test_db", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("/tmp/test_db")

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

func TestCreateRecord(t *testing.T) {
	dbPath := "/tmp/test_flatgas_create.db"
	defer os.Remove(dbPath)

	db, err := bbolt.Open(dbPath, 0666, nil)
	require.NoError(t, err)
	defer db.Close()

	type MyValue struct {
		Data string
	}

	schema := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	key := "user:123"
	value := MyValue{Data: "hello"}

	// Create
	err = db.Update(func(tx *bbolt.Tx) error {
		return CreateRecord(tx, key, value, schema)
	})
	require.NoError(t, err)

	// Validate record in Current
	err = db.View(func(tx *bbolt.Tx) error {
		rec, err := GetRecord[MyValue](tx, key, schema)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, ActionCreate, rec.Action)
		require.Equal(t, key, rec.Key)
		require.Equal(t, value.Data, rec.Value.Data)
		require.WithinDuration(t, time.Now(), rec.Timestamp, time.Second)
		return nil
	})

	err = db.View(func(tx *bbolt.Tx) error {
		audit, err := GetAuditRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, audit, 1)
		require.Equal(t, ActionCreate, audit[0].Action)

		journal, err := GetJournalRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, journal, 1)
		require.Equal(t, ActionCreate, journal[0].Action)

		return nil
	})

	err = db.Update(func(tx *bbolt.Tx) error {
		return CreateRecord(tx, key, value, schema)
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "already exists")
}

func TestCreateRollback(t *testing.T) {
	dbPath := "/tmp/test_flatgas_create_rollback.db"
	defer os.Remove(dbPath)

	db, err := bbolt.Open(dbPath, 0666, nil)
	require.NoError(t, err)
	defer db.Close()

	type MyValue struct {
		Data string
	}

	schema := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	key := "item:1"
	value := MyValue{Data: "alpha"}

	err = db.Update(func(tx *bbolt.Tx) error {
		return CreateRecord(tx, key, value, schema)
	})
	require.NoError(t, err)

	err = db.Update(func(tx *bbolt.Tx) error {
		return RollbackRecord[MyValue](tx, schema)
	})
	require.NoError(t, err)

	err = db.View(func(tx *bbolt.Tx) error {
		rec, err := GetRecord[MyValue](tx, key, schema)
		require.NoError(t, err)
		require.Nil(t, rec)

		current, err := GetCurrentRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, current, 0)

		journal, err := GetJournalRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, journal, 0)

		audit, err := GetAuditRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, audit, 2)
		require.Equal(t, ActionCreate, audit[0].Action)
		require.Equal(t, ActionRollback, audit[1].Action)

		return nil
	})
}

func TestUpdateRecord(t *testing.T) {
	dbPath := "/tmp/test_flatgas_update.db"
	defer os.Remove(dbPath)

	db, err := bbolt.Open(dbPath, 0666, nil)
	require.NoError(t, err)
	defer db.Close()

	type MyValue struct {
		Data string
	}

	schema := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	key := "item:42"
	initial := MyValue{Data: "before"}
	updated := MyValue{Data: "after"}

	err = db.Update(func(tx *bbolt.Tx) error {
		return UpdateRecord(tx, key, updated, schema)
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")

	err = db.Update(func(tx *bbolt.Tx) error {
		return CreateRecord(tx, key, initial, schema)
	})
	require.NoError(t, err)

	err = db.Update(func(tx *bbolt.Tx) error {
		return UpdateRecord(tx, key, updated, schema)
	})
	require.NoError(t, err)

	err = db.View(func(tx *bbolt.Tx) error {
		rec, err := GetRecord[MyValue](tx, key, schema)
		require.NoError(t, err)
		require.NotNil(t, rec)
		require.Equal(t, ActionUpdate, rec.Action)
		require.Equal(t, updated.Data, rec.Value.Data)

		journal, err := GetJournalRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, journal, 2)
		require.Equal(t, ActionCreate, journal[0].Action)
		require.Equal(t, ActionUpdate, journal[1].Action)

		audit, err := GetAuditRecords[MyValue](tx, schema)
		require.NoError(t, err)
		require.Len(t, audit, 2)
		require.Equal(t, ActionCreate, audit[0].Action)
		require.Equal(t, ActionUpdate, audit[1].Action)

		return nil
	})
}

func TestDeleteNonExistentRecord(t *testing.T) {
	dbPath := "/tmp/test_flatgas_delete_nonexistent.db"
	defer os.Remove(dbPath)

	db, err := bbolt.Open(dbPath, 0666, nil)
	require.NoError(t, err)
	defer db.Close()

	type MyValue struct {
		Data string
	}

	schema := StoreBuckets{
		Current: "current",
		Journal: "journal",
		Audit:   "audit",
	}

	key := "ghost:001"

	err = db.Update(func(tx *bbolt.Tx) error {
		return DeleteRecord[MyValue](tx, key, schema)
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}
