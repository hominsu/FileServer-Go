package main

import (
	"File-Uploader-go/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)

	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)

	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)

	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/file/update", handler.FileUpdateMetaHandler)

	http.HandleFunc("/file/delete", handler.FileDelHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server, err: ", err.Error())
	}
}
