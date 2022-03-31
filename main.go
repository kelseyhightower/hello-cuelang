package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

var (
	configFilePath string
)

type Config struct {
	HTTP     HTTPConfig     `json:"http"`
	Database DatabaseConfig `json:"database"`
}

type HTTPConfig struct {
	ListenPort int `json:"listen_port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func main() {
	flag.StringVar(&configFilePath, "c", "config.json", "Path to .json or .cue configuration file")
	flag.Parse()

	fmt.Println("Starting hello cue...")
	fmt.Printf("Loading configuration file: %s\n", configFilePath)
	c, err := configFromFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Running Configuration\n")
	fmt.Printf("    HTTP Port: %v\n", c.HTTP.ListenPort)
	fmt.Printf("    Database Host: %s\n", c.Database.Host)
	fmt.Printf("    Database User: %s\n", c.Database.User)
}

func configFromFile(path string) (*Config, error) {
	switch filepath.Ext(path) {
	case "":
		return nil, fmt.Errorf("missing file extension, must be .json or .cue")
	case ".json":
		return configFromJSONFile(path)
	case ".cue":
		return configFromCueFile(path)
	default:
		return nil, fmt.Errorf("invalid configuration file type: %s", path)
	}
}

func configFromCueFile(path string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ctx := cuecontext.New()
	v := ctx.CompileBytes(data)

	httpListenPort, _ := v.LookupPath(cue.ParsePath("config.http.listen_port")).Int64()
	databaseHost, _ := v.LookupPath(cue.ParsePath("config.database.host")).String()
	databaseUser, _ := v.LookupPath(cue.ParsePath("config.database.user")).String()
	databasePassword, _ := v.LookupPath(cue.ParsePath("config.database.password")).String()

	c.HTTP.ListenPort = int(httpListenPort)
	c.Database.Host = databaseHost
	c.Database.User = databaseUser
	c.Database.Password = databasePassword

	return &c, nil
}

func configFromJSONFile(path string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
