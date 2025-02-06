package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) HandlerMetrics(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html")
	rw.WriteHeader(200)
	text := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, cfg.FileserverHits.Load())
	rw.Write([]byte(text))
}
