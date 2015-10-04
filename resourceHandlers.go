package main
import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"strings"
)



type ResourceHandler struct {
	infoHandler  http.HandlerFunc
	statsHandler http.HandlerFunc
	cacheHandler http.HandlerFunc
}

func NewResourceHandler() *ResourceHandler {
	return &ResourceHandler{infoFunction, statsFunction, cacheFunction}
}

var infoFunction = func(w http.ResponseWriter, r *http.Request) {

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

}

var statsFunction = func(w http.ResponseWriter, r *http.Request) {

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
}

var cacheFunction = func(w http.ResponseWriter, r *http.Request) {
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

		var p StampyCacheRequest

		err := decoder.Decode(&p)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println(p.Value)

		if p.Value == "" {
			log.Println("Value cannot be null, p:", p)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		getBucket(key).processKeyViaRequest(key, &p)
		w.WriteHeader(http.StatusOK)

	case "DELETE":
		getBucket(key).deleteValueWithKeyIfPresent(key)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

}
