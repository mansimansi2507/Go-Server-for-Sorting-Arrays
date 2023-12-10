package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortSingle(arr []int) {
	sort.Ints(arr)
}
func sortConcurrent(arr []int, wg *sync.WaitGroup) {
	defer wg.Done()
	sort.Ints(arr)
}
func processSingle(w http.ResponseWriter, r *http.Request) {
	var request SortRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	for _, arr := range request.ToSort {
		sortSingle(arr)
	}

	endTime := time.Now()

	response := SortResponse{
		SortedArrays: request.ToSort,
		TimeNS:       int64(endTime.Sub(startTime).Nanoseconds()),
	}

	json.NewEncoder(w).Encode(response)
}
func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var request SortRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	for _, arr := range request.ToSort {
		wg.Add(1)
		go sortConcurrent(arr, &wg)
	}
	wg.Wait()

	endTime := time.Now()

	response := SortResponse{
		SortedArrays: request.ToSort,
		TimeNS:       int64(endTime.Sub(startTime).Nanoseconds()),
	}

	json.NewEncoder(w).Encode(response)
}
func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	port := ":8000"
	server := http.Server{
		Addr:    port,
		Handler: nil,
	}
	println("Server listening on", port)
	server.ListenAndServe()
}
