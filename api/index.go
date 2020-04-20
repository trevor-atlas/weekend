package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/skip2/go-qrcode"

	"github.com/trevor-atlas/weekend/api/router"
	"github.com/trevor-atlas/weekend/api/utils"
)

// Handler is the function that Now calls for every request
// Please note that only one function should be exported in this file!
func Handler(w http.ResponseWriter, r *http.Request) {
	defer utils.Track(utils.Runtime("handler"))
	log.Printf("handling request %s: %s", r.Method, r.URL.Path)

	router.NewStupidRouter("/api").
		GET("/encode", QREncode).
		GET("/encrypt", Encrypt).
		GET("/encrypt-with-qr", EncryptWithQR).
		GET("/decrypt", Decrypt).
		DEFAULT(func (w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "I can't do that.")
		}).
		Start(w, r)

	log.Println("finished")
}

func GetSize(r *http.Request) (int, error) {
	defaultSize := 256
	maxSize := 1024
	querySize := r.URL.Query().Get("size")
	if querySize == ""  {
		return defaultSize, nil
	}
	size, err := strconv.Atoi(querySize)
	if err != nil {
		return 0, err
	}
	if size < 32 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}
	return size, nil
}

func QREncode(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("content")
	if utf8.RuneCountInString(str) > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("content is too large, gtfo"))
		return
	}
	size, err := GetSize(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("size must be a positive integer greater than 32"))
		return
	}

	if str == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid text, content must contain a non empty string"))
		return
	}
	var png []byte
	png, err = qrcode.Encode(str, qrcode.Medium, size)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error converting text content to qr code"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "image/png")
	w.Write(png)
}

func Encrypt(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	password := r.URL.Query().Get("password")
	if message == "" || password == "" || (utf8.RuneCountInString(message) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	ciphertext := encrypt([]byte(base64.StdEncoding.EncodeToString([]byte(message))), password)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain")
	w.Write(ciphertext)
}

func EncryptWithQR(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	password := r.URL.Query().Get("password")
	if message == "" || password == "" || (utf8.RuneCountInString(message) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	ciphertext := encrypt([]byte(base64.StdEncoding.EncodeToString([]byte(message))), password)
	var png []byte
	png, err := qrcode.Encode(string(ciphertext), qrcode.Medium, 256)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error converting text content to qr code"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "image/png")
	w.Write(png)
}

func Decrypt(w http.ResponseWriter, r *http.Request) {
	cipher := r.URL.Query().Get("cipher")
	password := r.URL.Query().Get("password")
	if cipher == "" || password == "" || (utf8.RuneCountInString(cipher) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	plaintext := decrypt([]byte(base64.StdEncoding.EncodeToString([]byte(cipher))), password)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain")
	w.Write(plaintext)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
