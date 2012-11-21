package gokeeper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/skein"
	"hash"
)

type AES struct {
	ctx    cipher.Stream
	engine cipher.Block
}

func (aes *AES) Update(input []byte) []byte {
	buff := append([]byte(nil), input...)
	aes.ctx.XORKeyStream(buff, buff)
	return buff
}

func (aes *AES) BlockSize() int {
	return aes.engine.BlockSize()
}

func NewAES(key, iv []byte) (res *AES, err error) {
	engine, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ctx := cipher.NewCTR(engine, iv)
	return &AES{ctx: ctx, engine: engine}, nil
}

func BlockSizeAES() int {
	return aes.BlockSize
}

func SHA256(input []byte) []byte {
	ctx := sha256.New()
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func NewSkein256() hash.Hash {
	res, err := skein.New(256, 256)
	if err != nil {
		return nil
	}
	return res
}

func NewSkein512() hash.Hash {
	res, err := skein.New(512, 512)
	if err != nil {
		return nil
	}
	return res
}

func NewSkein1024() hash.Hash {
	res, err := skein.New(1024, 1024)
	if err != nil {
		return nil
	}
	return res
}

func Skein256(input []byte) []byte {
	ctx := NewSkein256()
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func Skein512(input []byte) []byte {
	ctx := NewSkein512()
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func Skein1024(input []byte) []byte {
	ctx := NewSkein1024()
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func Rand(size uint) []byte {
	buff := make([]byte, size)
	_, err := rand.Read(buff)
	if err != nil {
		return nil
	}
	return buff
}

func HMAC_SHA256(input, key []byte) []byte {
	ctx := hmac.New(sha256.New, key)
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func HMAC_SHA512(input, key []byte) []byte {
	ctx := hmac.New(sha512.New, key)
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func HMAC_Skein256(input, key []byte) []byte {
	ctx := hmac.New(NewSkein256, key)
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func HMAC_Skein512(input, key []byte) []byte {
	ctx := hmac.New(NewSkein512, key)
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func HMAC_Skein1024(input, key []byte) []byte {
	ctx := hmac.New(NewSkein1024, key)
	_, err := ctx.Write(input)
	if err != nil {
		return nil
	}
	return ctx.Sum([]byte(""))
}

func PBKDF2_SHA1(password, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = Rand(8)
	}
	return salt, pbkdf2.Key(password, salt, 10000, 32, sha1.New)
}

func PBKDF2_Skein256(password, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = Rand(8)
	}
	return salt, pbkdf2.Key(password, salt, 10000, 32, NewSkein256)
}
