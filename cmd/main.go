package main

import (
	"flag"
	httpinterface "storage/internal/httpInterface"
	"storage/internal/storage"
)

func main() {
	var err error

	format := flag.String(".format", ".map", "format of save file")
	addr := flag.String("addr", "", "http-server addr")

	flag.Parse()

	storage.DataFileFormat = *format

	httpinterface.Storage, err = storage.NewStorage()
	if err != nil {
		panic(err)
	}

	err = httpinterface.Listen(*addr)
	if err != nil {
		panic(err)
	}
}
