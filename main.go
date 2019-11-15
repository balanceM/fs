package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"log"
	"net/http"
)

var config *Conf
func init()  {
	config = GetConf()
}

func main() {

	log.Println("logStarting...")
	//http.HandleFunc("/upload", upload)
	authHandle := httpauth.SimpleBasicAuth(config.AuthUser, config.AuthPassword)

	s := http.FileServer(http.Dir(config.FilesDir))
	static_fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix( "/static/", static_fs))
	http.Handle("/files/", authHandle(http.StripPrefix("/files/", s)))
	http.Handle("/upload", authHandle(http.HandlerFunc(upload)))
	http.Handle("/delete", authHandle(http.HandlerFunc(delFile)))
	err := http.ListenAndServeTLS(config.ServePort, "balancem.ltd_chain.crt", "balancem.ltd_key.key",nil)
	if err != nil {
		fmt.Println(err)
	}
}

