package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func (s *SyncDict) Get(key string) (val interface{}, exists bool) {
	value, ok := s.m.Load(key)
	return value, ok
}

func (s *SyncDict) Len() int {
	len := 0
	s.m.Range(func(key, value any) bool {
		len++
		return true
	})
	return len
}

func (s *SyncDict) Put(key string, value interface{}) (result int) {
	_, ok := s.m.Load(key)
	s.m.Store(key, value)
	if ok {
		return 0
	}
	return 1
}

func (s *SyncDict) PutIfAbsent(key string, value interface{}) (result int) {
	_, ok := s.m.Load(key)
	if ok {
		return 0
	}
	s.m.Store(key, value)
	return 1
}

func (s *SyncDict) PutIfExists(key string, value interface{}) (result int) {
	_, ok := s.m.Load(key)
	if ok {
		s.m.Store(key, value)
		return 1
	}
	return 0
}

func (s *SyncDict) Remove(key string) (result int) {
	_, ok := s.m.Load(key)
	s.m.Delete(key)
	if ok {
		return 1
	}
	return 0
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncDict) Keys() []string {
	result := make([]string, s.Len())
	i := 0
	s.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

func (s *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, s.Len())
	for i := 0; i < limit; i++ {
		s.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (s *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, s.Len())
	i := 0
	s.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return i != limit
	})
	return result
}

func (s *SyncDict) Clear() {
	*s = *MakeSyncDict()
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}
