package main

import (
	"time"
	"log"
	"errors"
	"sync"
)


type StampyBucket struct {

	bucketIndex       int
	keyValueCache     map[string]StampyBucketEntry
	ttlIndex          map[string]bool

	stampyBucketStats *StampyBucketStats

	cacheMutex        sync.RWMutex
	ttlCacheMutex     sync.RWMutex
}


func (s *StampyBucket) processKeyViaRequest(key string, r *StampyCacheRequest) {

	now := time.Now()
	var cacheEntry StampyBucketEntry

	cacheEntry.EntryValue = r.Value
	cacheEntry.CreationTime = now
	cacheEntry.LastAccessedTime = now

	log.Println("Putting/Updating entry for key:", key, "with value:", r.Value)

	if r.TimeToLive != 0 {

		log.Println("TTL value has been provided, valid until: ", r.getExpiryDate())

		s.ttlCacheMutex.Lock()

		cacheEntry.ExpiryTime = r.getExpiryDate()

		s.keyValueCache[key] = cacheEntry
		s.ttlIndex[key] = true

		s.ttlCacheMutex.Unlock()
		return
	}

	s.cacheMutex.Lock()

	s.keyValueCache[key] = cacheEntry
	s.stampyBucketStats.KeyPuts++

	s.cacheMutex.Unlock()
}

func (s *StampyBucket) getValueWithKey(key string) (StampyBucketEntry, error) {

	log.Println("Fetching entry with key:", key)

	now := time.Now()

	if s.ttlIndex[key] {
		// key with ttl set
		s.ttlCacheMutex.RLock()
		value, ok := s.keyValueCache[key]

		if ok {

			if value.ExpiryTime.Before(now) {
				// key valid but expired
				s.stampyBucketStats.incrementExpiredKeyHits()
				log.Println("Entry with key:", key, "has been expired.")

				defer s.deleteValueWithKeyIfPresent(key)
				return value, errors.New("Expired key")
			}

			// key valid
			s.stampyBucketStats.incrementKeyHits()
			value.LastAccessedTime = now
			return value, nil
		}

		// missing key
		return value, errors.New("Missing key")

	}

	s.cacheMutex.RLock()

	value, ok := s.keyValueCache[key]

	if ok {
		s.stampyBucketStats.incrementKeyHits()
	} else {
		s.stampyBucketStats.incrementAbsentKeyHits()
	}

	s.cacheMutex.RUnlock()

	if !ok {
		log.Println("Entry with key:", key, "not found.")
		return value, errors.New("Missing key")
	}

	return value, nil
}

func (s *StampyBucket) deleteValueWithKeyIfPresent(key string) {

	_, ok := s.keyValueCache[key]

	if ok {
		s.cacheMutex.Lock()

		delete(s.keyValueCache, key)
		delete(s.ttlIndex, key)

		s.stampyBucketStats.incrementKeyDeletes()

		log.Println("Deleted key from cache:", key)

		s.cacheMutex.Unlock()
	}
}

func (s *StampyBucket) deleteExpiredKeys() {


	if len(s.ttlIndex) == 0 {
		return
	}

	now := time.Now()

	for i, _ := range s.ttlIndex {
		if s.keyValueCache[i].ExpiryTime.Before(now) {
			s.deleteValueWithKeyIfPresent(i)
			s.stampyBucketStats.incrementExpiredKeys()
		}
	}
}


