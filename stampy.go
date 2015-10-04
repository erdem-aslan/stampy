package main

import (
	"net/http"
	"log"
	"flag"
	"fmt"
	"runtime"
	"time"
	"hash/fnv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const stampyPath string = "/stampy"
const versionPath string = stampyPath + "/v1"
const infoPath string = versionPath + "/info"
const cachePath string = versionPath + "/cache"
const cacheRootPath string = cachePath + "/"
const memoryDivider float32 = 1024

var (
	stampyInfo StampyInfo
	m runtime.MemStats
	buckets map[int]StampyBucket
	config StampyConfig
)


func main() {

	log.SetPrefix("<Stampy> " + log.Prefix())



	var ipFlag string
	var portFlag int
	var bucketsCount int
	var configFilePath string

	flag.StringVar(&ipFlag, "ip", DefaultIp, "A valid IPv4 address for serving restful interface, ex: 127.0.0.1")
	flag.IntVar(&portFlag, "port", DefaultPort, "An unoccupied port for serving restful interface")
	flag.IntVar(&bucketsCount, "buckets", DefaultBucketCount, "Number of buckets for keys to be evenly distributed, " +
	"higher numbers will increase concurrency with additional memory overhead")

	flag.StringVar(&configFilePath, "configFile", "", "All options are also configurable via config file in YAML format.")

	flag.Parse()

	if configFilePath != "" {
		config = loadConfiguration(configFilePath)
	} else {
		config = StampyConfig{bucketsCount, ipFlag, portFlag}
	}

	initializeBuckets(config.Buckets)
	resolveStampyInformation()

	resourceHandler := NewResourceHandler()

	registerStampyHandlers(resourceHandler)


	log.Println("Stampy starting on", ipFlag, "with port", portFlag)
	log.Fatal(http.ListenAndServe(ipFlag + ":" + fmt.Sprint(portFlag), nil))

}

func loadConfiguration(configFilePath string) StampyConfig {

	file, err := ioutil.ReadFile(configFilePath)


	if err != nil {
		log.Fatalln(err)

	}

	var config StampyConfig

	yamlError := yaml.Unmarshal(file, &config)

	if yamlError != nil {
		log.Fatalln(yamlError)
	}

	log.Printf("Configuration: %v\n", config)

	return config
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
	stampyInfo.StampyBucketCount = config.Buckets

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

	if memory > memoryDivider {
		memory = memory / float32(memoryDivider)
		memoryUnit = "mb"

		if memory > memoryDivider {
			memory = memory / float32(memoryDivider)
			memoryUnit = "gb"
		}
	}
	return
}

func getBucket(key string) *StampyBucket {
	hash := fnv.New32()
	hash.Write([]byte(key))
	s := buckets[int(hash.Sum32()) % int(config.Buckets)]
	return &s
}


func registerStampyHandlers(resourceHandler *ResourceHandler) {

	http.HandleFunc(infoPath, resourceHandler.infoHandler)
	http.HandleFunc(cachePath, resourceHandler.statsHandler)
	http.HandleFunc(cacheRootPath, resourceHandler.cacheHandler)

}








