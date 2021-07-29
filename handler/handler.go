package handler

import (
	"File-Uploader-go/meta"
	"File-Uploader-go/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler 处理上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// TODO: 返回 html 页面
		data, err := ioutil.ReadFile("static/view/upload.html")
		if err != nil {
			io.WriteString(w, "Internal server error")
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// TODO: 接收文件流并存储在本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to get data, err: ", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:01:02"),
		}
		fmt.Println(head.Filename)

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Println("Failed to create file, err: ", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("Failed to save data into file, err: ", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		err = meta.UpdateFileMetaDB(fileMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler 获取文件源信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form["fileHash"][0]
	//fMeta := meta.GetFileMeta(fileHash)
	fMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUpdateMetaHandler 更新文件元信息
func FileUpdateMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileHash := r.Form.Get("fileHash")
	newFileName := r.Form.Get("fileName")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	curFileMeta := meta.GetFileMeta(fileHash)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DownloadHandler 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("fileHash")
	fm := meta.GetFileMeta(fileHash)

	file, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileDelHandler 删除文件及其元信息
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("fileHash")

	fMeta := meta.GetFileMeta(fileHash)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileHash)

	w.WriteHeader(http.StatusOK)
}
