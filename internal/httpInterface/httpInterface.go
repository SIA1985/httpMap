package httpinterface

import (
	"net/http"
	"storage/internal/storage"
)

var Storage *storage.Storage

func Listen(addr string) error {

	http.HandleFunc("GET /storage/{key}", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("PUT /storage/{key}/{data}", func(w http.ResponseWriter, r *http.Request) {

	})

	return http.ListenAndServe(addr, nil)
}
