package storageprove

import (
    //"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipfs/repo"
	"gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore"
	keccak "gx/ipfs/QmQPWTeQJnJE7MYu6dJTiNTQRNuqBr41dis6UgY6Uekmgd/keccakpg"
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
)
var ds repo.Datastore
var lastKey datastore.Key
var aesLock *sync.Mutex

func SetlastKey(key datastore.Key) {
	//lastKey = key
}

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

func SetDataStore(d repo.Datastore) {
	ds = d
	aesLock = new(sync.Mutex)
	go func() {
		count := 0
		fmt.Println("Daemon timer started")
		for {
			select {
			case <-time.After(time.Second * 5):
				count += 5
				//fmt.Println("daemon live for seconds: ", count)
				//check()
				//Get(lastKey)
				/*
				service := "127.0.0.1:17888" //os.Args[1]
				//tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
				//checkError(err)
				conn, err := net.Dial("tcp", service)
				checkError(err)
				_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
				checkError(err)
				result, err := ioutil.ReadAll(conn)
				checkError(err)
				fmt.Println(string(result))
				//*/
			}
		}
	}()
}

func check(){
	aesLock.Lock()
	defer aesLock.Unlock()
	const content = "/home/long/tmp/foo.txt"
	var line string
	var err error
	var lenRead int
	line2 := "/home/long/tmp/ipfstest3/blocks/B5/CIQDOZU3EAGXWK3PLVFOFOZOAE5USX3XM6I5CSHSQGTML2BAGN7MB5I.data"
	rand.Seed(time.Now().Unix())
	lineNum := rand.Intn(1000)
	lineNum = lineNum - lineNum%2 +1
	if line, err = rsl(content, lineNum); err == nil {
		fmt.Println(" line:", lineNum)
		fmt.Println(line, len(line), len(line2), line[:len(line)-1])
	} else {
		fmt.Println("rsl:", err)
	}

	rec, err := os.OpenFile(line[:len(line)-1], os.O_RDWR, 0660);
	//rec, err := ioutil.TempFile("/home/long", "temp")
	if err != nil {
		fmt.Println("failed to open file: ", line, err)
	}
	rd := make([]byte, 256*1024 + 16)
	rec.Seek(0, 0)
	if lenRead, err = rec.Read(rd); err != nil && err != io.EOF {
		fmt.Println("not EOF:", err)
	}
	fmt.Println("rd", rd)
	rd = rd[:lenRead]
	h := keccak.New256()
	h.Write(rd)
	hash := h.Sum(nil)
	fmt.Println("hash: ", hash)

	if line, err = rsl(content, lineNum+1); err == nil {
		fmt.Println("line:", lineNum+1)
		fmt.Println(line, len(line), len(line2), line[:len(line)-1])
	} else {
		fmt.Println("rsl:", err)
	}

	if err := rec.Close(); err != nil {
		return
	}

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



func Get(key datastore.Key) (value interface{}, err error) {
	data, err := ds.Get(key)
	if err != nil {
		fmt.Println("storageprove: error get key")
	}
	fmt.Println(data)
	return data, err
}

