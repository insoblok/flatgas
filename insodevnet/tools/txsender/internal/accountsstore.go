package internal

import (
	"encoding/json"
	"fmt"
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
	return filepath.Join(base, "wallet", "kvstore", "accounts.db")
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
