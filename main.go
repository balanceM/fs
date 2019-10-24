package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/upload", upload)
	//fmt.Println("Starting...")
	log.Println("logStarting...")
	//http.HandleFunc("/upload", upload)
	s := http.FileServer(http.Dir("/home/banana"))
	http.Handle("/files/", http.StripPrefix("/files/", s))
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":3000", nil)
}

const uploadHTML = `
<html>
  <head>
	<title>选择文件</title>
  </head>
  <body>
	<form enctype="multipart/form-data" action="/upload" method="post">
	  <input type="file" name="uploadfile" />
	  <input type="submit" value="上传" />
  	</form>
  </body>
</html>
`

const destLocalPath = "/home/banana/Public/"

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(uploadHTML))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("GET method....")
		index(w, r)
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

	localpath := fmt.Sprintf("%s%s", destLocalPath, handler.Filename)
	localfd, err := os.Create(localpath)
	if err != nil {
		fmt.Println("3333", err)
		w.Write([]byte("upload failed!"))
		return
	}

	//localfd, err := os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE, 0666)
	//if err != nil {
	//	fmt.Println(err)
	//	w.Write([]byte("upload failed!"))
	//	return
	//}
	defer localfd.Close()

	io.Copy(localfd, clientfd)
	w.Write([]byte("upload finish ~.~ !"))
}
