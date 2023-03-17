package pkg

import (
	context "context"
	"errors"
	"log"
	sync "sync"

	grpc "google.golang.org/grpc"
)

type BlockStore struct {
	BlockMap map[string]*Block
	DHT      []Pair
	RingNum  int
	Position int
	Addr     string
	sync.RWMutex

	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	bs.RWMutex.RLock()
	value, ok := bs.BlockMap[blockHash.Hash]
	bs.RWMutex.RUnlock()
	if ok {
		return value, nil
	} else {
		return nil, errors.New("hash " + blockHash.Hash + " is not in the block store!")
	}
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	bs.RWMutex.Lock()
	hash := GetHash(block.BlockData)
	bs.BlockMap[hash] = block
	bs.RWMutex.Unlock()
	return &Success{Flag: true}, nil
}

func (bs *BlockStore) GetBlockStoreAddr(ctx context.Context, blockHash *BlockHash) (*BlockStoreAddr, error) {
	pos, err := GetChordPosition(bs.RingNum, blockHash.Hash)
	if err != nil {
		return nil, err
	}
	if bs.Position == pos {
		return &BlockStoreAddr{Addr: bs.Addr}, nil
	}
	pos = (pos + bs.RingNum - bs.Position) % bs.RingNum
	addr := ""
	bs.RWMutex.RLock()
	if pos <= bs.DHT[0].Key {
		bs.RWMutex.RUnlock()
		return &BlockStoreAddr{Addr: bs.DHT[0].Value}, nil
	}
	for i := 1; i < len(bs.DHT); i++ {
		if bs.DHT[i].Key >= pos {
			addr = bs.DHT[i-1].Value
			break
		}
	}
	bs.RWMutex.RUnlock()
	if addr == "" {
		addr = bs.DHT[len(bs.DHT)-1].Value
	}
	log.Println("Nexthop:", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := NewBlockStoreClient(conn)
	BlockAddr, err := c.GetBlockStoreAddr(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	conn.Close()
	return BlockAddr, nil
}

func (bs *BlockStore) CleanBlockStore(ctx context.Context, blockhashes *BlockHashes) (*Success, error) {
	newMap := make(map[string]*Block)
	bs.RWMutex.RLock()
	for _, hash := range blockhashes.Hashes {
		newMap[hash] = bs.BlockMap[hash]
	}
	bs.RWMutex.RUnlock()
	bs.RWMutex.Lock()
	bs.BlockMap = newMap
	bs.RWMutex.Unlock()
	return &Success{Flag: true}, nil
}

func NewBlockStore(DHT []Pair, N int, pos int, addr string) *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
		DHT:      DHT,
		RingNum:  N,
		Position: pos,
		Addr:     addr,
	}
}
