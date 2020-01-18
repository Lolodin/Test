package Encode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

//Шифруем текст, используя время в качестве ключа
//return key, encodetext
func EncodeAes(plaintext []byte) ([]byte, []byte) {
	r := time.Now().String()
	randomsalt := r[20:]
	fmt.Println(randomsalt)
	randomsaltByte := []byte(randomsalt)
	var key = make([]byte, hex.EncodedLen(len(randomsalt)))
	hex.Encode(key, randomsaltByte)
	fmt.Println("key :", string(key[:16]))
	block, err := aes.NewCipher(key[:16])
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return key[:16], ciphertext

}

//return decode text
func DecodeAes(key, ciphertext []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}


	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}
