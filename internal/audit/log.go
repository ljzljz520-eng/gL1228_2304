package audit

import (
	"fmt"
	"stagebeam/internal/model"
	"sync"
	"time"
)

type Log struct {
	mu      sync.RWMutex
	entries []model.AuditEntry
	nextID  int64
}

func NewLog() *Log {
	return &Log{entries: make([]model.AuditEntry, 0, 32)}
}

func (log *Log) Record(showID, action, detail string) model.AuditEntry {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.nextID++
	entry := model.AuditEntry{ID: fmt.Sprintf("audit-%06d", log.nextID), ShowID: showID, Action: action, Detail: detail, Sequence: log.nextID, CreatedAt: time.Unix(0, log.nextID*int64(time.Millisecond)).UTC()}
	log.entries = append(log.entries, entry)
	return entry
}

func (log *Log) ForShow(showID string) []model.AuditEntry {
	log.mu.RLock()
	defer log.mu.RUnlock()
	result := make([]model.AuditEntry, 0)
	for _, entry := range log.entries {
		if entry.ShowID == showID {
			result = append(result, entry)
		}
	}
	return result
}

func (log *Log) Latest(showID string) (model.AuditEntry, bool) {
	entries := log.ForShow(showID)
	if len(entries) == 0 {
		return model.AuditEntry{}, false
	}
	return entries[len(entries)-1], true
}

func (log *Log) Count() int {
	log.mu.RLock()
	defer log.mu.RUnlock()
	return len(log.entries)
}

func (log *Log) Since(sequence int64) []model.AuditEntry {
	log.mu.RLock()
	defer log.mu.RUnlock()
	result := make([]model.AuditEntry, 0)
	for _, entry := range log.entries {
		if entry.Sequence > sequence {
			result = append(result, entry)
		}
	}
	return result
}

func (log *Log) Actions(showID string) []string {
	entries := log.ForShow(showID)
	actions := make([]string, 0, len(entries))
	for _, entry := range entries {
		actions = append(actions, entry.Action)
	}
	return actions
}

func (log *Log) ClearShow(showID string) int {
	log.mu.Lock()
	defer log.mu.Unlock()
	kept := log.entries[:0]
	removed := 0
	for _, entry := range log.entries {
		if entry.ShowID == showID {
			removed++
			continue
		}
		kept = append(kept, entry)
	}
	log.entries = kept
	return removed
}
