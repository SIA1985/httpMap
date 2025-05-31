package main

import (
	"flag"
	httpinterface "storage/internal/httpInterface"
	"storage/internal/storage"
	"time"
)

func main() {
	var err error

	pathToDataFile := flag.String("saveFile", "", "path to save file")
	addr := flag.String("addr", "", "http-server addr")
	savePeriod := flag.Uint("savePeriod", 0, "secs to drop ram to file")

	httpinterface.Storage, err = storage.NewStorage(*pathToDataFile)
	if err != nil {
		panic(err)
	}

	go func() {
		for range time.Tick(time.Duration(*savePeriod) * time.Second) {
			httpinterface.Storage.Save()
		}
	}()

	err = httpinterface.Listen(*addr)
	if err != nil {
		panic(err)
	}
}
