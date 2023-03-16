package pkg

import (
	context "context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MetaStoreInterface interface {
	GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error)

	UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error)

	GetDefualtBlockStoreAddrs(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddrs, error)

	GetStatus(ctx context.Context, _ *emptypb.Empty) (*Status, error)

	AppendEntries(ctx context.Context, input *AppendEntryInput) (*AppendEntryOutput, error)

	RequestVote(ctx context.Context, cstate *CandidateState) (*VoteResult, error)
}

type BlockStoreInterface interface {
	GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error)

	PutBlock(ctx context.Context, block *Block) (*Success, error)

	GetBlockStoreAddr(ctx context.Context, blockHash *BlockHash) (*BlockStoreAddr, error)

	CleanBlockStore(ctx context.Context, blockhashes *BlockHashes) (*Success, error)
}

type ClientInterface interface {
	// MetaStore
	GetFileInfoMap(mAddr string, serverFileInfoMap *map[string]*FileMetaData) error

	UpdateFile(mAddr string, fileMetaData *FileMetaData, latestVersion *int32) error

	GetDefualtBlockStoreAddrs(blockStoreAddrs *[]string) error

	GetStatus(mAddr string, status *Status) error

	// BlockStore
	GetBlock(blockHash string, blockStoreAddr string, block *Block) error

	PutBlock(block *Block, blockStoreAddr string, succ *bool) error

	GetBlockStoreAddr(blockHash *BlockHash, lockStoreAddr string, targetBlockStoreAddr *string) error
}
