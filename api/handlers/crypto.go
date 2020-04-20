package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/skip2/go-qrcode"
)

func Encrypt(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	password := r.URL.Query().Get("password")
	if message == "" || password == "" || (utf8.RuneCountInString(message) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	log.Printf("message: %s, password: %s", message, password)
	ciphertext := base64.StdEncoding.EncodeToString(encrypt([]byte(message), password))

	log.Printf("cipher: %s", ciphertext)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte(ciphertext))
}

func EncryptWithQR(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	password := r.URL.Query().Get("password")
	if message == "" || password == "" || (utf8.RuneCountInString(message) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	ciphertext := base64.StdEncoding.EncodeToString(encrypt([]byte(message), password))
	var png []byte
	png, err := qrcode.Encode(ciphertext, qrcode.Medium, 256)
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
	payload := r.URL.Query().Get("cipher")
	password := r.URL.Query().Get("password")
	if payload == "" || password == "" || (utf8.RuneCountInString(payload) > 10000 || utf8.RuneCountInString(password) > 1000) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error decoding cipher"))
		return
	}
	plaintext := decrypt(cipherBytes, password)
	log.Printf("payload: %s, password: %s, plaintext: %s", payload, password, plaintext)
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
	// need this to be deterministic because we don't actually store anything
	// if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
	// 	panic(err.Error())
	// }
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
