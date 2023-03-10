package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type config struct {
	MetaServerList  []string
	BlockServerList []string
	MetaHoldNum     int
}

func main() {

	configFile := flag.String("f", "", "(required) Config file, absolute path")
	flag.Parse()

	workDir, err := os.Getwd()
	log.Println("Servers start at", workDir)

	serverConfig, err := getConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// raft meta store

	for _, metaAddr := range serverConfig.MetaServerList {
		go launchMetaServer(metaAddr, serverConfig.BlockServerList[:serverConfig.MetaHoldNum])
	}

	// chord block store

	DHT := initDHTMap(serverConfig.BlockServerList)

	for _, blockAddr := range serverConfig.BlockServerList {
		go launchBlockServer(blockAddr, DHT[blockAddr])
	}
}

func getConfig(fileName string) (*config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sConfig config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sConfig)
	if err != nil {
		return nil, err
	}
	return &sConfig, nil
}

// raft meta server
func launchMetaServer(addr string, blockAddrs []string) {

}

// chord block server
func launchBlockServer(addr string, DHT []string) {

}

func initDHTMap(addrs []string) map[string][]string {
	return nil
}
