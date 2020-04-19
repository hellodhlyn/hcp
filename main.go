package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type filePipe struct {
	data      []byte
	dataSize  int64
	readCount int
}

var filePool = sync.Map{}

func Download(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fileKey := p.ByName("key")
	pipe, ok := filePool.Load(fileKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(pipe.(*filePipe).dataSize, 10))
	_, err := io.Copy(w, bytes.NewReader(pipe.(*filePipe).data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = file.Close()

	fileKey := p.ByName("key")
	filePool.Store(fileKey, &filePipe{data: fileBytes, dataSize: handler.Size})

	for {
		select {
		case <-r.Context().Done():
			filePool.Delete(fileKey)
			return
		default:
			continue
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/:key", Download)
	router.POST("/:key", Upload)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	fmt.Println("Server listening port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
