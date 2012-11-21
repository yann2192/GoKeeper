package gokeeper

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
)

type Storage struct {
	path string
	blob struct {
		Data map[string][]byte
		Hash []byte
	}
}

func (s *Storage) Save() error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err2 := encoder.Encode(s.blob)
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
	err2 := decoder.Decode(&s.blob)
	if err2 != nil {
		return err2
	}
	file.Close()
	return nil
}

func (s *Storage) Get(key string) ([]byte, error) {
	buffer := s.blob.Data[key]
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
	s.blob.Data[key] = result.Bytes()
	return nil
}

func (s *Storage) Validate(key []byte) bool {
	if s.blob.Hash == nil {
		s.blob.Hash = Skein1024(key)
        return true
	} else {
		if bytes.Equal(Skein1024(key), s.blob.Hash) {
			return true
		}
	}
	return false
}

func (s *Storage) Data() map[string][]byte {
	return s.blob.Data
}

func NewStorage(path string) (res *Storage) {
	//res = &Storage{path: path, blob:*NewBlob()}
	res = &Storage{path: path}
	res.blob.Data = make(map[string][]byte)
	res.Load()
	return res
}
