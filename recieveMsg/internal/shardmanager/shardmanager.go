package shardmanager

import (
	"errors"
	"sync"
)

var ShardNow int

type Shard struct {
	Address string
	Number  int
}

type Manager struct {
	size int
	ss   *sync.Map
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
	}
}
func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}
func (m *Manager) Shard() (*Shard, error) {
	n := ShardNow
	length := lenSyncMap(m.ss)
	if n > length-1 {
		n = 0
		ShardNow = 0
	}
	if s, ok := m.ss.Load(n); ok {
		ShardNow++
		return s.(*Shard), nil
	}

	return nil, ErrorShardNotFound
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}
