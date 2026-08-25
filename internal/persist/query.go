package persist

import (
	"fmt"
	"go.etcd.io/bbolt"
	"stagebeam/internal/model"
)

func (store *Store) ListShows() ([]model.Show, error) {
	shows := make([]model.Show, 0)
	err := store.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(showsBucket).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var record model.ShowRecord
			if err := decode(value, &record); err != nil {
				return err
			}
			shows = append(shows, recordToShow(record))
			return nil
		})
	})
	return shows, err
}

func (store *Store) SaveAudit(entry model.AuditEntry) error {
	data, err := encode(entry)
	if err != nil {
		return err
	}
	bucketName := []byte("audit")
	return store.db.Update(func(tx *bbolt.Tx) error {
		bucket, bucketErr := tx.CreateBucketIfNotExists(bucketName)
		if bucketErr != nil {
			return bucketErr
		}
		key := fmt.Sprintf("%s:%06d", entry.ShowID, entry.Sequence)
		return bucket.Put([]byte(key), data)
	})
}

func (store *Store) AuditForShow(showID string) ([]model.AuditEntry, error) {
	entries := make([]model.AuditEntry, 0)
	err := store.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("audit"))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(key, value []byte) error {
			if len(key) < len(showID) || string(key[:len(showID)]) != showID {
				return nil
			}
			var entry model.AuditEntry
			if err := decode(value, &entry); err != nil {
				return err
			}
			entries = append(entries, entry)
			return nil
		})
	})
	return entries, err
}
