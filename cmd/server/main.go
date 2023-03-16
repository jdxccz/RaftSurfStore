package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"sort"
	"surfstore/pkg"

	"google.golang.org/grpc"
)

const INF int = 999999999

type ServerConfig struct {
	MetaServerList      []string
	BlockServerList     []string
	DefaultBlockAddrNum int
	BlockCircleNodeNum  int
}

func main() {

	configFile := flag.String("f", "", "(required) Config file, absolute path")
	flag.Parse()

	workDir, err := os.Getwd()
	log.Println("Servers start at", workDir)

	serverConfig, err := getServerConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// raft meta store

	for id := range serverConfig.MetaServerList {
		go launchMetaServer(id, serverConfig.MetaServerList, serverConfig.BlockServerList[:serverConfig.DefaultBlockAddrNum])
	}

	// chord block store

	DHTs, PosMap, err := initDHTMap(serverConfig.BlockServerList, serverConfig.BlockCircleNodeNum)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// log.Println(DHTs)

	for _, blockAddr := range serverConfig.BlockServerList {
		go launchBlockServer(blockAddr, DHTs[blockAddr], serverConfig.BlockCircleNodeNum, PosMap[blockAddr])
	}

	for {
	}
}

func getServerConfig(fileName string) (*ServerConfig, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sConfig ServerConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sConfig)
	if err != nil {
		return nil, err
	}
	return &sConfig, nil
}

// raft meta server
func launchMetaServer(id int, metaAddr []string, blockAddrs []string) {

}

// chord block server
func launchBlockServer(addr string, DHT []pkg.Pair, N int, pos int) {
	log.Println("Start Block Store Server: ", addr)
	blockStore := pkg.NewBlockStore(DHT, N, pos, addr)
	grpcServer := grpc.NewServer()
	pkg.RegisterBlockStoreServer(grpcServer, blockStore)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func initDHTMap(addrs []string, N int) (map[string][]pkg.Pair, map[string]int, error) {
	crypto := "SurfStore"
	var pairs []pkg.Pair
	DHT := make(map[string][]pkg.Pair)
	posMap := make(map[string]int)
	for _, addr := range addrs {
		hash := pkg.GetHash([]byte(crypto + addr))
		pos, err := pkg.GetChordPosition(N, hash)
		if err != nil {
			return nil, nil, err
		}
		posMap[addr] = pos
		pairs = append(pairs, pkg.Pair{Key: pos, Value: addr})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Key < pairs[j].Key })
	for _, pair := range pairs {
		dht := []pkg.Pair{}
		i := 1
		for i < N {
			pos := (pair.Key + i) % N
			successor := pkg.Pair{Key: -1, Value: ""}
			// could change to binary search
			for _, pr := range pairs {
				if pr.Key >= pos {
					successor = pr
					break
				}
			}
			if successor.Value == "" {
				successor = pairs[0]
			}
			successor.Key = (successor.Key + N - pair.Key) % N
			dht = append(dht, successor)
			i *= 2
		}
		DHT[pair.Value] = dht
	}
	return DHT, posMap, nil
}
