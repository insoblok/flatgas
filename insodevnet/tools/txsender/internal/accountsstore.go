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
	Metadata map[string]string      `json:"meta"`
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

func WriteTxHistoryLogEntry(tx *bbolt.Tx, entry JournalEntry) error {
	journal := tx.Bucket([]byte("journal"))
	if journal == nil {
		var err error
		journal, err = tx.CreateBucket([]byte("journal"))
		if err != nil {
			return fmt.Errorf("create journal bucket: %w", err)
		}
	}

	key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal buket entry: %w", err)
	}
	return journal.Put(key, data)
}

func ReadAlias(db *bbolt.DB, alias string) (*AliasRecord, error) {
	var record AliasRecord

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("aliases"))
		if bucket == nil {
			return fmt.Errorf("aliases bucket not found")
		}

		data := bucket.Get([]byte(alias))
		if data == nil {
			return fmt.Errorf("alias '%s' not found", alias)
		}

		if err := json.Unmarshal(data, &record); err != nil {
			return fmt.Errorf("failed to decode alias data: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &record, nil
}

// WithUpdateAlias loads an alias, modifies it via fn, and saves it back.
func WithUpdateAlias(db *bbolt.DB, alias string, fn func(*AliasRecord) error) error {
	return db.Update(func(tx *bbolt.Tx) error {
		aliases := tx.Bucket([]byte("aliases"))
		if aliases == nil {
			return fmt.Errorf("aliases bucket not found")
		}

		data := aliases.Get([]byte(alias))
		if data == nil {
			return fmt.Errorf("alias not found: %s", alias)
		}

		var record AliasRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return fmt.Errorf("unmarshal alias record: %w", err)
		}

		// Modify in-place
		if err := fn(&record); err != nil {
			return err
		}

		record.Updated = time.Now()

		newData, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("marshal updated record: %w", err)
		}

		if err := aliases.Put([]byte(alias), newData); err != nil {
			return fmt.Errorf("update alias: %w", err)
		}

		// Record journal entry
		entry := JournalEntry{
			Action:    ActionUpdateMeta,
			Alias:     Alias(alias),
			Timestamp: time.Now(),
			Data:      &record,
		}
		err = WriteTxHistoryLogEntry(tx, entry)
		if err != nil {
			return fmt.Errorf("write journal entry: %w", err)
		}

		return WriteTxAuditLogEntry(tx, entry)
	})
}

// SaveAliasRecord saves an AliasRecord under the given alias,
// and appends the action to both journal and auditlog buckets.
func SaveAliasRecord(tx *bbolt.Tx, alias string, record AliasRecord, action Action) error {
	// --- Aliases bucket ---
	aliases := tx.Bucket([]byte("aliases"))
	if aliases == nil {
		return fmt.Errorf("aliases bucket not found")
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal alias record: %w", err)
	}
	if err := aliases.Put([]byte(alias), data); err != nil {
		return fmt.Errorf("failed to store alias record: %w", err)
	}

	// --- Journal bucket ---
	journal := tx.Bucket([]byte("journal"))
	if journal == nil {
		j, err := tx.CreateBucketIfNotExists([]byte("journal"))
		if err != nil {
			return fmt.Errorf("failed to create journal bucket: %w", err)
		}
		journal = j
	}

	journalEntry := JournalEntry{
		Action:    action,
		Alias:     Alias(alias),
		Timestamp: time.Now(),
		Data:      &record,
	}

	if err := WriteTxHistoryLogEntry(tx, journalEntry); err != nil {
		return fmt.Errorf("failed to write hisotry log: %w", err)
	}

	// --- Audit log bucket ---
	if err := WriteTxAuditLogEntry(tx, journalEntry); err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}

	return nil
}
