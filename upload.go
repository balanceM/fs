package main

import (
	"fmt"
	"github.com/juju/ratelimit"
	"io"
	"net/http"
	"os"
	"html/template"
)

func index(w http.ResponseWriter, r *http.Request, message string) {
	upload_template, err := template.ParseFiles("html/upload.html")
	if err != nil {
		fmt.Println("parse file err:",err)
		w.Write([]byte("parse htmlfile failed!"))
		return
	}
	upload_template.Execute(w, map[string]string{"Message": message})
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		index(w, r, "Ready!")
		return
	}
	fmt.Println("POST method....")

	r.ParseMultipartForm(32 << 20)
	clientfd, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("2222", err)
		w.Write([]byte("upload failed!"))
		return
	}
	defer clientfd.Close()

	localpath := fmt.Sprintf("%s%s", config.DestLocalPath, handler.Filename)
	localfd, err := os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("upload failed!"))
		return
	}
	defer localfd.Close()

	//---》 ratelimit 限速10m/s
	// Bucket adding 100KB every second, holding max 100KB
	bucket := ratelimit.NewBucketWithRate(100*1024, 100*1024)
	io.Copy(localfd, ratelimit.Reader(clientfd, bucket))
	// 《-- ratelimit
	//index(w, r, fmt.Sprintf("[%s] uploaded!", handler.Filename))
}

func showVideos(w http.ResponseWriter, r *http.Request) {
	video_template, err := template.ParseFiles("html/abx.html")
	if err != nil {
		fmt.Println("parse file err:",err)
		w.Write([]byte("parse htmlfile failed!"))
		return
	}
	video_template.Execute(w, nil)
}