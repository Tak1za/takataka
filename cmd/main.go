package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/fasthash/fnv1a"
)

type content struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type cache struct {
	l      sync.Mutex
	holder map[uint32]int
}

type entry struct {
	value     []byte
	createdAt time.Time
}

type myHandler struct {
	l           sync.Mutex
	totalShards int
	allShards   []cache
	entries     []entry
}

type redisHandler struct {
	rdb *redis.Client
	ctx context.Context
}

func (mh *myHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[5:]
	hashedKey := fnv1a.HashString32(key)
	shardNumber := hashedKey % uint32(mh.totalShards)
	fetchedOffset := mh.allShards[shardNumber].holder[hashedKey]
	w.Write(mh.entries[fetchedOffset].value)
}

func (mh *myHandler) Add(w http.ResponseWriter, r *http.Request) {
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

	hashedIndex := fnv1a.HashString32(c.Key)

	//add to queue
	byteData, _ := json.Marshal(c.Value)
	offset := mh.appendAndGetOffset(byteData)

	//add to cache
	mh.appendToCache(hashedIndex, offset)

	w.WriteHeader(http.StatusCreated)
}

func (mh *myHandler) appendToCache(hashedIndex uint32, offset int) {
	shardNumber := hashedIndex % uint32(mh.totalShards)
	mh.allShards[shardNumber].l.Lock()
	defer mh.allShards[shardNumber].l.Unlock()
	mh.allShards[shardNumber].holder[hashedIndex] = offset
}

func main() {
	shards := 100
	allShards := make([]cache, shards)
	for i := 0; i < shards; i++ {
		allShards[i] = cache{
			holder: make(map[uint32]int),
		}
	}
	mux := http.NewServeMux()
	mh := myHandler{
		allShards:   allShards,
		totalShards: shards,
		entries:     make([]entry, 0),
	}
	rh := redisHandler{
		ctx: context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
	mh.entries = append(mh.entries, entry{value: nil})
	mux.HandleFunc("/add", mh.Add)
	mux.HandleFunc("/get/", mh.Get)
	mux.HandleFunc("/redis/add", rh.redisAdd)
	mux.HandleFunc("/redis/get/", rh.redisGet)

	http.ListenAndServe(":3000", mux)
}

func (mh *myHandler) appendAndGetOffset(val []byte) int {
	mh.l.TryLock()

	defer mh.l.Unlock()

	mh.entries = append(mh.entries, entry{
		value:     val,
		createdAt: time.Now(),
	})
	index := len(mh.entries) - 1

	return index
}

func (rh *redisHandler) redisGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[11:]
	val, err := rh.rdb.Get(rh.ctx, key).Result()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(val))
}

func (rh *redisHandler) redisAdd(w http.ResponseWriter, r *http.Request) {
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

	data, _ := json.Marshal(c.Value)

	if err := rh.rdb.Set(rh.ctx, c.Key, data, time.Duration(1<<63-1)).Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
