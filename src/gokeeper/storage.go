package gokeeper

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
)

type Storage struct {
	path string
	data map[string][]byte
}

func (s *Storage) Save(masterkey []byte) error {
	var buffer bytes.Buffer
	iv := Rand(BlockSizeAES())
	ctx, err := NewAES(masterkey, iv)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(&buffer)
	err2 := encoder.Encode(s.data)
	if err2 != nil {
		return err2
	}
	file, err3 := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0600)
	if err3 != nil {
		return err3
	}
	_, err4 := file.Write(Skein1024(masterkey))
	if err4 != nil {
		return err4
	}
	_, err5 := file.Write(iv)
	if err5 != nil {
		return err5
	}
	_, err6 := file.Write(ctx.Update(buffer.Bytes()))
	if err6 != nil {
		return err6
	}
	file.Close()
	return nil
}

func (s *Storage) Load(masterkey []byte) error {
	var buffer bytes.Buffer
	file, err := os.OpenFile(s.path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	var masterhash []byte = make([]byte, 128)
	_, err2 := file.Read(masterhash)
	if err2 != nil {
		return err2
	}
	if !bytes.Equal(Skein1024(masterkey), masterhash) {
		return errors.New("Bad Master Key")
	}
	var iv []byte = make([]byte, BlockSizeAES())
	_, err3 := file.Read(iv)
	if err3 != nil {
		return err3
	}
	ciphertext, err4 := readAll(file)
	if err4 != nil {
		return err4
	}
	ctx, err5 := NewAES(masterkey, iv)
	if err5 != nil {
		return err5
	}
	_, err6 := buffer.Write(ctx.Update(ciphertext))
	if err6 != nil {
		return err6
	}
	decoder := gob.NewDecoder(&buffer)
	err7 := decoder.Decode(&s.data)
	if err7 != nil {
		return err7
	}
	file.Close()
	return nil
}

func (s *Storage) Get(key string, masterkey []byte) ([]byte, error) {
	buffer := s.data[key]
	if buffer == nil {
		return nil, errors.New("Unknown key")
	}
	var iv []byte = buffer[:BlockSizeAES()]
	var ciphertext []byte = buffer[BlockSizeAES():]
	ctx, err := NewAES(masterkey, iv)
	if err != nil {
		return nil, err
	}
	return ctx.Update(ciphertext), nil
}

func (s *Storage) Put(key string, data []byte, masterkey []byte) error {
	iv := Rand(BlockSizeAES())
	ctx, err := NewAES(masterkey, iv)
	if err != nil {
		return err
	}
	s.data[key] = append(iv, ctx.Update(data)...)
	return nil
}

func (s *Storage) UpdateKey(oldmasterkey, masterkey []byte) error {
	for key, _ := range s.data {
		uncrypt_data, err := s.Get(key, oldmasterkey)
		if err != nil {
			return err
		}
		err = s.Put(key, uncrypt_data, masterkey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Data() map[string][]byte {
	return s.data
}

func NewStorage(path string, masterkey []byte) (res *Storage, err error) {
	res = &Storage{path: path, data: make(map[string][]byte)}
	err = res.Load(masterkey)
	return res, err
}
