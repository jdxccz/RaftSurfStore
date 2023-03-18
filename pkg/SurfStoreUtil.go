package pkg

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Pair struct {
	Key   int
	Value string
}

func GetChordPosition(N int, hash string) (int, error) {
	digits := int(math.Log2(float64(N))/4) + 1
	hashTail := hash[len(hash)-digits:]
	i, err := strconv.ParseInt(hashTail, 16, 64)
	if err != nil {
		return -1, err
	}
	return int(i) % N, nil
}

func GetHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func Monitor(m *MetaStore) {
	for {
		log.Println("Id:", m.serverId, "Status:", m.status, "Term:", m.term)
		time.Sleep(1 * time.Second)
	}
}

// ReadLog

// WriteLog

/* File Path Related */
func ConcatPath(baseDir, fileDir string) string {
	return baseDir + "/" + fileDir
}

const createTable string = `create table if not exists indexes (
		fileName TEXT, 
		version INT,
		hashIndex INT,
		hashValue TEXT
	);`

const insertTuple string = `INSERT INTO indexes VALUES (?,?,?,?);`

// WriteMetaFile writes the file meta map back to local metadata file index.db
func WriteMetaFile(fileMetas map[string]*FileMetaData, baseDir string) error {
	// remove index.db file if it exists
	outputMetaPath := ConcatPath(baseDir, DEFAULT_META_FILENAME)
	if _, err := os.Stat(outputMetaPath); err == nil {
		e := os.Remove(outputMetaPath)
		if e != nil {
			log.Fatal("Error During Meta Write Back")
		}
	}
	db, err := sql.Open("sqlite3", outputMetaPath)
	if err != nil {
		log.Fatal("Error During Meta Write Back")
	}
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal("Error During Meta Write Back")
	}
	statement.Exec()
	for _, v := range fileMetas {
		for i := 0; i < len(v.BlockHashList); i++ {
			// values := []string{v.Filename, strconv.Itoa(int(v.Version)), strconv.Itoa(int(i)), v.BlockHashList[i]}
			// stmt := insertTuple + strings.Join(values, ",") + ");"
			_, err := db.Exec(insertTuple, v.Filename, strconv.Itoa(int(v.Version)), strconv.Itoa(int(i)), v.BlockHashList[i])
			if err != nil {
				log.Println("insert file index occur error:", err)
				return err
			}
		}
	}
	log.Println("Write Meta Files successfully!")
	return nil
}

const getDistinctFileName string = `SELECT DISTINCT fileName FROM indexes;`

const getTuplesByFileName string = `SELECT version,hashIndex,hashValue FROM indexes WHERE fileName = ?;`

// LoadMetaFromMetaFile loads the local metadata file into a file meta map.
func LoadMetaFromMetaFile(baseDir string) (fileMetaMap map[string]*FileMetaData, e error) {
	metaFilePath, _ := filepath.Abs(ConcatPath(baseDir, DEFAULT_META_FILENAME))
	fileMetaMap = make(map[string]*FileMetaData)
	metaFileStats, e := os.Stat(metaFilePath)
	if e != nil || metaFileStats.IsDir() {
		return fileMetaMap, nil
	}
	db, err := sql.Open("sqlite3", metaFilePath)
	if err != nil {
		log.Fatal("Error When Opening Meta")
	}
	Rows, err := db.Query(getDistinctFileName)
	defer Rows.Close()
	if err != nil {
		log.Println("Get Names occur error:", err)
		return nil, err
	}
	fileNames := []string{}
	for Rows.Next() {
		var filename string
		if err := Rows.Scan(&filename); err != nil {
			log.Println("San Names occur error:", err)
			return nil, err
		}
		fileNames = append(fileNames, filename)
	}
	for _, fileName := range fileNames {
		v := -999
		hashmap := make(map[int]string)
		var idx int
		var hash string
		hashblocks := []string{}
		Rows, err = db.Query(getTuplesByFileName, fileName)
		if err != nil {
			log.Println("Get Features occur error:", err)
			return nil, err
		}
		for Rows.Next() {
			if err := Rows.Scan(&v, &idx, &hash); err != nil {
				log.Println("San Features occur error:", err)
				return nil, err
			}
			hashmap[idx] = hash
		}
		var keys []int
		for k := range hashmap {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			hashblocks = append(hashblocks, hashmap[k])
		}
		fileMetaMap[fileName] = &FileMetaData{
			Filename:      fileName,
			Version:       int32(v),
			BlockHashList: hashblocks,
		}
	}
	log.Println("Load Meta Files successfully!")
	return fileMetaMap, nil
}

// PrintMetaMap prints the contents of the metadata map.
func PrintMetaMap(metaMap map[string]*FileMetaData) {

	fmt.Println("--------BEGIN PRINT MAP--------")

	for _, filemeta := range metaMap {
		fmt.Println("\t", filemeta.Filename, filemeta.Version)
		for _, blockHash := range filemeta.BlockHashList {
			fmt.Println("\t", blockHash)
		}
	}

	fmt.Println("---------END PRINT MAP--------")

}

func CompareStrings(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
