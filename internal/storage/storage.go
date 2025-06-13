package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var DataFileFormat string

func getPWD() (pwd string, err error) {
	pwd, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}

	return
}

func getPWDfile(file string) (pwdFile string, err error) {
	var pwd string
	pwd, err = getPWD()
	if err != nil {
		return
	}

	pwdFile = filepath.Join(pwd, file+DataFileFormat)
	return
}

type StorageMap map[string]string

type Storage struct {
	/*file -> [key -> value]*/
	data  map[string]StorageMap
	mutex sync.RWMutex
}

func NewStorage() (s *Storage, err error) {
	s = &Storage{
		data: make(map[string]StorageMap),
	}

	pwd, err := getPWD()
	if err != nil {
		return
	}

	items, err := os.ReadDir(pwd)
	if err != nil {
		return
	}

	var files []string
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if strings.HasSuffix(item.Name(), DataFileFormat) {
			cutFile, _ := strings.CutSuffix(item.Name(), DataFileFormat)
			files = append(files, cutFile)
		}
	}

	for _, file := range files {
		err = s.AddDataFile(file)
		if err != nil {
			return
		}
	}

	return
}

func (s *Storage) Files() (files []string, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	it := maps.Keys(s.data)
	for file := range it {
		files = append(files, file)
	}

	return

}

func (s *Storage) Keys(file string) (keys []string, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, exists := s.data[file]; !exists {
		err = fmt.Errorf("no such file '%s'", file)
		return
	}

	it := maps.Keys(s.data[file])
	for key := range it {
		keys = append(keys, key)
	}

	return
}

func (s *Storage) AddDataFile(file string) (err error) {
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

	storageMap := make(map[string]string)
	defer func() {
		s.data[file] = storageMap
	}()

	var fileInfo os.FileInfo
	fileInfo, err = saveFile.Stat()
	if err != nil {
		return
	}

	if fileInfo.Size() == 0 {
		return
	}

	var rawData []byte
	rawData, err = os.ReadFile(pwdFile)
	if err != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(rawData))
	err = decoder.Decode(&storageMap)
	if err != nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return
}

func (s *Storage) Store(file string, key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.data[file]; !exists {
		return fmt.Errorf("no such file '%s'", file)
	}

	if len(key) == 0 {
		return fmt.Errorf("empty key")
	}

	s.data[file][key] = value
	return nil
}

func (s *Storage) Load(file string, key string) (value string, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var exists bool

	if _, exists = s.data[file]; !exists {
		err = fmt.Errorf("no such file '%s'", file)
		return
	}

	if value, exists = s.data[file][key]; !exists {
		err = fmt.Errorf("no such key '%s'", key)
		return
	}

	return value, nil
}

func (s *Storage) Save(file string) (err error) {
	if _, exists := s.data[file]; !exists {
		err = fmt.Errorf("no such file '%s'", file)
		return
	}

	var buffer bytes.Buffer

	enconder := gob.NewEncoder(&buffer)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	err = enconder.Encode(s.data[file])
	if err != nil {
		return
	}

	pwdFile, err := getPWDfile(file)
	if err != nil {
		return
	}

	var saveFile *os.File
	saveFile, err = os.OpenFile(pwdFile, os.O_WRONLY|os.O_TRUNC, 0644)
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
