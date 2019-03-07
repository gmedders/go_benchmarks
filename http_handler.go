package functest

import (
	"encoding/json"
	// "log"
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
func startServer(srvErr chan error) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf(r.Body)
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	go func(e chan error) {
		// defer wg.Done()
		// ListenAndServe() returns ErrServerClosed on graceful close
		e <- srv.ListenAndServe()
	}(srvErr)

	return srv
}
