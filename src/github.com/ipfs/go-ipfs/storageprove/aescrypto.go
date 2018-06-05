package storageprove
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	//"fmt"
	"io"
	"strings"
//home/long/dev/ipfs/go-ipfs/src/gx/ipfs/QmQPWTeQJnJE7MYu6dJTiNTQRNuqBr41dis6UgY6Uekmgd/keccakpg
	keccak "gx/ipfs/QmQPWTeQJnJE7MYu6dJTiNTQRNuqBr41dis6UgY6Uekmgd/keccakpg"
	"gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
	//"os"
	//"bufio"
)

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}


func Encrypt(key1 []byte, text1 []byte) ([]byte, string, error) {   //return cyphertext and hash
	key := []byte("LKHlhb899Y09olUi")
	text := string(text1)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{0}, "", err
	}

	msg := Pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{0}, "", err
	}
	for i,_ := range iv {
		iv[i] = 0
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))


	h := keccak.New256()
	h.Write([]byte(finalMsg))
	hash := h.Sum(nil)

	return []byte(finalMsg), base32.RawStdEncoding.EncodeToString(hash), nil
}

func Decrypt(key1 []byte, text1 []byte) ([]byte, error) {
	key := []byte("LKHlhb899Y09olUi")
	text := string(text1)
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{0}, err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return []byte{0}, err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return []byte{0}, errors.New("blocksize must be multipe of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return []byte{0}, err
	}

	return unpadMsg, nil
}


func PKCS5Padding(ciphertext []byte, length int, blockSize int) ([]byte, int) {
	padding := blockSize - length%blockSize
	padtext :=make([]byte, length+padding)
	copy(padtext, ciphertext)
	for  i:=0; i<padding; i++{
		padtext[length + i] = byte(padding)
	}

	return padtext, padding
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(key1 []byte, buf []byte) ([]byte, string, error) { //return cyphertext and hash
	key := GetAesKey()//[]byte("LKHlhb899Y09olUi")
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{0}, "", err
	}
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
	var blockMode cipher.BlockMode

		blockMode = cipher.NewCBCEncrypter(block,iv)


	//buf := make([]byte,1024)  //AES BlockSize = 16
	var buf_out []byte

	n := len(buf)

	buf1, padLength := PKCS5Padding(buf,n,16)
	buf_out = make([]byte,n + padLength)
	blockMode.CryptBlocks(buf_out, buf1[:(n + padLength)])


	h := keccak.New256()
	h.Write(buf_out)
	hash := h.Sum(nil)

	return buf_out, base32.RawStdEncoding.EncodeToString(hash), nil
}


func AesDecrypt(key1 []byte, buf []byte) ([]byte, error) { //return cyphertext and hash
	key := GetAesKey();//[]byte("LKHlhb899Y09olUi")
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{0}, err
	}
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
	var blockMode cipher.BlockMode

	blockMode = cipher.NewCBCDecrypter(block,iv)
	n := len(buf)
	buf_out := make([]byte,n)
	blockMode.CryptBlocks(buf_out, buf[:n])
	ntmp :=n
	// get the next block

	// end file
	buf_out = PKCS5UnPadding(buf_out[:ntmp])
	return buf_out, nil

}
/*
func AesBlockCipher(keydata []byte, flag int, buf []byte) ([]byte, error){

	// enc file

	block, err := aes.NewCipher(keydata)
	if err != nil {
		return []byte{0}, err
	}
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
	var blockMode cipher.BlockMode
	if flag==1{
		blockMode = cipher.NewCBCEncrypter(block,iv)
	}else{
		blockMode = cipher.NewCBCDecrypter(block,iv)
	}

	//buf := make([]byte,1024)  //AES BlockSize = 16
	//buf_out := make([]byte,1024)
	n := len(buf)

		if flag ==1 { //enc


			if (n % 16)!=0 {
				//padding
				padLength := PKCS5Padding(buf,n,16)
				buf_out := make([]byte,n + padLength)
				blockMode.CryptBlocks(buf_out, buf[:(n + padLength)])
			}else{
				buf_out := make([]byte,n)
				blockMode.CryptBlocks(buf_out, buf)
			}
		}else{  //dec
			buf_out := make([]byte,n)
			blockMode.CryptBlocks(buf_out, buf[:n])
			ntmp :=n
			// get the next block

				// end file
			buf_out = PKCS5UnPadding(buf_out[:ntmp])
		}


	return buf_out, nil
}
*/

