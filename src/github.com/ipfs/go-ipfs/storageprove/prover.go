package storageprove

import (
    //"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipfs/repo"
	"gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore"
	"fmt"
	"time"
	//"os"
	//"io"
	//"gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore/fs"
	"sync"
	"os"
	"bufio"
	"errors"
	"io"

	"math/rand"
	"gx/ipfs/QmdYwCmx8pZRkzdcd8MhmLJqYVoVTC1aGsy5Q4reMGLNLg/atomicfile"
	"path/filepath"
	"github.com/ipfs/go-ipfs/repo/config"
	"net"
	"encoding/json"
	"strconv"
	"gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
	"strings"
	"gx/ipfs/Qmc74pRHvndTDAB5nXztWAV7vs5j1obvCb9ejfQzXp9USX/retry-datastore"
	"syscall"
)

var ds *retry.Datastore
var lastKey datastore.Key
var aesLock *sync.Mutex
var countLock *sync.Mutex


var prefixDs = string("/pos/")
var pathTypeVersionDs = string("/p1")
var hashTypeVersionDs = string("/h1")
var countDs = string("/pos/c/count")

type prover struct{
    lock sync.Locker
}

func GetAesKey() []byte{
	aesLock.Lock()
	defer aesLock.Unlock()
	key := []byte("LKHlhb899Y09olUi")
	return key
}

func checkError(err error){
	if err != nil {
		fmt.Println("error is: ", err)
	}
}

func InitProver(filePath string) error {
	cf, err := atomicfile.New(filepath.Join(filePath, "count"), 0660)
	if err != nil {
		return err
	}
	line := strconv.Itoa(0)
	if _, err := cf.Write([]byte(line + "\n")); err != nil {
		fmt.Println("failed to write count file ")
		return err
	}
	defer cf.Close()
	return nil
}


func GetProverFile() (string, error) {
	var pathRoot string
	var err error
	if pathRoot, err = config.PathRoot(); err != nil {
		return "", err
    }

    return filepath.Join(pathRoot, "prover"), nil
}

func getCountFile() (string, error) {
	var pathRoot string
	var err error
	if pathRoot, err = config.PathRoot(); err != nil {
		return "", err
	}

	return filepath.Join(pathRoot, "count"), nil
}

func getCount() (int, error) {
	//var bCount []byte
	//var err error

	key := datastore.NewKey(countDs)
	bCount, err := ds.Get(key)
	//TODO: count should deal with care
	if err == datastore.ErrNotFound {
		fmt.Errorf("storageprove: err get count:", err)
		err = ds.Put(key, []byte(strconv.Itoa(0)))
		if err != nil {
			return -1, err
		}
		return 0, nil
	}
	if  err != nil {
		fmt.Errorf("storageprove: err get count:", err)
		return 0, err
	}
	return strconv.Atoi(string(bCount.([]byte)))
/*
	var countFile, line string
	var err error
	if countFile, err = getCountFile(); err != nil {
		return 0, err
	}
	if line, err = rsl(countFile, 1); err != nil {
		fmt.Println("failed to read count file ", err)
		return 0, err
	}
	return strconv.Atoi(line[:len(line)-1])
*/
}

func increaseCount() (error) {


	count, err := getCount()
	if err != nil {
		fmt.Println("storageprove: err get count", err)
		return err
	}
	count ++
	key := datastore.NewKey(countDs)
	err = ds.Put(key, []byte(strconv.Itoa(count)))
	if err != nil {
		fmt.Println("storageprove: err put count", err)
		return err
	}
	return nil
	/*
	var countFile, line string
	var err error
	var cf *os.File
	var count int
	if countFile, err = getCountFile(); err != nil {
		return err
	}
	if line, err = rsl(countFile, 1); err != nil {
		fmt.Println("failed to read count file ", err)
		return err
	}
	fmt.Println("count: ", line)
	if count,err = strconv.Atoi(line[:len(line)-1]); err != nil {
		fmt.Println("count file not int", err)
	}
	cf, err = os.OpenFile(countFile, os.O_RDWR, 0660);
	if err != nil {
		fmt.Println("failed to open count file ", err)
		return err
	}
	count ++
	line = strconv.Itoa(count)
	if _, err := cf.Write([]byte(line + "\n")); err != nil {
		fmt.Println("failed to write count file ")
		return err
	}
	if err := cf.Close(); err != nil {
		return err
	}
	return nil
	*/
}

func PutRecord(path string, hash string, sync bool) error {
	var count int
	var err error

	countLock.Lock()
	defer countLock.Unlock()

	if count, err = getCount(); err != nil {
		fmt.Println("storageprove: getCount failed", err)
		return err
	}

	pathKey := datastore.NewKey(prefixDs + strconv.Itoa(count) + pathTypeVersionDs)
	Put(pathKey, []byte(path))

	hashKey := datastore.NewKey(prefixDs + strconv.Itoa(count) + hashTypeVersionDs)
	Put(hashKey, []byte(hash))

	if err := increaseCount(); err != nil {
		fmt.Println("storageprove: increaseCount failed", err)
		return err
	}

	return nil
}

func isTrue(err error) bool {
	return true
}

func isTooManyFDError(err error) bool {
	perr, ok := err.(*os.PathError)
	if ok && perr.Err == syscall.EMFILE {
		return true
	}

	return false
}

