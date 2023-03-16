package pkg

import (
	context "context"
	"errors"
	"math"
	"math/rand"
	sync "sync"
	"time"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MetaStore struct {
	FileMetaMap            map[string]*FileMetaData
	DefaultBlockStoreAddrs []string
	status                 int // 0 = leader 1 = candidate 2 = follower
	leaderAddr             string
	term                   int64
	log                    []*UpdateOperation
	serverList             []string
	serverId               int
	commitIndex            int
	prevLogIndexList       []int
	// lastVoteTerm           int
	logLock     sync.RWMutex
	metaMapLock sync.RWMutex
	statusLock  sync.RWMutex
	// logCommitChan          chan int
	comTimoutState bool
	UnimplementedMetaStoreServer
}

func (m *MetaStore) GetStatus(ctx context.Context, _ *emptypb.Empty) (*Status, error) {
	m.statusLock.RLock()
	var ls *Status
	if m.status == 0 {
		ls = &Status{Flag: true, Addr: ""}
	} else {
		ls = &Status{Flag: false, Addr: m.leaderAddr}
	}
	m.statusLock.RUnlock()
	return ls, nil
}

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	// ensure leader
	m.statusLock.RLock()
	if m.status != 0 {
		m.statusLock.RUnlock()
		return nil, errors.New("Not Leader")
	}
	m.statusLock.RUnlock()
	// add log
	m.logLock.Lock()
	id := len(m.log)
	m.log = append(m.log, &UpdateOperation{Term: m.term, Op: "GetFileInfoMap", FileMetaData: nil})
	m.logLock.Unlock()

	for m.commitIndex < id {
		time.Sleep(1 * time.Millisecond)
	}

	m.statusLock.RLock()
	if m.status != 0 {
		m.statusLock.RUnlock()
		return nil, errors.New("Fake Leader")
	}
	m.statusLock.RUnlock()

	// readmap
	m.metaMapLock.RLock()
	copyFMM := make(map[string]*FileMetaData)
	for k, v := range m.FileMetaMap {
		copyFMM[k] = v
	}
	m.metaMapLock.RUnlock()
	return &FileInfoMap{FileInfoMap: copyFMM}, nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	// ensure leader
	m.statusLock.RLock()
	if m.status != 0 {
		m.statusLock.RUnlock()
		return nil, errors.New("Not Leader")
	}
	m.statusLock.RUnlock()
	// add log
	m.logLock.Lock()
	id := len(m.log)
	m.log = append(m.log, &UpdateOperation{Term: m.term, Op: "UpdateFile", FileMetaData: fileMetaData})
	m.logLock.Unlock()

	for m.commitIndex < id {
		time.Sleep(1 * time.Millisecond)
	}

	m.statusLock.RLock()
	if m.status != 0 {
		m.statusLock.RUnlock()
		return nil, errors.New("Fake Leader")
	}
	m.statusLock.RUnlock()

	// update map
	return updateMetaMap(m, fileMetaData)
}

func (m *MetaStore) GetDefualtBlockStoreAddrs(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddrs, error) {
	return &BlockStoreAddrs{BlockStoreAddrs: m.DefaultBlockStoreAddrs}, nil
}

func (m *MetaStore) AppendEntries(ctx context.Context, input *AppendEntryInput) (*AppendEntryOutput, error) {
	succ := true
	matchIndex := -1
	// if follower term > leader term => new leader nothing change
	if input.Term < m.term {
		return &AppendEntryOutput{ServerId: int64(m.serverId), Term: m.term, Success: false, MatchedIndex: -1}, nil
	} else if input.Term > m.term {
		m.statusLock.Lock()
		m.status = 2
		m.statusLock.Unlock()
		m.term = input.Term
		succ = false
	} else {
		m.statusLock.Lock()
		m.status = 2
		m.statusLock.Unlock()
	}

	// find match point
	if input.PrevLogIndex <= -1 {
		m.logLock.Lock()
		m.log = input.Entries
		m.logLock.Unlock()
		matchIndex += len(input.Entries)
		NewcommitIndex := int(math.Min(float64(matchIndex), float64(input.LeaderCommit)))
		commitLog(m, NewcommitIndex)
		m.commitIndex = NewcommitIndex
		return &AppendEntryOutput{ServerId: int64(m.serverId), Term: m.term, Success: succ, MatchedIndex: int64(matchIndex)}, nil
	}

	if len(m.log)-1 >= int(input.PrevLogIndex) && m.log[input.PrevLogIndex].Term == input.PrevLogTerm {
		matchIndex = int(input.PrevLogIndex)
	} else {
		succ = false
	}

	// if term equal
	if matchIndex >= 0 {
		m.logLock.Lock()
		m.log = m.log[:matchIndex+1]
		m.log = append(m.log, input.Entries...)
		m.logLock.Unlock()
		matchIndex += len(input.Entries)
		NewcommitIndex := int(math.Min(float64(matchIndex), float64(input.LeaderCommit)))
		commitLog(m, NewcommitIndex)
		m.commitIndex = NewcommitIndex
	}

	return &AppendEntryOutput{ServerId: int64(m.serverId), Term: m.term, Success: succ, MatchedIndex: int64(matchIndex)}, nil
}

