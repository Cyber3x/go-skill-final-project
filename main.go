package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"example.com/main/middleware"
	"gopkg.in/yaml.v3"
)

type config struct {
	JMBAG string `ymal:"jmbag"`
	HTTP  struct {
		Address string `ymal:"address"`
		Port    int    `ymal:"port"`
	} `ymal:"http"`
	Users []struct {
		Name     string `ymal:"name"`
		JMBAG    string `ymal:"jmbag"`
		Password string `ymal:"password"`
	} `ymal:"users"`
}

var ConfigData config
var StudentFilePath string

func main() {
	err := ParseYAMLFile("./config.yaml", &ConfigData)

	if err != nil {
		log.Fatalln(err)
	}

	StudentFilePath, err = filepath.Abs("./studentData.txt")

	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/jmbag", JMBAGHandler)
	http.HandleFunc("/sum", middleware.BasicAuth(SumHandler, ConfigData.Users))
	http.HandleFunc("/multiply", middleware.BasicAuth(MultiplyHandler, ConfigData.Users))
	http.HandleFunc("/fetch", middleware.BasicAuth(FetchHandler, ConfigData.Users))
	http.HandleFunc("/"+ConfigData.JMBAG, middleware.BasicAuth(JMBAGFileHandler, ConfigData.Users))

	fmt.Println("starting server")

	address := ConfigData.HTTP.Address + ":" + strconv.Itoa(ConfigData.HTTP.Port)

	err = http.ListenAndServe(address, nil)
	if err != nil {

		log.Fatal(err)
	}
}

// This will parse the file at an absolute path
//
// Parsed data will be put in out interface
func ParseYAMLFile(path string, out interface{}) error {
	config_path, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	data, err := os.ReadFile(config_path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)

	if err != nil {
		return err
	}

	return nil
}
