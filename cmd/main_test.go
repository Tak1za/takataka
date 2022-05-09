package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() *myHandler {
	return &myHandler{
		entries: make([][]byte, 0),
		cache:   make(map[int]int),
	}
}

func BenchmarkAddStringValue(b *testing.B) {
	mh := setup()
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	mh := setup()
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	mh := setup()
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	mh := setup()
	val := b.N
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	router.Handle("/get/", mh)
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
	mh := setup()
	val := b.N
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	router.Handle("/get/", mh)
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
	mh := setup()
	val := b.N
	router := http.NewServeMux()
	router.Handle("/add", mh)
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
	router.Handle("/get/", mh)
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