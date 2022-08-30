package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/h2non/bimg"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":3000", nil)
}

// The mime type of the image is changed, it is compressed and then saved in the specified folder.
func imageProcessing(buffer []byte, quality int) ([]byte, error) {

	mimeType := http.DetectContentType(buffer)
	// check mime type
	fmt.Println(mimeType)
	if !strings.Contains(mimeType, "image") {
		return buffer, nil
	}

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return converted, err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return processed, err
	}

	return processed, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("picture")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fileBytes, err := imageProcessing(bytes, 40)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}
