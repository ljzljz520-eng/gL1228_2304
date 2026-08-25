package persist

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"stagebeam/internal/model"
)

var ErrNotFound = errors.New("show not found")

var (
	showsBucket    = []byte("shows")
	settingsBucket = []byte("settings")
	gesturesBucket = []byte("gestures")
	layersBucket   = []byte("layers")
)

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open stage store: %w", err)
	}
	store := &Store{db: db}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (store *Store) initialize() error {
	return store.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{showsBucket, settingsBucket, gesturesBucket, layersBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("create bucket: %w", err)
			}
		}
		return nil
	})
}

func (store *Store) Close() error {
	if store == nil || store.db == nil {
		return nil
	}
	return store.db.Close()
}

func (store *Store) SaveShow(show model.Show) error {
	data, err := encode(showToRecord(show))
	if err != nil {
		return err
	}
	return store.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(showsBucket).Put([]byte(show.ID), data); err != nil {
			return fmt.Errorf("save show: %w", err)
		}
		settingData, settingErr := encode(model.SettingRecord{ShowID: show.ID, Settings: show.Settings})
		if settingErr != nil {
			return settingErr
		}
		if err := tx.Bucket(settingsBucket).Put([]byte(show.ID), settingData); err != nil {
			return fmt.Errorf("save settings: %w", err)
		}
		return nil
	})
}

func (store *Store) LoadShow(id string) (model.Show, error) {
	var record model.ShowRecord
	err := store.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(showsBucket).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &record)
	})
	if err != nil {
		return model.Show{}, err
	}
	return recordToShow(record), nil
}

func (store *Store) SaveGesture(record model.GestureRecord) error {
	data, err := encode(record)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%d", record.ShowID, record.Sequence)
	return store.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(gesturesBucket).Put([]byte(key), data) })
}

func (store *Store) SaveLayerSnapshot(snapshot model.LayerSnapshot) error {
	data, err := encode(snapshot)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%d", snapshot.ShowID, snapshot.Frame)
	return store.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(layersBucket).Put([]byte(key), data) })
}

func (store *Store) LoadSetting(id string) (model.SettingRecord, error) {
	var record model.SettingRecord
	err := store.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(settingsBucket).Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &record)
	})
	return record, err
}

func (store *Store) RemoveShow(id string) error {
	return store.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(showsBucket).Delete([]byte(id)); err != nil {
			return err
		}
		if err := tx.Bucket(settingsBucket).Delete([]byte(id)); err != nil {
			return err
		}
		return nil
	})
}

func (store *Store) SaveSetting(record model.SettingRecord) error {
	data, err := encode(record)
	if err != nil {
		return err
	}
	return store.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(settingsBucket).Put([]byte(record.ShowID), data)
	})
}

func (store *Store) Health() error {
	if store == nil || store.db == nil {
		return fmt.Errorf("store is closed")
	}
	return store.db.View(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{showsBucket, settingsBucket, gesturesBucket, layersBucket} {
			if tx.Bucket(bucket) == nil {
				return fmt.Errorf("bucket %s is missing", bucket)
			}
		}
		return nil
	})
}
