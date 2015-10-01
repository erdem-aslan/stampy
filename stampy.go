package main

import (
	"net/http"
	"log"
	"flag"
	"fmt"
	"runtime"
	"encoding/json"
	"time"
)

const versionPath string = "/v1"
const infoPath string = versionPath + "/info"
const cachePath string = versionPath + "/cache"
const cacheRootPath string = cacheRootPath + "/"
const memoryDivider float32 = 1024

var (
	stampyInfo StampyInfo
	m runtime.MemStats
	cache StampyCache
)


func main() {

	log.SetPrefix("<Stampy> " + log.Prefix())

	initializeStampyCache()
	resolveStampyInformation()
	registerStampyHandlers()

	var ipFlag string
	var portFlag int

	flag.StringVar(&ipFlag, "ip", DefaultIp, "A valid IPv4 address for serving restful interface, ex: 127.0.0.1")
	flag.IntVar(&portFlag, "port", DefaultPort, "An unoccupied port for serving restful interface")

	flag.Parse()

	log.Println("Stampy starting on", ipFlag, "with port", portFlag)
	log.Fatal(http.ListenAndServe(ipFlag + ":" + fmt.Sprint(portFlag), nil))

}

func initializeStampyCache() {
	cache.keyValueCache = make(map[string]StampyCacheEntry, 10000)
}

func resolveStampyInformation() {

	stampyInfo.Name = "Stampy, Elephant in the room"
	stampyInfo.Version = Version
	stampyInfo.Os = fmt.Sprint(runtime.GOOS, "-", runtime.GOARCH)
	stampyInfo.GoVersion = runtime.Version()
	stampyInfo.CpuCores = runtime.NumCPU()

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


/**
	Registers handlers for REST interface of Stampy the Mighty Elephant
 */
func registerStampyHandlers() {


	http.HandleFunc(infoPath, func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case "GET":
			json, err := json.Marshal(stampyInfo)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(json)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc(cachePath, func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case "GET":
			json, err := json.Marshal(cache.stampyCacheStats)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(json)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc(cacheRootPath, func(w http.ResponseWriter, r *http.Request) {

	})

}







