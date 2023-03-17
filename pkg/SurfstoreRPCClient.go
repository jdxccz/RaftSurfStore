package pkg

import (
	context "context"
	"math/rand"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCClient struct {
	MetaStoreAddrs []string
	BaseDir        string
	BlockSize      int
	BsAddrs        []string
}

const timeout int = 1

func (surfClient *RPCClient) GetBlock(blockHash string, blockStoreAddr string, block *Block) error {
	// connect to the server
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	b, err := c.GetBlock(ctx, &BlockHash{Hash: blockHash})
	if err != nil {
		conn.Close()
		return err
	}
	block.BlockData = b.BlockData
	block.BlockSize = b.BlockSize

	// close the connection
	return conn.Close()
}

func (surfClient *RPCClient) PutBlock(block *Block, blockStoreAddr string, succ *bool) error {
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	s, err := c.PutBlock(ctx, block)
	if err != nil {
		conn.Close()
		return err
	}
	*succ = s.Flag
	return conn.Close()
}

func (surfClient *RPCClient) GetBlockStoreAddr(blockHash *BlockHash, blockStoreAddr string, targetBlockStoreAddr *string) error {
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	b, err := c.GetBlockStoreAddr(ctx, blockHash)
	if err != nil {
		conn.Close()
		return err
	}
	*targetBlockStoreAddr = b.Addr
	return conn.Close()
}

func (surfClient *RPCClient) GetFileInfoMap(serverFileInfoMap *map[string]*FileMetaData) error {
	idx := rand.Intn(len(surfClient.MetaStoreAddrs))
	mAddr := surfClient.MetaStoreAddrs[idx]
	st := Status{Flag: false}
	for !st.Flag {
		err := surfClient.GetStatus(mAddr, &st)
		if err == nil {
			mAddr = st.Addr
		} else {
			idx := rand.Intn(len(surfClient.MetaStoreAddrs))
			mAddr = surfClient.MetaStoreAddrs[idx]
			st = Status{Flag: false}
		}
	}

	conn, err := grpc.Dial(mAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewMetaStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	f, err := c.GetFileInfoMap(ctx, &emptypb.Empty{})
	if err != nil {
		conn.Close()
		return err
	} else {
		*serverFileInfoMap = f.FileInfoMap
		return conn.Close()
	}
}

func (surfClient *RPCClient) UpdateFile(fileMetaData *FileMetaData, latestVersion *int32) error {
	idx := rand.Intn(len(surfClient.MetaStoreAddrs))
	mAddr := surfClient.MetaStoreAddrs[idx]
	st := Status{Flag: false}
	for !st.Flag {
		err := surfClient.GetStatus(mAddr, &st)
		if err == nil {
			mAddr = st.Addr
		} else {
			idx := rand.Intn(len(surfClient.MetaStoreAddrs))
			mAddr = surfClient.MetaStoreAddrs[idx]
			st = Status{Flag: false}
		}
	}

	conn, err := grpc.Dial(mAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewMetaStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	v, err := c.UpdateFile(ctx, fileMetaData)
	if err != nil {
		conn.Close()
		return err
	} else {
		*latestVersion = v.Version
		return conn.Close()
	}
}

func (surfClient *RPCClient) GetDefualtBlockStoreAddrs(blockStoreAddrs *[]string) error {
	idx := rand.Intn(len(surfClient.MetaStoreAddrs))
	conn, err := grpc.Dial(surfClient.MetaStoreAddrs[idx], grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewMetaStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	bsAddr, err := c.GetDefualtBlockStoreAddrs(ctx, &emptypb.Empty{})
	if err != nil {
		conn.Close()
		return err
	} else {
		*blockStoreAddrs = bsAddr.BlockStoreAddrs
	}
	return conn.Close()
}

func (surfClient *RPCClient) GetStatus(mAddr string, status *Status) error {
	conn, err := grpc.Dial(mAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c := NewMetaStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	st, err := c.GetStatus(ctx, &emptypb.Empty{})
	if err != nil {
		conn.Close()
		return err
	} else {
		*status = Status{Flag: st.Flag, Addr: st.Addr}
	}
	return conn.Close()
}

// Create an Surfstore RPC client
func NewSurfstoreRPCClient(addrs []string, baseDir string, blockSize int, bsAddrs []string) RPCClient {

	return RPCClient{
		MetaStoreAddrs: addrs,
		BaseDir:        baseDir,
		BlockSize:      blockSize,
		BsAddrs:        bsAddrs,
	}
}
