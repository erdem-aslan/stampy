package main
import (
	"time"
	"log"
	"errors"
	"sync"
)


type StampyCache struct {

	keyValueCache    map[string]StampyCacheEntry
	stampyCacheStats StampyCacheStats
	cacheMutex sync.RWMutex
}


func (s StampyCache) putKeyWithValue(key string, value string) {

	log.Println("Putting/Updating entry for key:", key, "with value:", value)

	s.cacheMutex.Lock()

	now := time.Now()

	s.keyValueCache[key] = StampyCacheEntry{value, now}
	s.stampyCacheStats.KeyPuts++

	s.cacheMutex.Unlock()
}

func (s StampyCache) getValueWithKey(key string) (StampyCacheEntry, error) {

	log.Println("Fetching entry with key:", key)

	s.cacheMutex.RLock()

	value, ok := s.keyValueCache[key]

	if ok {
		s.stampyCacheStats.KeyHits++
	} else {
		s.stampyCacheStats.AbsentKeyHits++
	}

	s.cacheMutex.RUnlock()

	if !ok {
		log.Println("Entry with key:", key, "not found.")
		return value, errors.New("Missing key")
	}

	log.Printf("Entry with key:%s found, %b", key, value)

	return value, nil

}

func (s StampyCache) deleteValueWithKey(key string) {

	log.Println("Deleting entry with key:", key)

	s.cacheMutex.Lock()

	delete(s.keyValueCache, key)
	s.stampyCacheStats.KeyDeletes++

	s.cacheMutex.Unlock()

}

type StampyCacheStats struct {

	KeyPuts        uint64
	KeyDeletes     uint64
	KeyHits        uint64
	AbsentKeyHits  uint64
}

type StampyCacheEntry struct {
	entryValue   string
	creationDate time.Time
}

