package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goji/httpauth"
)

var config *Conf

func init() {
	config = GetConf()
}

func main() {

	log.Println("logStarting...")
	//http.HandleFunc("/upload", upload)
	authHandle := httpauth.SimpleBasicAuth(config.AuthUser, config.AuthPassword)

	s := http.FileServer(http.Dir(config.FilesDir))
	static_fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static_fs))
	http.Handle("/files/", authHandle(http.StripPrefix("/files/", s)))
	http.Handle("/upload", authHandle(http.HandlerFunc(upload)))
	http.Handle("/delete", authHandle(http.HandlerFunc(delFile)))
	http.Handle("/videos", authHandle(http.HandlerFunc(showVideos)))

	err := http.ListenAndServe(config.ServePort, nil)
	//err := http.ListenAndServeTLS(config.ServePort, "balancem.ltd_chain.crt", "balancem.ltd_key.key",nil)
	if err != nil {
		fmt.Println("ServerError: %v", err)
	}
}
