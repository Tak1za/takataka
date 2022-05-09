package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type content struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type myHandler struct {
	l       sync.Mutex
	entries [][]byte
	cache   map[int]int
}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var c content
		if err := json.Unmarshal(body, &c); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		idx := mh.appendAndGetIndex(body)

		mh.cache[len(mh.entries)-1] = idx
		w.WriteHeader(http.StatusCreated)
	} else {
		key := r.URL.Path[5:]
		intKey, _ := strconv.Atoi(key)
		offset := mh.cache[intKey]
		fetchedByteValue := mh.entries[offset]
		w.Header().Add("Content-Type", "application/json")
		w.Write(fetchedByteValue)
	}

}

func main() {
	mux := http.NewServeMux()
	mh := myHandler{
		entries: make([][]byte, 0),
		cache:   make(map[int]int),
	}
	mux.Handle("/add", &mh)
	mux.Handle("/get/", &mh)

	http.ListenAndServe(":3000", mux)
}

func (mh *myHandler) appendAndGetIndex(val []byte) int {
	mh.l.TryLock()

	defer mh.l.Unlock()

	mh.entries = append(mh.entries, val)
	index := len(mh.entries) - 1

	return index
}