func (m *MetaStore) RequestVote(ctx context.Context, candidateState *CandidateState) (*VoteResult, error) {
	vote := false

	if candidateState.Term > m.term {
		m.term = candidateState.Term
		m.statusLock.Lock()
		m.status = 2
		m.statusLock.Unlock()
		Pterm := int64(-1)
		if len(m.log) > 0 {
			Pterm = m.log[len(m.log)-1].Term
		}
		if candidateState.LastLogTerm > Pterm || (candidateState.LastLogTerm == Pterm && int(candidateState.LastLogIndex) >= len(m.log)-1) {
			vote = true
		}
	}

	return &VoteResult{Term: m.term, VoteGranted: vote}, nil
}

func NewMetaStore(id int, metaStoreAddrs []string, blockStoreAddrs []string) *MetaStore {

	prevlogindexlist := []int{}
	for range metaStoreAddrs {
		prevlogindexlist = append(prevlogindexlist, 0)
	}
	mt := &MetaStore{
		FileMetaMap:            map[string]*FileMetaData{},
		DefaultBlockStoreAddrs: blockStoreAddrs,
		status:                 2,
		term:                   0,
		log:                    make([]*UpdateOperation, 0),
		commitIndex:            0,
		serverList:             metaStoreAddrs,
		serverId:               id,
		leaderAddr:             "",
		comTimoutState:         false,
		// lastVoteTerm:           -1,
	}
	mt.log = append(mt.log, &UpdateOperation{Term: 0, FileMetaData: nil})

	go metaStoreExec(mt)
	return mt
}

func metaStoreExec(m *MetaStore) {
	communicateInterval := 100
	communicateTimeout := 250
	electionTimeout := [2]int{150, 300}
	for {
		m.statusLock.RLock()
		st := m.status
		m.statusLock.RUnlock()
		switch st {
		case 0:
			// leader
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(communicateInterval)*time.Millisecond)
			defer cancel()
			go sendHeartBeat(ctx, m)
			time.Sleep(time.Duration(communicateInterval) * time.Millisecond)
		case 1:
			// candidate
			rand.Seed(time.Now().UnixNano())
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rand.Intn(electionTimeout[1]-electionTimeout[0]+1)+electionTimeout[0])*time.Millisecond)
			defer cancel()
			holdElection(ctx, m)
			m.statusLock.RLock()
			if m.status == 1 {
				m.term++
			}
			m.statusLock.RUnlock()
		case 2:
			// follower
			if m.comTimoutState {
				m.statusLock.Lock()
				m.status = 1
				m.statusLock.Unlock()
				m.term += 1
			} else {
				m.comTimoutState = true
				time.Sleep(time.Duration(communicateTimeout) * time.Millisecond)
			}
		}
	}
}

func sendHeartBeat(ctx context.Context, m *MetaStore) {
	for id := range m.serverList {
		if id == m.serverId {
			continue
		}
		go appendExec(ctx, m, id)
	}
}

