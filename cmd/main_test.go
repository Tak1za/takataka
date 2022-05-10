package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
)

func setupMyHandler() *myHandler {
	shards := 100
	allShards := make([]cache, shards)
	for i := 0; i < shards; i++ {
		allShards[i] = cache{
			holder: make(map[uint32]int),
		}
	}
	return &myHandler{
		totalShards: 100,
		entries:     make([]entry, 0),
		allShards:   allShards,
	}
}

func setupRedisHandler() *redisHandler {
	return &redisHandler{
		ctx: context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func cleanupRedis(rh *redisHandler) {
	rh.rdb.FlushAll(rh.ctx)
}

func BenchmarkAddStringValue(b *testing.B) {
	mh := setupMyHandler()
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		value := fmt.Sprintf("solo%d", i)
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/add", 201)
		}
	}
}

func BenchmarkAddIntValue(b *testing.B) {
	mh := setupMyHandler()
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		value := i
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/add", 201)
		}
	}
}

func BenchmarkAddComplexValue(b *testing.B) {
	mh := setupMyHandler()
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		complex := struct {
			a string
			b string
		}{
			a: fmt.Sprintf("test%d", i),
			b: fmt.Sprintf("test%d", i),
		}
		input := content{
			Key:   key,
			Value: complex,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/add", 201)
		}
	}
}

func BenchmarkGetStringValue(b *testing.B) {
	mh := setupMyHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		value := fmt.Sprintf("solo%d", i)
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/get/", mh.Get)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/get", 200)
		}
	}
}

func BenchmarkGetIntValue(b *testing.B) {
	mh := setupMyHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		value := i
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/get/", mh.Get)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/get", 200)
		}
	}
}

func BenchmarkGetComplexValue(b *testing.B) {
	mh := setupMyHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/add", mh.Add)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		complex := struct {
			a string
			b string
		}{
			a: fmt.Sprintf("test%d", i),
			b: fmt.Sprintf("test%d", i),
		}
		input := content{
			Key:   key,
			Value: complex,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/add", 201)
		}
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/get/", mh.Get)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/get", 200)
		}
	}
}

func BenchmarkRedisAddStringValue(b *testing.B) {
	rh := setupRedisHandler()
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		value := fmt.Sprintf("solo%d", i)
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/redis/add", 201)
		}
	}

	cleanupRedis(rh)
}

func BenchmarkRedisAddIntValue(b *testing.B) {
	rh := setupRedisHandler()
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		value := i
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/redis/add", 201)
		}
	}

	cleanupRedis(rh)
}

func BenchmarkRedisAddComplexValue(b *testing.B) {
	rh := setupRedisHandler()
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		complex := struct {
			a string
			b string
		}{
			a: fmt.Sprintf("test%d", i),
			b: fmt.Sprintf("test%d", i),
		}
		input := content{
			Key:   key,
			Value: complex,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/redis/add", 201)
		}
	}

	cleanupRedis(rh)
}

func BenchmarkRedisGetStringValue(b *testing.B) {
	rh := setupRedisHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		value := fmt.Sprintf("solo%d", i)
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/redis/get/", rh.redisGet)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/redis/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/redis/get/", 200)
		}
	}

	cleanupRedis(rh)
}

func BenchmarkRedisGetIntValue(b *testing.B) {
	rh := setupRedisHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		value := i
		input := content{
			Key:   key,
			Value: value,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/redis/get/", rh.redisGet)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/redis/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/redis/get/", 200)
		}
	}

	cleanupRedis(rh)
}

func BenchmarkRedisGetComplexValue(b *testing.B) {
	rh := setupRedisHandler()
	val := b.N
	router := http.NewServeMux()
	router.HandleFunc("/redis/add", rh.redisAdd)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		complex := struct {
			a string
			b string
		}{
			a: fmt.Sprintf("test%d", i),
			b: fmt.Sprintf("test%d", i),
		}
		input := content{
			Key:   key,
			Value: complex,
		}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(input)
		req, _ := http.NewRequest("POST", "/redis/add", &buf)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 201 {
			b.Errorf("%d in API %s; expected %d\n", w.Code, "/add", 201)
		}
	}

	b.ResetTimer()

	router = http.NewServeMux()
	router.HandleFunc("/redis/get/", rh.redisGet)
	for i := 0; i < val; i++ {
		key := fmt.Sprint(i)
		getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/redis/get/%s", key), nil)
		wg := httptest.NewRecorder()
		router.ServeHTTP(wg, getRequest)
		if wg.Code != 200 {
			b.Errorf("%d in API %s; expected %d\n", wg.Code, "/redis/get/", 200)
		}
	}

	cleanupRedis(rh)
}
