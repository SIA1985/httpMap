package main

import (
	"flag"
	httpinterface "storage/internal/httpInterface"
	"storage/internal/storage"
)

func main() {
	var err error

	pathToDataFile := flag.String("saveFile", "", "path to save file")
	addr := flag.String("addr", "", "http-server addr")

	httpinterface.Storage, err = storage.NewStorage(*pathToDataFile)
	if err != nil {
		panic(err)
	}

	err = httpinterface.Listen(*addr)
	if err != nil {
		panic(err)
	}
}
