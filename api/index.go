package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/skip2/go-qrcode"

	"github.com/trevor-atlas/weekend/api/handlers"
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
		GET("/encrypt", handlers.Encrypt).
		GET("/encrypt-with-qr", handlers.EncryptWithQR).
		GET("/decrypt", handlers.Decrypt).
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

