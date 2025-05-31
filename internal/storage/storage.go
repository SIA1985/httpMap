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

	mutex sync.RWMutex
}

func NewStorage(pathToDataFile string) (s *Storage, err error) {
	var rawData []byte
	s = &Storage{
		pathToDataFile: pathToDataFile,
		data:           make(map[string][]byte),
	}

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

func (s *Storage) Store(key string, value []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(key) == 0 {
		return fmt.Errorf("empty key")
	}

	s.data[key] = value
	return nil
}

func (s *Storage) Load(key string) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.data[key]
	if exists {
		return nil, fmt.Errorf("no such key '%s'", key)
	}

	return value, nil
}

func (s *Storage) Save() (err error) {
	var buffer bytes.Buffer

	enconder := gob.NewEncoder(&buffer)

	s.mutex.RLock()

	err = enconder.Encode(s.data)
	if err != nil {
		return
	}

	s.mutex.RUnlock()

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
