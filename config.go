package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type Conf struct {
	DestLocalPath string
	FilesDir string
	ServePort string

	AuthUser string
	AuthPassword string
}

//var config *Conf

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