func initRetryDataStore(ds repo.Datastore) *retry.Datastore {
	rds := &retry.Datastore{
	Batching:    ds,
	Delay:       time.Millisecond * 200,
	Retries:     6,
	TempErrFunc: isTooManyFDError,
	}
	return rds
}
var Crypt AesCryptI
var Ra *AesCrypt
func SetDataStore(d repo.Datastore) {

	ds = initRetryDataStore(d)		//measure datastore
	aesLock = new(sync.Mutex)
	countLock = new(sync.Mutex)

	useSdsc := false
	if (useSdsc) {
		Crypt = NewSdsc("E")// windows
		//Crypt = NewSdsc("/dev/sdb") // linux
	} else {
		Crypt = NewAes()
	}
	Ra = &AesCrypt {
		AesCryptI:   Crypt,
		Delay:       time.Millisecond * 200,
		Retries:     6,
		TempErrFunc: isTrue,
	}


	type proveData struct {
		//Id string `json:"id"`
		Count int `json:"count"`
	}

	go func() {
		count := int(500)
		var buf []byte
		var err error
		//fmt.Println("Daemon timer started")
		for {
			select {
			case <-time.After(time.Hour * 1 ):
				//count += 5
				//fmt.Println("daemon live for seconds: ", count)
				//fmt.Println("pos check ")
				//*
				if count, err = check(); err != nil {
					//TODO: send fail msg to server?
					//fmt.Println("pos failed!")
				}
				//*/
				//Get(lastKey)
				//*
				//fmt.Println("pos succeed: count ", count)
				pd := proveData{Count:count}
				if buf, err= json.Marshal(pd); err!=nil{
					//fmt.Println("marshel err: ", err)
				}
				//fmt.Println("pd: ",buf)
				service := "127.0.0.1:17888" //os.Args[1]
				//tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
				//checkError(err)
				conn, err := net.Dial("tcp", service)
				//checkError(err)
				if err != nil {
					continue;
				}
				_, err = conn.Write(buf)
				//checkError(err)
				//result, err := ioutil.ReadAll(conn)
				//checkError(err)
				//fmt.Println(string(result))
				//*/
			}
		}
	}()
}

func check() (int, error){
	aesLock.Lock()
	defer aesLock.Unlock()
	var line string
	var path, hash string
	var err error
	var lenRead, countTotal int
	var hash2 []byte

	//var bHash, bPath []byte

	if countTotal, err = getCount(); err != nil {
		fmt.Println("storageprove: failed to get count", err)
		return 0, err
	}

	if countTotal == 0 {
		fmt.Println("storageprove: not init?")
		return 0, err
	}

	rand.Seed(time.Now().Unix())
	count := rand.Intn(countTotal-1)

//for count = 0; count < countTotal; count ++ {

	fmt.Printf("checking %d, total %d\n", count, countTotal)
	pathKey := datastore.NewKey(prefixDs + strconv.Itoa(count) + pathTypeVersionDs)
	if path, err = Get(pathKey); err != nil {
		fmt.Println("storageprove: failed to get key:", pathKey, err)
		return 0, err
	}

	hashKey := datastore.NewKey(prefixDs + strconv.Itoa(count) + hashTypeVersionDs)
	if hash, err = Get(hashKey); err != nil {
		fmt.Println("storageprove: failed to get key:", hashKey, err)
		return 0, err
	}



	rec, err := os.OpenFile(path, os.O_RDWR, 0660);
	//rec, err := ioutil.TempFile("/home/long", "temp")
	if err != nil {
		fmt.Println("storageprove: failed to open file: ", line, err)
		return 0, err
	}
	defer rec.Close()
	rd := make([]byte, 256*1024+16)
	rec.Seek(0, 0)
	if lenRead, err = rec.Read(rd); err != nil && err != io.EOF {
		fmt.Println("storageprove: not EOF:", err)
		return 0, err
	}
	//fmt.Println("rd", rd)
	rd = rd[:lenRead]
	if hash2, err = Crypt.Hash(rd); err != nil {
		fmt.Println("storageprove: hash failed:", err)
		return 0, err
	}

	strHash := base32.RawStdEncoding.EncodeToString(hash2)
	//fmt.Println("hash & len: ", strHash, len(strHash))

	if strings.Compare(hash, strHash) != 0 {
		fmt.Println("hash cmp failed: ", hash, strHash)
		return 0, errors.New("storageprove: pos failed to compare")
	}
//}

    return count, nil
}


func rsl(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}


func SetlastKey(key datastore.Key) {
	//fmt.Println("lastkey:", key)
	//ds.NewKey()
	for i:=0; i<1; i++ {
		//Put(key, 0)
	}
	//Get(key)
}

func Put(key datastore.Key, value interface{}) (error) {
	//key1 := datastore.NewKey("/pos/s/haha")
	//fmt.Println("loong puting", key.String(), string(value.([]byte)))
	err := ds.Put(key, value)
	if err != nil {
		fmt.Errorf("storageprove: error put key, %s", key.String(), err)
		panic("put error!")
		return err
	}
	_, err = Get(key)

	return err
}

func Get(key datastore.Key) (string, error) {
	//fmt.Println("loong geting", key.String())
	data, err := ds.Get(key)
	if err != nil {
		fmt.Errorf("storageprove: error get key: %s", key.String(), err)
		//return "", err
	}

	//fmt.Println("loong get",string(data.([]byte)))
	return string(data.([]byte)), err
}