func appendExec(ctx context.Context, m *MetaStore, id int) {
	var AppendIn *AppendEntryInput
	var AppendOut *AppendEntryOutput
	conn, err := grpc.Dial(m.serverList[id], grpc.WithTimeout(100*time.Millisecond))
	defer conn.Close()
	if err != nil {
		return
	}
	follower := NewMetaStoreClient(conn)
	var Pterm int64
	var entries []*UpdateOperation
	if m.prevLogIndexList[id] <= -1 {
		Pterm = -1
	} else {
		Pterm = m.log[m.prevLogIndexList[id]].Term
	}
	if m.prevLogIndexList[id] == len(m.log)-1 {
		entries = []*UpdateOperation{}
	} else {
		entries = []*UpdateOperation{m.log[m.prevLogIndexList[id]+1]}
	}
	AppendIn = &AppendEntryInput{Term: m.term, PrevLogIndex: int64(m.prevLogIndexList[id]), PrevLogTerm: Pterm, Entries: entries, LeaderCommit: int64(m.commitIndex)}
	AppendOut, err = follower.AppendEntries(ctx, AppendIn)
	if err == nil {
		if AppendOut.Term > m.term {
			m.statusLock.Lock()
			m.status = 2
			m.statusLock.Unlock()
			return
		}
		if AppendOut.MatchedIndex < AppendIn.PrevLogIndex {
			m.prevLogIndexList[id]--
		} else {
			m.prevLogIndexList[id]++
		}
		if AppendOut.MatchedIndex >= int64(len(m.log))-1 {
			m.prevLogIndexList[id] = len(m.log) - 1
		}
	}
}

func holdElection(ctx context.Context, m *MetaStore) {
	ch := make(chan int, len(m.serverList)-1)
	for id := range m.serverList {
		if id == m.serverId {
			continue
		}
		go electionExec(ctx, m, id, ch)
	}
	succ := 0
	for {
		st := <-ch
		if st == 0 {
			succ++
		} else if st == 2 {
			m.statusLock.Lock()
			m.status = 2
			m.statusLock.Unlock()
			break
		}
		if succ >= int(len(m.serverList)/2) {
			m.statusLock.Lock()
			m.status = 0
			m.statusLock.Unlock()
			break
		}
	}
}

func electionExec(ctx context.Context, m *MetaStore, id int, ch chan int) {
	var cstate *CandidateState
	var result *VoteResult
	conn, err := grpc.Dial(m.serverList[id], grpc.WithTimeout(100*time.Millisecond))
	defer conn.Close()
	if err != nil {
		return
	}
	follower := NewMetaStoreClient(conn)
	Pterm := int64(-1)
	if len(m.log) > 0 {
		Pterm = m.log[len(m.log)-1].Term
	}
	cstate = &CandidateState{Term: m.term, CandidateId: int64(m.serverId), LastLogIndex: int64(len(m.log) - 1), LastLogTerm: Pterm}
	result, err = follower.RequestVote(ctx, cstate)
	if err == nil {
		if result.Term > m.term {
			m.term = result.Term
			ch <- 2
		} else {
			if result.VoteGranted {
				ch <- 0
			} else {
				ch <- 1
			}
		}
	}
}

func commitLog(m *MetaStore, newidx int) {
	m.logLock.RLock()
	for i := m.commitIndex + 1; i < newidx+1; i++ {
		switch m.log[i].Op {
		case "GetFileInfoMap":
			continue
		case "UpdateFile":
			updateMetaMap(m, m.log[i].FileMetaData)
		}
	}
	m.logLock.RUnlock()
}

func updateMetaMap(m *MetaStore, fileMetaData *FileMetaData) (*Version, error) {
	var err error
	var v int32
	m.metaMapLock.Lock()
	if _, ok := m.FileMetaMap[fileMetaData.Filename]; ok {
		if fileMetaData.Version == m.FileMetaMap[fileMetaData.Filename].Version+1 {
			err = nil
			m.FileMetaMap[fileMetaData.Filename] = fileMetaData
			v = fileMetaData.Version
		} else {
			err = nil
			v = -1
		}
	} else {
		if fileMetaData.Version != 1 {
			err = errors.New("Wrong Version to Create a file!")
		} else {
			err = nil
		}
		m.FileMetaMap[fileMetaData.Filename] = &FileMetaData{Filename: fileMetaData.Filename, Version: int32(1), BlockHashList: fileMetaData.BlockHashList}
		v = 1
	}
	m.metaMapLock.Unlock()
	return &Version{Version: v}, err
}
