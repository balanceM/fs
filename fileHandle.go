package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

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
