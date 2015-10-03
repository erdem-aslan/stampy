package main
import (
	"fmt"
	"time"
)

type StampyCacheRequest struct {
	Value      string `json:"value"`
	TimeToLive int64 `json:"timeToLive"`
	expiryDate time.Time
}

func (r *StampyCacheRequest) String() string {
	return fmt.Sprintf("%v", r)
}

func (s *StampyCacheRequest) getExpiryDate() time.Time {

	if s.expiryDate.IsZero() {
		s.expiryDate = time.Now().Add(time.Duration(s.TimeToLive) * time.Second)
	}

	return s.expiryDate
}

type StampyBucketEntry struct {
	EntryValue       string `json:"value"`
	CreationTime     time.Time `json:"creationDate"`
	LastAccessedTime time.Time `json:"lastAccessed"`
	ExpiryTime       time.Time `json:"expiryTime"`
}

type StampyBucketStats struct {
	KeyPuts        uint64 `json:"keyPuts"`
	KeyDeletes     uint64 `json:"keyDeletes"`
	KeyHits        uint64 `json:"keyHits"`
	AbsentKeyHits  uint64 `json:"absentKeyHits"`
	ExpiredKeys    uint64 `json:"expiredKeys"`
	ExpiredKeyHits uint64 `json:"expiredKeyHits"`
}

func (stats *StampyBucketStats) incrementKeyPuts() {
	stats.KeyPuts += 1
}

func (stats *StampyBucketStats) incrementKeyDeletes() {
	stats.KeyDeletes += 1
}
func (stats *StampyBucketStats) incrementKeyHits() {
	stats.KeyHits += 1
}
func (stats *StampyBucketStats) incrementAbsentKeyHits() {
	stats.AbsentKeyHits += 1
}
func (stats *StampyBucketStats) incrementExpiredKeys() {
	stats.ExpiredKeys += 1
}
func (stats *StampyBucketStats) incrementExpiredKeyHits() {
	stats.ExpiredKeyHits += 1
}


