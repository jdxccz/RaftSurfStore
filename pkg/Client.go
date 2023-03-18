package pkg

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) []string {
	originindex, err := LoadMetaFromMetaFile(client.BaseDir)
	if err != nil {
		log.Fatalf("Error Occurs when load local map: %v", err)
	}
	log.Println("load meta data successfully!")
	localindex, err := synclocal(originindex, client.BaseDir, client.BlockSize)
	// PrintMetaMap(localindex)
	if err != nil {
		log.Fatalf("Error Occurs when sync local dic: %v", err)
	}
	log.Println("sync local successfully!")
	localindex, bsAddrs, err := syncRemote(localindex, client)
	if err != nil {
		log.Fatalf("Error Occurs when sync remote dic: %v", err)
	}
	log.Println("sync remote successfully!")
	err = WriteMetaFile(localindex, client.BaseDir)
	if err != nil {
		log.Fatalf("Error Occurs when write local map: %v", err)
	}
	// log.Println("write meta data successfully!")
	// PrintMetaMap(localindex)
	// remoteindex := make(map[string]*FileMetaData)
	// client.GetFileInfoMap(&remoteindex)
	// PrintMetaMap(remoteindex)
	return bsAddrs
}

func synclocal(oldindex map[string]*FileMetaData, basedir string, blocksize int) (newindex map[string]*FileMetaData, err error) {
	buf := make([]byte, blocksize)
	newindex = make(map[string]*FileMetaData)

	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == "index.db" {
			continue
		}
		fi, err := os.Open(ConcatPath(basedir, file.Name()))
		defer fi.Close()
		if err != nil {
			return nil, err
		}
		hashes := []string{}
		for {
			n, err := fi.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			} else {
				if err == io.EOF {
					break
				}
				hashes = append(hashes, GetHash(buf[:n]))
			}
		}
		if info, _ := os.Stat(ConcatPath(basedir, file.Name())); info.Size() == 0 {
			hashes = []string{EMPTYFILE_HASHVALUE}
		}
		_, ok := oldindex[file.Name()]
		if ok {
			if CompareStrings(oldindex[file.Name()].BlockHashList, hashes) {
				newindex[file.Name()] = oldindex[file.Name()]
			} else {
				newindex[file.Name()] = &FileMetaData{Filename: file.Name(), Version: oldindex[file.Name()].Version + 1, BlockHashList: hashes}
			}
		} else {
			newindex[file.Name()] = &FileMetaData{Filename: file.Name(), Version: int32(1), BlockHashList: hashes}
		}
	}
	for k, v := range oldindex {
		if _, ok := newindex[k]; ok {
			continue
		}
		if len(v.BlockHashList) == 1 && v.BlockHashList[0] == TOMBSTONE_HASHVALUE {
			newindex[k] = v
			continue
		}
		newindex[k] = &FileMetaData{Filename: k, Version: v.Version + 1, BlockHashList: []string{TOMBSTONE_HASHVALUE}}
	}
	return newindex, nil
}

func syncRemote(oldindex map[string]*FileMetaData, client RPCClient) (newindex map[string]*FileMetaData, bsAddrs []string, err error) {
	newindex = make(map[string]*FileMetaData)

	if len(client.BsAddrs) != 0 {
		bsAddrs = client.BsAddrs
	} else {
		err = client.GetDefualtBlockStoreAddrs(&bsAddrs)
		// log.Println("bsaddr:", bsAddr)
		if err != nil {
			log.Println("Get Default BS address failed")
			return nil, nil, err
		}
	}
	remoteindex := make(map[string]*FileMetaData)
	err = client.GetFileInfoMap(&remoteindex)
	if err != nil {
		log.Println("Get remote index failed")
		return nil, nil, err
	}
	for name, remotedata := range remoteindex {
		if localdata, ok := oldindex[name]; ok {
			if localdata.Version < remotedata.Version {
				newindex[name] = remotedata
				if err := writeFile(client, remotedata, bsAddrs); err != nil {
					log.Println("write new version file failed")
					return nil, nil, err
				}
			} else if localdata.Version > remotedata.Version {
				if err := updateFile(client, localdata, newindex, bsAddrs); err != nil {
					log.Println("update new version file failed")
					return nil, nil, err
				}
			} else if !CompareStrings(localdata.BlockHashList, remotedata.BlockHashList) {
				newindex[name] = remotedata
				if err := writeFile(client, remotedata, bsAddrs); err != nil {
					log.Println("write quick version file failed")
					return nil, nil, err
				}
			} else {
				newindex[name] = localdata
			}
		} else {
			newindex[name] = remotedata
			if err := writeFile(client, remotedata, bsAddrs); err != nil {
				log.Println("write new file failed")
				return nil, nil, err
			}
		}
	}
	for name, localdata := range oldindex {
		if _, ok := remoteindex[name]; ok {
			continue
		}
		if err := updateFile(client, localdata, newindex, bsAddrs); err != nil {
			log.Println("update new file failed")
			return nil, nil, err
		}
	}
	return newindex, bsAddrs, nil
}

