package shardmanager

import (
	"errors"
	"fmt"
	"math"
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

type Shard struct {
	Address string
	Number  int
}

type Manager struct {
	ss *sync.Map
	sr *sync.Map
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager() *Manager {
	return &Manager{
		ss: &sync.Map{},
		sr: &sync.Map{},
	}
}
func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}

func (m *Manager) AddReplica(s *Shard) {
	m.sr.Store(s.Number, s)
}

func (m *Manager) Shard(orderId string) (*Shard, error) {

	length := lenSyncMap(m.ss)

	hash := hash.Int8String(orderId)

	output := hash % int8(length)

	fmt.Println("Shard: ", math.Abs(float64(output)))

	if s, ok := m.ss.Load(int(math.Abs(float64(output)))); ok {

		return s.(*Shard), nil
	}

	return nil, ErrorShardNotFound
}

func (m *Manager) ShardReplica(orderId string) (*Shard, error) {

	length := lenSyncMap(m.sr)

	hash := hash.Int8String(orderId)

	output := hash % int8(length)

	fmt.Println("Replica: ", math.Abs(float64(output)))

	if s, ok := m.sr.Load(int(math.Abs(float64(output)))); ok {

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
