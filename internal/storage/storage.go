package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"
)

func getPWDfile(file string) (pwdFile string, err error) {
	var pwd string
	pwd, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	pwdFile = filepath.Join(pwd, file)
	return
}

type Storage struct {
	data    map[string][]byte
	pwdFile string

	mutex sync.RWMutex
}

func NewStorage(file string) (s *Storage, err error) {
	s = &Storage{
		data: make(map[string][]byte),
	}

	err = s.SetDataFile(file)
	if err != nil {
		return
	}

	return
}

func (s *Storage) Keys() (keys []string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	it := maps.Keys(s.data)
	for key := range it {
		keys = append(keys, key)
	}

	return
}

func (s *Storage) SetDataFile(file string) (err error) {
	var pwdFile string
	pwdFile, err = getPWDfile(file)
	if err != nil {
		return
	}

	var saveFile *os.File
	saveFile, err = os.OpenFile(pwdFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer saveFile.Close()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = make(map[string][]byte)
	s.pwdFile = pwdFile

	var fileInfo os.FileInfo
	fileInfo, err = saveFile.Stat()
	if err != nil {
		return
	}

	if fileInfo.Size() == 0 {
		return
	}

	var rawData []byte
	rawData, err = os.ReadFile(s.pwdFile)
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
	if !exists {
		return nil, fmt.Errorf("no such key '%s'", key)
	}

	return value, nil
}

func (s *Storage) Save() (err error) {
	var buffer bytes.Buffer

	enconder := gob.NewEncoder(&buffer)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	err = enconder.Encode(s.data)
	if err != nil {
		return
	}

	var saveFile *os.File
	saveFile, err = os.OpenFile(s.pwdFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer saveFile.Close()

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