func writeFile(client RPCClient, data *FileMetaData, bsAddrs []string) error {
	if len(data.BlockHashList) == 1 && data.BlockHashList[0] == TOMBSTONE_HASHVALUE {
		e := os.Remove(ConcatPath(client.BaseDir, data.Filename))
		if e != nil {
			return e
		}
		return nil
	}
	buf := []byte{}
	var block Block

	for _, hash := range data.BlockHashList {
		addr := ""
		err := client.GetBlockStoreAddr(&BlockHash{Hash: hash}, bsAddrs[0], &addr)
		if err != nil {
			return err
		}
		bsAddrs = append([]string{addr}, bsAddrs[:len(bsAddrs)-1]...)
		err = client.GetBlock(hash, addr, &block)
		if err != nil {
			return err
		}
		buf = append(buf, block.BlockData[:block.BlockSize]...)
	}

	if _, err := os.Stat(ConcatPath(client.BaseDir, data.Filename)); err == nil {
		e := os.Remove(ConcatPath(client.BaseDir, data.Filename))
		if e != nil {
			return e
		}
	}
	err := ioutil.WriteFile(ConcatPath(client.BaseDir, data.Filename), buf, 0666)
	if err != nil {
		return err
	}
	return nil
}

func updateFile(client RPCClient, localdata *FileMetaData, newindex map[string]*FileMetaData, bsAddrs []string) error {
	var v int32
	name := localdata.Filename
	err := client.UpdateFile(localdata, &v)
	if err != nil {
		log.Println("UpdateFile: client update error!")
		return err
	}
	if v == -1 {
		// update fail
		newinfomap := make(map[string]*FileMetaData)
		err := client.GetFileInfoMap(&newinfomap)
		if err != nil {
			log.Println("UpdateFile: get new map failed when update failed")
			return err
		}
		newindex[name] = newinfomap[name]
		if err := writeFile(client, newinfomap[name], bsAddrs); err != nil {
			log.Println("UpdateFile: write file failed when update failed")
			return err
		}
	} else {
		var suss Success
		newindex[name] = localdata
		if len(localdata.BlockHashList) == 1 && localdata.BlockHashList[0] == TOMBSTONE_HASHVALUE {
			return nil
		}
		fi, err := os.Open(ConcatPath(client.BaseDir, name))
		defer fi.Close()
		if err != nil {
			log.Println("UpdateFile: file failed to open")
			return err
		}
		buf := make([]byte, client.BlockSize)
		blockMap := make(map[string]*Block)
		hashes := []string{}
		for {
			n, err := fi.Read(buf)
			if err != nil && err != io.EOF {
				log.Println("UpdateFile: file read failed")
				return err
			} else {
				if err == io.EOF {
					break
				}
				hash := GetHash(buf[:n])
				blockMap[hash] = &Block{BlockData: make([]byte, n), BlockSize: int32(n)}
				copy(blockMap[hash].BlockData, buf[:n])
				hashes = append(hashes, hash)
			}
		}
		for hash, data := range blockMap {
			addr := ""
			err := client.GetBlockStoreAddr(&BlockHash{Hash: hash}, bsAddrs[0], &addr)
			if err != nil {
				log.Println("UpdateFile: failed to get bs addr")
				return err
			}
			bsAddrs = append([]string{addr}, bsAddrs[:len(bsAddrs)-1]...)
			err = client.PutBlock(data, addr, &suss.Flag)
			if err != nil {
				log.Println("UpdateFile:failed to put the block online")
				return err
			}
		}
	}
	return nil
}
