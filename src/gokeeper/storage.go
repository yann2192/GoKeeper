package gokeeper

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
)

type Storage struct {
	path string
	Data map[string][]byte
}

func (s *Storage) Save() error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err2 := encoder.Encode(s.Data)
	if err2 != nil {
		return err2
	}
	file.Close()
	return nil
}

func (s *Storage) Load() error {
	file, err := os.OpenFile(s.path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(file)
	err2 := decoder.Decode(&s.Data)
	if err2 != nil {
		return err2
	}
	file.Close()
	return nil
}

func (s *Storage) Get(key string) ([]byte, error) {
	buffer := s.Data[key]
	if buffer == nil {
		return nil, errors.New("Unknown key")
	}
	var iv []byte = buffer[:BlockSizeAES()]
	var ciphertext []byte = buffer[BlockSizeAES():]
	ctx, err := NewAES(KEY, iv)
	if err != nil {
		return nil, err
	}
	return ctx.Update(ciphertext), nil
}

func (s *Storage) Put(key string, data []byte) error {
	var result bytes.Buffer
	iv := Rand(uint(BlockSizeAES()))
	result.Write(iv)
	ctx, err := NewAES(KEY, iv)
	if err != nil {
		return err
	}
	result.Write(ctx.Update(data))
	s.Data[key] = result.Bytes()
	return nil
}

func NewStorage(path string) (res *Storage) {
	res = &Storage{path: path, Data: make(map[string][]byte)}
	res.Load()
	return res
}
