package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

type Storage struct {
	data           map[string][]byte
	pathToDataFile string

	mutex sync.Mutex
}

func NewStorage(pathToDataFile string) (s *Storage, err error) {
	var rawData []byte
	s = &Storage{pathToDataFile: pathToDataFile}

	rawData, err = os.ReadFile(pathToDataFile)
	if err != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(rawData))
	err = decoder.Decode(&s.data)
	if err != nil {
		return
	}

	return
}

func (s *Storage) Save() (err error) {
	var buffer bytes.Buffer

	enconder := gob.NewEncoder(&buffer)

	s.mutex.Lock()

	err = enconder.Encode(s.data)
	if err != nil {
		return
	}

	s.mutex.Unlock()

	/*todo: сохранение добавление*/
	var saveFile *os.File
	saveFile, err = os.OpenFile(s.pathToDataFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}

	data := buffer.Bytes()

	var n int
	n, err = saveFile.Write(data)
	if n != len(data) {
		err = fmt.Errorf("saved only %d bytes", n)
		return
	}
	if err != nil {
		return
	}

	return
}
