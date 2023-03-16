package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"surfstore/pkg"
)

type ClientConfig struct {
	MetaServerList []string
	Basedir        string
	BlockSize      int
}

func main() {
	configFile := flag.String("f", "", "(required) Config file, absolute path")
	flag.Parse()

	workDir, err := os.Getwd()
	log.Println("Client start at", workDir)

	clientConfig, err := getClientConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient := pkg.NewSurfstoreRPCClient(clientConfig.MetaServerList, clientConfig.Basedir, clientConfig.BlockSize)
	addr := ""
	err = rpcClient.GetBlockStoreAddr(&pkg.BlockHash{Hash: "1A3A5A1A"}, "localhost:8081", &addr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Println(addr)
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
