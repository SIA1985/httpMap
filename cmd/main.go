package main

import (
	"flag"
	httpinterface "storage/internal/httpInterface"
	"storage/internal/storage"
)

func main() {
	var err error

	file := flag.String("file", "", "name of save file")
	addr := flag.String("addr", "", "http-server addr")

	flag.Parse()

	httpinterface.Storage, err = storage.NewStorage(*file)
	if err != nil {
		panic(err)
	}

	err = httpinterface.Listen(*addr)
	if err != nil {
		panic(err)
	}
}
