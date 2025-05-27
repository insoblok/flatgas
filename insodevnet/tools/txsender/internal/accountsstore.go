package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

type Action string
type Alias string
type KvAddress string

func (a Alias) String() string {
	return string(a)
}

func (a Action) String() string {
	return string(a)
}

func (a KvAddress) String() string {
	return string(a)
}

type AliasRecord struct {
	Alias    Alias                  `json:"alias"`
	Address  KvAddress              `json:"address"`
	Keystore map[string]interface{} `json:"keystore"`
	Metadata map[string]interface{} `json:"meta"`
	Created  time.Time              `json:"created"`
	Updated  time.Time              `json:"updated"`
}

func GetDBFilePath(base string) string {
	dir := filepath.Join(base, "wallet", "kvstore")
	if err := os.MkdirAll(dir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create db directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(dir, "accounts.db")
}

const (
	ActionCreate     Action = "create"
	ActionDelete     Action = "delete"
	ActionUpdateMeta Action = "update-meta"
	ActionRollback   Action = "rollback" // optional
)

type JournalEntry struct {
	Action    Action       `json:"action"`
	Alias     Alias        `json:"alias"`
	Timestamp time.Time    `json:"timestamp"`
	Data      *AliasRecord `json:"data,omitempty"`
}

// WriteJournalEntry appends a new journal entry to the journal bucket using timestamp-based key
func WriteJournalEntry(db *bbolt.DB, entry JournalEntry) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("journal"))
		if err != nil {
			return fmt.Errorf("create journal bucket: %w", err)
		}
		key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal journal entry: %w", err)
		}
		return b.Put(key, data)
	})
}
func WriteJournalEntryWithBucket(b *bbolt.Bucket, entry JournalEntry) error {
	key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal journal entry: %w", err)
	}
	return b.Put(key, data)
}

// WriteAuditLogEntry appends a new audit log entry to the auditlog bucket using timestamp-based key
func WriteAuditLogEntry(db *bbolt.DB, entry JournalEntry) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("auditlog"))
		if err != nil {
			return fmt.Errorf("create auditlog bucket: %w", err)
		}
		key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal audit log entry: %w", err)
		}
		return b.Put(key, data)
	})
}
func WriteTxAuditLogEntry(tx *bbolt.Tx, entry JournalEntry) error {
	audit := tx.Bucket([]byte("auditlog"))
	if audit == nil {
		var err error
		audit, err = tx.CreateBucket([]byte("auditlog"))
		if err != nil {
			return fmt.Errorf("create auditlog bucket: %w", err)
		}
	}

	key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal buket entry: %w", err)
	}
	return audit.Put(key, data)
}
