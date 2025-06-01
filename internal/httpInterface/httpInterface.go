package httpinterface

import (
	"fmt"
	"net/http"
	"storage/internal/storage"
)

var Storage *storage.Storage

func Listen(addr string) error {

	http.HandleFunc("GET /storage/{key}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		key := r.PathValue("key")

		var value []byte
		value, err = Storage.Load(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)
	})

	http.HandleFunc("PUT /storage/{key}/{data}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		key := r.PathValue("key")
		value := r.PathValue("data")

		err = Storage.Store(key, []byte(value))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(addr, nil)
}
