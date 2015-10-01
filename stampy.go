package main

import (
	"net/http"
	"log"
	"flag"
	"fmt"
	"runtime"
	"encoding/json"
	"time"
	"strings"
	"hash/fnv"
	"bytes"
)

const versionPath string = "/v1"
const infoPath string = versionPath + "/info"
const cachePath string = versionPath + "/cache"
const cacheRootPath string = cachePath + "/"
const memoryDivider float32 = 1024

var (
	stampyInfo StampyInfo
	m runtime.MemStats
	buckets map[int]StampyBucket
	bucketsCount int
)


func main() {

	log.SetPrefix("<Stampy> " + log.Prefix())


	var ipFlag string
	var portFlag int

	flag.StringVar(&ipFlag, "ip", DefaultIp, "A valid IPv4 address for serving restful interface, ex: 127.0.0.1")
	flag.IntVar(&portFlag, "port", DefaultPort, "An unoccupied port for serving restful interface")
	flag.IntVar(&bucketsCount, "buckets", DefaultBucketCount, "Number of buckets for keys to be evenly distributed," +
	"higher numbers will increase concurrency with a memory overhead")

	flag.Parse()

	initializeBuckets(bucketsCount)
	resolveStampyInformation()
	registerStampyHandlers()


	log.Println("Stampy starting on", ipFlag, "with port", portFlag)
	log.Fatal(http.ListenAndServe(ipFlag + ":" + fmt.Sprint(portFlag), nil))

}

func initializeBuckets(bucketCount int) {

	buckets = make(map[int]StampyBucket, bucketCount)

	for i := 0; i < bucketCount; i++ {

		var bucket StampyBucket

		bucket.keyValueCache = make(map[string]StampyBucketEntry)
		bucket.ttlIndex = make(map[string]bool)
		bucket.stampyBucketStats = &StampyBucketStats{0, 0, 0, 0, 0, 0}
		bucket.bucketIndex = i

		ttlTimer := time.NewTicker(time.Minute)

		go func() {
			for {
				<-ttlTimer.C
				bucket.deleteExpiredKeys()
			}
		}()

		buckets[i] = bucket
	}

}

func resolveStampyInformation() {

	stampyInfo.Name = "Stampy, Elephant in the room"
	stampyInfo.Version = Version
	stampyInfo.Os = fmt.Sprint(runtime.GOOS, "-", runtime.GOARCH)
	stampyInfo.CpuCores = runtime.NumCPU()
	stampyInfo.StampyBucketCount = bucketsCount

	memory, memoryUnit := resolveStampyMemoryUsage()

	stampyInfo.MemoryUsage = fmt.Sprintf("%.2f %s", memory, memoryUnit)

	stampyInfo.Started = time.Now()

	ticker := time.NewTicker(time.Second * 10)

	go func() {

		for {
			<-ticker.C
			memory, memoryUnit := resolveStampyMemoryUsage()
			stampyInfo.MemoryUsage = fmt.Sprintf("%.2f %s", memory, memoryUnit)
		}
	}()
}

func resolveStampyMemoryUsage() (memory float32, memoryUnit string) {

	runtime.ReadMemStats(&m)

	memory = float32(m.Alloc) / memoryDivider
	memoryUnit = "kb"

	if memory > 1024 {
		memory = memory / float32(1024)
		memoryUnit = "mb"

		if memory > 1024 {
			memory = memory / float32(1024)
			memoryUnit = "gb"
		}
	}
	return
}

func getBucket(key string) *StampyBucket {
	hash := fnv.New32()
	hash.Write([]byte(key))
	s := buckets[int(hash.Sum32()) % int(bucketsCount)]
	return &s
}


/**
	Registers handlers for REST interface of Stampy the Mighty Elephant
 */
func registerStampyHandlers() {


	http.HandleFunc(infoPath, func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case "GET":
			payload, err := json.Marshal(stampyInfo)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var indented bytes.Buffer
			json.Indent(&indented, payload, "", "\t")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(indented.Bytes())

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc(cachePath, func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case "GET":
			var totalStats StampyBucketStats

			for _, v := range buckets {
				totalStats.KeyHits += v.stampyBucketStats.KeyHits
				totalStats.AbsentKeyHits += v.stampyBucketStats.AbsentKeyHits
				totalStats.ExpiredKeyHits += v.stampyBucketStats.ExpiredKeyHits
				totalStats.KeyPuts += v.stampyBucketStats.KeyPuts
				totalStats.KeyDeletes += v.stampyBucketStats.KeyDeletes
				totalStats.ExpiredKeys += v.stampyBucketStats.ExpiredKeys
			}

			payload, err := json.Marshal(totalStats)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var indented bytes.Buffer
			json.Indent(&indented, payload, "", "\t")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(indented.Bytes())

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc(cacheRootPath, func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, cacheRootPath)

		switch r.Method {

		case "GET":
			entry, cacheError := getBucket(key).getValueWithKey(key)

			if cacheError != nil {
				// key is missing
				w.WriteHeader(http.StatusNotFound)
				return
			}

			payload, err := json.Marshal(entry)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(payload)

		case "PUT":

			decoder := json.NewDecoder(r.Body)

			var p StampyPayload

			err := decoder.Decode(&p)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Println(p.Payload)

			if p.Payload == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			getBucket(key).putKeyWithValue(key, p.Payload, p.ValidUntil)
			w.WriteHeader(http.StatusOK)

		case "DELETE":
			getBucket(key).deleteValueWithKeyIfPresent(key)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}

	})

}

type StampyPayload struct {
	Payload    string `json:"payload"`
	ValidUntil time.Time `json:"validUntil"`
}







