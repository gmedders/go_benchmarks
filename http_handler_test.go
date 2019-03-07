package functest

import (
	"context"
	"encoding/json"
	"net/http"
	// "fmt"
	// "sync"
	"testing"
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
	srv := startServer(srvErr)

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

// func BenchmarkIncrementIntByValue(b *testing.B) {
// 	// set up http handler
// 	b.ResetTimer()
//
// 	counter := 0
// 	for i := 0; i < b.N; i++ {
// 		counter = increment(counter)
// 	}
// }