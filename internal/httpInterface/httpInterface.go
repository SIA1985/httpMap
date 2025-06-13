package httpinterface

import (
	"fmt"
	"net/http"
	"storage/internal/storage"
)

var Storage *storage.Storage

func Listen(addr string) error {

	http.HandleFunc("GET /files", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		files, err := Storage.Files()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, files)
	})

	http.HandleFunc("GET /keys/{file}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		file := r.PathValue("file")

		keys, err := Storage.Keys(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, keys)
	})

	http.HandleFunc("PUT /storage/{file}/{key}/{data}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		var err error

		file := r.PathValue("file")
		key := r.PathValue("key")
		value := r.PathValue("data")

		err = Storage.Store(file, key, value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = Storage.Save(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /storage/{file}/{key}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		file := r.PathValue("file")
		key := r.PathValue("key")

		var value string
		value, err := Storage.Load(file, key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)
	})

	return http.ListenAndServe(addr, nil)
}
