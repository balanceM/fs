package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Conf struct {
	DestLocalPath string
	FilesDir string
	ServePort string
}

var config *Conf

func main() {
	config = GetConf()

	//mux := http.NewServeMux()
	//mux.HandleFunc("/upload", upload)
	//fmt.Println("Starting...")
	log.Println("logStarting...")
	//http.HandleFunc("/upload", upload)
	s := http.FileServer(http.Dir(config.FilesDir))
	http.Handle("/files/", http.StripPrefix("/files/", s))
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(config.ServePort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func GetConf() *Conf{
	data, _ := ioutil.ReadFile("config.yml")
	t := Conf{}
	err := yaml.Unmarshal(data, &t)
	fmt.Println("初始数据", t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &t
}

type indexParam struct {
	Message string
}

func index(w http.ResponseWriter, r *http.Request, message string) {
	upload_template, err := template.ParseFiles("html/upload.html")
	if err != nil {
		fmt.Println("parse file err:",err)
		w.Write([]byte("parse htmlfile failed!"))
		return
	}
	p := &indexParam{
		Message: message,
	}
	upload_template.Execute(w, p)
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

	io.Copy(localfd, clientfd)
	index(w, r, fmt.Sprintf("[%s] uploaded!", handler.Filename))
	//w.Header().Set("Location", "/up)load")
	//w.WriteHeader(301)
}
