package main

import (
	"benchmarks/utils"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	// "sync"
	// "utils"
)

func main() {
	srvErr := make(chan error)
	srv := utils.StartServer(":8080", srvErr)

	message := map[string]int{
		"counter": 0,
	}

	bytesRepresentation, _ := json.Marshal(message)

	_, err := http.Post("http://localhost:8080/count", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalf("Error on posting: %v", err)
	}
	time.Sleep(2 * time.Second)
	// return

	srv.Shutdown(context.Background())

	// Capture the error returned by ListenAndServe when we close the server
	<-srvErr
}
