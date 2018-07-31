package storageprove

import (
	"fmt"
	"time"
)

type AesCryptI interface {

	AesEncrypt(key1 []byte, buf []byte) ([]byte, string, error)

	AesDecrypt(key1 []byte, buf []byte) ([]byte, error)

	Hash(buf []byte) ([]byte, error)
}

type AesCrypt struct {
	TempErrFunc func(error) bool
	Retries     int
	Delay       time.Duration

	AesCryptI
}

var errFmtString = "ran out of retries trying to get past temporary error: %s"

func (a *AesCrypt) runOp(op func() error) error {
	err := op()
	if err == nil || !a.TempErrFunc(err) {
		return err
	}

	for i := 0; i < a.Retries; i++ {
		time.Sleep(time.Duration(i+1) * a.Delay)

		err = op()
		if err == nil || !a.TempErrFunc(err) {
			return err
		}
	}

	return fmt.Errorf(errFmtString, err)
}

func (a *AesCrypt) AesEncrypt(key1 []byte, buf []byte) ([]byte, string, error) {
	var cypher []byte
	var hash string
	err := a.runOp(func() error {
		var err error
		cypher, hash, err = a.AesCryptI.AesEncrypt(key1, buf)
		return err
	})

	return cypher, hash, err
}

func (a *AesCrypt) AesDecrypt(key1 []byte, buf []byte) ([]byte, error) {
	var plain []byte
	err := a.runOp(func() error {
		var err error
		plain, err = a.AesCryptI.AesDecrypt(key1, buf)
		return err
	})

	return plain, err
}

func (a *AesCrypt) Hash(buf []byte) ([]byte, error) {
	var hash []byte
	err := a.runOp(func() error {
		var err error
		hash, err = a.AesCryptI.Hash(buf)
		return err
	})
	return hash, err
}
