package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/goji/httpauth"
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

	AuthUser string
	AuthPassword string
}

var config *Conf

func main() {
	config = GetConf()

	log.Println("logStarting...")
	//http.HandleFunc("/upload", upload)
	authHandle := httpauth.SimpleBasicAuth(config.AuthUser, config.AuthPassword)

	s := http.FileServer(http.Dir(config.FilesDir))
	static_fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix( "/static/", static_fs))
	http.Handle("/files/", authHandle(http.StripPrefix("/files/", s)))
	http.Handle("/upload", authHandle(http.HandlerFunc(upload)))
	http.Handle("/delete", authHandle(http.HandlerFunc(delFile)))
	err := http.ListenAndServe(config.ServePort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type FileInfo struct {
	Type string
	Name string
	Path string
	Files []*FileInfo
}

func filesshow(w http.ResponseWriter, r *http.Request) {
	var funcMaps = template.FuncMap{
		"add": func(a, b string) string {
			return a+b
		},
	}

	files_template, err := template.New("filesshow.html").Funcs(funcMaps).ParseFiles("html/filesshow.html", "html/ul_list.html")
	if err != nil {

		fmt.Println(err)
		w.Write([]byte("files_template failed!"))
	}

	files, err := getFiles(w, config.DestLocalPath)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("error!"))
		return
	}
	data := struct {
		Files []*FileInfo
	}{files}
	files_template.Execute(w, data)
}

func getFiles(w http.ResponseWriter, dirpath string) ([]*FileInfo, error){
	rd, err := ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var files []*FileInfo
	for _, fi := range rd {
		fileInfo := &FileInfo{}
		if fi.IsDir() {
			fileInfo.Type = "dir"
			fmt.Println(dirpath+fi.Name()+"/")
			fileInfo.Files, err = getFiles(w, dirpath+fi.Name()+"/")
			if err != nil {
				return nil, err
			}
		} else {
			fileInfo.Type = "file"
		}
		fileInfo.Name = fi.Name()
		fileInfo.Path = dirpath
		files = append(files, fileInfo)
	}
	return files, nil
}

func delFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		filesshow(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~")
	filepath := r.PostFormValue("filepath")
	fmt.Println(filepath)
	info, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Failed!"))
		return
	}

	// remove
	if info.IsDir() {
		err = os.RemoveAll(filepath)
	} else {
		err = os.Remove(filepath)
	}
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Failed!"))
		return
	}
	w.Write([]byte("Succeed!"))
	return
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
