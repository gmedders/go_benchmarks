package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "sync"
)

func incrementInJSON(payload []byte) []byte {
	var result map[string]int

	json.Unmarshal(payload, &result)
	result["counter"]++
	byteJSON, _ := json.Marshal(result)

	return byteJSON
}

// from https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
func StartServer(addr string, srvErr chan error) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In the server handler")

		var result map[string]int
		json.NewDecoder(r.Body).Decode(&result)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("400 - Unable to deserial request body"))
		// }
		result["counter"]++

		resultJSON, _ := json.Marshal(result)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte("500 - Failed to serialize response"))
		// }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	})

	srv := &http.Server{Addr: addr, Handler: mux}

	go func(e chan error) {
		// defer wg.Done()
		// ListenAndServe() returns ErrServerClosed on graceful close
		e <- srv.ListenAndServe()
	}(srvErr)

	return srv
}
