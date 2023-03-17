package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"surfstore/pkg"
)

type ClientConfig struct {
	MetaServerList []string
	BlockSize      int
	BlockAddrs     []string
}

func main() {
	configFile := flag.String("f", "", "(required) Config file, absolute path")
	baseDir := flag.String("b", "", "(required) Base dir, absolute path")
	flag.Parse()

	workDir, err := os.Getwd()
	log.Println("Client start at", workDir)

	clientConfig, err := getClientConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	Client := pkg.NewSurfstoreRPCClient(clientConfig.MetaServerList, *baseDir, clientConfig.BlockSize, clientConfig.BlockAddrs)

	clientConfig.BlockAddrs = pkg.ClientSync(Client)

	err = writeClientConfig(*configFile, clientConfig)

	if err != nil {
		log.Fatalf("fail to write back new BS addrs!%s", err)
	}
}

func getClientConfig(fileName string) (*ClientConfig, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sConfig ClientConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sConfig)
	if err != nil {
		return nil, err
	}
	return &sConfig, nil
}

func writeClientConfig(fileName string, config *ClientConfig) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, jsonBytes, 0644); err != nil {
		return err
	}
	return nil
}
