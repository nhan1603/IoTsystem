package iotsystem

import (
	"sync"
	"time"
)

// Per run/process, keep this around (e.g., as a field on cassandraImpl)
type TsDeduper struct {
	mu   sync.Mutex
	last map[string]int64 // device_id -> last used ms since epoch
}

func NewTsDeduper() *TsDeduper { return &TsDeduper{last: make(map[string]int64)} }

// Next returns a time that is unique & monotonic-per-device at ms resolution.
func (d *TsDeduper) Next(deviceID string, candidate time.Time) time.Time {
	ms := candidate.UnixMilli()
	d.mu.Lock()
	defer d.mu.Unlock()
	if last, ok := d.last[deviceID]; ok && ms <= last {
		ms = last + 1 // bump by 1ms (timestamp type only stores ms)
	}
	d.last[deviceID] = ms
	return time.UnixMilli(ms)
}
