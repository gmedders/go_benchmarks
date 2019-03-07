package utils

import (
	"bytes"
	"context"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestIncrementIntInJSON(t *testing.T) {
	message := map[string]int{
		"counter": 0,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Error marshling to bytesRepresentation\n")
	}

	countTo := 50000
	for i := 0; i < countTo; i++ {
		bytesRepresentation = incrementInJSON(bytesRepresentation)
	}
	json.Unmarshal(bytesRepresentation, &message)
	if message["counter"] != countTo {
		t.Errorf("Error in incrementing using JSON: %d != %d\n", message["counter"], countTo)
	}
}

func BenchmarkIncrementIntInJSON(b *testing.B) {
	message := map[string]int{
		"counter": 0,
	}
	bytesRepresentation, _ := json.Marshal(message)

	b.ResetTimer()
	countTo := b.N
	for i := 0; i < countTo; i++ {
		bytesRepresentation = incrementInJSON(bytesRepresentation)
	}
}

func TestGracefulCloseHTTPServer(t *testing.T) {
	srvErr := make(chan error)
	srv := StartServer(":8081", srvErr)

	err := srv.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Error shutting down server: %s\n", err)
	}

	// Capture the error returned by ListenAndServe when we close the server
	listenAndServerError := <-srvErr
	if listenAndServerError != http.ErrServerClosed {
		t.Errorf("Unexpected error when gracefully shutting down server: %s\n", err)
	}
}

func TestIncrementIntViaHTTP(t *testing.T) {
	srvErr := make(chan error)
	srv := StartServer(":8080", srvErr)

	message := map[string]int{
		"counter": 0,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Error marshling to bytesRepresentation\n")
	}

	countTo := 50000
	for i := 0; i < countTo; i++ {
		resp, _ := http.Post("http://localhost:8080/count", "application/json", bytes.NewBuffer(bytesRepresentation))
		json.NewDecoder(resp.Body).Decode(&message)
		resp.Body.Close()
		bytesRepresentation, _ = json.Marshal(message)
	}
	if message["counter"] != countTo {
		t.Errorf("Error in incrementing using JSON: %d != %d\n", message["counter"], countTo)
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Error shutting down server: %s\n", err)
	}

	listenAndServerError := <-srvErr
	if listenAndServerError != http.ErrServerClosed {
		t.Errorf("Unexpected error when gracefully shutting down server: %s\n", err)
	}
}

func TestNoOpHTTP(t *testing.T) {
	srvErr := make(chan error)
	srv := StartServer(":8082", srvErr)

	bytesRepresentation, err := json.Marshal("")
	if err != nil {
		t.Errorf("Error marshling to bytesRepresentation\n")
	}

	resp, err := http.Post("http://localhost:8082/noop", "application/json", bytes.NewBuffer(bytesRepresentation))
	resp.Body.Close()
	if err != nil {
		t.Errorf("Failed to post to noop\n")
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Error shutting down server: %s\n", err)
	}

	listenAndServerError := <-srvErr
	if listenAndServerError != http.ErrServerClosed {
		t.Errorf("Unexpected error when gracefully shutting down server: %s\n", err)
	}
}

func BenchmarkIncrementIntViaHTTP(b *testing.B) {
	srvErr := make(chan error)
	srv := StartServer(":8080", srvErr)

	message := map[string]int{
		"counter": 0,
	}

	bytesRepresentation, _ := json.Marshal(message)
	for i := 0; i < b.N; i++ {
		resp, _ := http.Post("http://localhost:8080/count", "application/json", bytes.NewBuffer(bytesRepresentation))
		json.NewDecoder(resp.Body).Decode(&message)
		bytesRepresentation, _ = json.Marshal(message)
	}

	_ = srv.Shutdown(context.Background())
	<-srvErr
}

func BenchmarkIncrementIntViaPersistentHTTP(b *testing.B) {
	srvErr := make(chan error)
	srv := StartServer(":8080", srvErr)

	message := map[string]int{
		"counter": 0,
	}

	client := &http.Client{Timeout: time.Second * 10}
	bytesRepresentation, _ := json.Marshal(message)
	for i := 0; i < b.N; i++ {
		resp, _ := client.Post("http://localhost:8080/count", "application/json", bytes.NewBuffer(bytesRepresentation))
		json.NewDecoder(resp.Body).Decode(&message)
		resp.Body.Close()
		bytesRepresentation, _ = json.Marshal(message)
	}

	_ = srv.Shutdown(context.Background())
	<-srvErr
}

func BenchmarkNoOpHTTP(b *testing.B) {
	srvErr := make(chan error)
	srv := StartServer(":8080", srvErr)

	bytesRepresentation, _ := json.Marshal("")

	client := &http.Client{Timeout: time.Second * 10}
	// bytesRepresentation, _ := json.Marshal(message)
	for i := 0; i < b.N; i++ {
		resp, _ := client.Post("http://localhost:8080/noop", "application/json", bytes.NewBuffer(bytesRepresentation))
		ioutil.ReadAll(resp.Body)
	}

	_ = srv.Shutdown(context.Background())
	<-srvErr
}

// func BenchmarkIncrementIntByValue(b *testing.B) {
// 	// set up http handler
// 	b.ResetTimer()
//
// 	counter := 0
// 	for i := 0; i < b.N; i++ {
// 		counter = increment(counter)
// 	}
// }
