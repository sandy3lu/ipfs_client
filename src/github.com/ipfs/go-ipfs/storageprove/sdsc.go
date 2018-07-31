package storageprove

import (
	"fmt"
	"errors"

	"syscall"
	"unsafe"
	"gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
)

var ErrSdscCom = errors.New("sdsc: access failed")

var (
	dll32, _   = syscall.LoadDLL("LKT4209H_64.dll")

	c_Open, _     = dll32.FindProc("EK_Open")
	c_Reset, _    = dll32.FindProc("EK_Reset")
	c_APDU, _     = dll32.FindProc("EK_Exchange_APDU")
	c_Close, _         = dll32.FindProc("EK_Close")
)

type AesSdsc struct {
	AesCryptI
}

func (a *AesSdsc) AesEncrypt(key1 []byte, buf []byte) ([]byte, string, error) {
	return AesEncrypt(key1, buf)
}

func (a *AesSdsc) AesDecrypt(key1 []byte, buf []byte) ([]byte, error) {
	return AesDecrypt(key1, buf)
}

func (a *AesSdsc) Hash(buf []byte) ([]byte, error) {
	return Hash(buf)
}

func NewSdsc(dev string) *AesSdsc {
	Sdsc_Connect(dev)
	a := &AesSdsc{
	}
	return a
}

func Sdsc_Connect(title string) {

	diskName:= []byte(title)
	ret, _, _ := c_Open.Call(uintptr((uint)(diskName[0])))
	fmt.Println("connect: ret ", ret)
}

func Sdsc_GetATR(title string, l int) {

	ttitle := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	ret, _, _ := c_Reset.Call(uintptr(unsafe.Pointer(&l)),uintptr(unsafe.Pointer(&ttitle[0])))

	fmt.Printf("atr: %x, len %d title: %s, ret %d\n", ttitle, l, title, ret)
}
/*
func Sdsc_SendAPDU(CMD string, CMDlen int, Res string, Reslen int) {
	cCMD := C.CString(CMD)
	cCMDlen := C.uint(CMDlen)
	clen := C.uint(Reslen)
	cRes := C.CString(Res)
	C.SDSmartCard_SendAPDU2(cCMD, cCMDlen, cRes, &clen)
	fmt.Println("get res: ", cRes, clen)
}
*/
func Sdsc_SendAPDUA(CMD []byte, CMDlen int, Res []byte, Reslen int, Time int)  (uint, error) {
	ret, _, _:= c_APDU.Call(uintptr(len(CMD)),uintptr(unsafe.Pointer(&CMD[0])),uintptr(unsafe.Pointer(&Reslen)),uintptr(unsafe.Pointer(&Res[0])))

	fmt.Printf("get res: len %d, ret %d, res %x, \n", Reslen, ret, Res)
//	if ret != 0{  //TODO:ret is not 0 when return
	////		fmt.Println("sdsc: send apdua error, ret ", ret)
	////		return 0, ErrSdscCom
	////	}
	return uint(Reslen), nil
}

//return cyphertext, the SW(0x9000) not included
//caller must make sure plainLen >= 16
func SdscAesEncrypt(plain []byte, plainLen int) ([]byte, uint, error) {

	commandHdr := []byte{0x80, 0x08, 0x00, 0x00, byte(plainLen + 2), 0x07, byte(plainLen)}
	commandHdr = append(commandHdr, plain...)

	cCypher := make([]byte, plainLen + 2 + 2)	//0x9000 + aabb
	respLen, err := Sdsc_SendAPDUA(commandHdr, len(commandHdr), cCypher, 0, 10)
	if err != nil{
		fmt.Println("sdsc error, ret ", err)
		return nil, 0, err
	}
	cypher := make([]byte, respLen - 2 - 2)	//0x9000 + aabb
	for i := uint(0); i < respLen - 2 -2; i++ {
		cypher[i] = byte(cCypher[i + 2])
	}
	fmt.Printf("sdsc aes encrypt result:len %d,  %x, \n", respLen, cypher)
	return cypher, respLen - 2 - 2, nil
}

func AesEncrypt(key1 []byte, buf []byte) ([]byte, string, error) { //return cyphertext and hash(string)
	aesLock.Lock()
	defer aesLock.Unlock()
	if len(buf) < 16 {
		return buf, "", nil
	}
	bufLen := len(buf)
	if bufLen > cyhperHdrLen {
		bufLen = cyhperHdrLen
	}

	bufLen = bufLen - (bufLen % 16)
	//only head cyhperHdrLen(bytes) is used because of sdsc's performance limitation
	buf1 := buf[:bufLen]
	cypher, _, err := SdscAesEncrypt(buf1, len(buf1))
	if err != nil {
		fmt.Println("SdscAesEncrypt fail", err)
		return []byte{}, "", err
	}

	hash, _, err := SdscHash(1, cypher, len(cypher))
	if err != nil {
		fmt.Println("SdscHash fail", err)
		return []byte{}, "", err
	}

	//buf_out := make([]byte, len(buf))
	cypher = append(cypher, buf[bufLen:]...)

	return cypher, base32.RawStdEncoding.EncodeToString(hash), nil
}

//return plaintext, the SW(0x9000) not included
func SdscAesDecrypt(cypher []byte, cypherLen int) ([]byte, uint, error) {
	//80080000120510112233445566778899AABBCCDDEEFF00
	commandHdr := []byte{0x80, 0x08, 0x00, 0x00, byte(cypherLen + 2), 0x08, byte(cypherLen)}

	commandHdr = append(commandHdr, cypher...)

	cPlain := make([]byte, cypherLen + 2 + 2)	//0x9000 + aabb
	respLen, err := Sdsc_SendAPDUA(commandHdr, len(commandHdr), cPlain, 0, 10)
	if err != nil{
		fmt.Println("sdsc error, ret ", err)
		return nil, 0, err
	}
	plain := make([]byte, respLen - 2 - 2)
	for i := uint(0); i < respLen - 2 - 2; i++ {
		plain[i] = byte(cPlain[i + 2])
	}
	fmt.Printf("sdsc aes decrypt result: %x, len %d \n", cypher, respLen)
	return plain, respLen - 2 - 2, nil
}

func AesDecrypt(key1 []byte, buf []byte) ([]byte, error) { //return cyphertext and hash
	aesLock.Lock()
	defer aesLock.Unlock()
	if len(buf) < 16 {
		return buf, nil
	}

	bufLen := len(buf)
	if bufLen > cyhperHdrLen {
		bufLen = cyhperHdrLen
	}

	bufLen = bufLen - (bufLen % 16)
	//only head cyhperHdrLen(bytes) is used because of sdsc's performance limitation
	buf1 := buf[:bufLen]
	plain, _, err := SdscAesDecrypt(buf1, len(buf1))
	if err != nil {
		fmt.Println("SdscAesDecrypt fail", err)
		return []byte{}, err
	}

	//buf_out := make([]byte, len(buf))
	plain = append(plain, buf[bufLen:]...)
	return plain, nil

}

//return len includes SW(0x9000)
func SdscHash(htype byte, data []byte, lenData int) ([]byte, uint, error) {
	//80 08 00 00 07 09 00
	command := []byte{0x80, 0x08, 0x00, 0x00, byte(lenData + 2), 0x09, byte(htype)}
	command = append(command, data...)
	cHash := make([]byte, 34)	//0x9000
	var respLen uint
	var err error

	if respLen, err = Sdsc_SendAPDUA(command, len(command), cHash, 0, 2); err != nil {
		fmt.Println("sdsc error, ret ", err)
		return nil, 0, err
	}

	fmt.Printf("sdsc hash result: %x, len %d \n", cHash, respLen)
	return cHash, respLen, nil
}
/*
var maxSeq byte
var maxPerCommand byte
//return len includes SW(0x9000)
func SdscHash2(htype byte, data []byte, lenData int) ([]byte, uint, error) {
	//80 08 00 00 07 09 00
	var command []C.uchar
	var j int
	var respLen, i uint
	var err error
	cHash := make([]C.uchar, 36) //0x9000
	maxSeq = 0xFF
	maxPerCommand = 0xFC
	j = lenData / int(maxPerCommand)

	if lenData > int(maxPerCommand) * int(maxSeq -1) {
		return nil, 0, errors.New("SdscHash2: input data too long")
	}

	if( j == 0 ){
		return SdscHash(1, data, lenData)
	}
	for jj := 0; jj < j; jj ++ {
		command = []C.uchar{0x80, 0x08, 0x00, 0x00, 0xFF, 0x0A, C.uchar(jj+0), C.uchar(maxPerCommand)}
		cData := make([]C.uchar, maxPerCommand)
		for i := 0; i < int(maxPerCommand); i++ {
			cData[i] = C.uchar(data[jj * int(maxPerCommand) + i])
		}
		command = append(command, cData...)

		//fmt.Printf("cmmand: %x\n", command)
		//begin := time.Now()

		if respLen, err = Sdsc_SendAPDUA(command, len(command), cHash, 0, 1); err != nil {
			fmt.Println("sdsc error, ret ", err)
			return nil, 0, err
		}

		//end := time.Now().Nanosecond()
		//fmt.Println("hash: ", time.Since(begin))
	}
	command = []C.uchar{0x80, 0x08, 0x00, 0x00, C.uchar(lenData % int(maxPerCommand)+3), 0x0A, C.uchar(0xFF), C.uchar(lenData % int(maxPerCommand))}
	cData := make([]C.uchar, lenData % int(maxPerCommand))
	for i := 0; i < int(lenData % int(maxPerCommand)); i++ {
		cData[i] = C.uchar(data[j*int(maxPerCommand)+i])
	}
	command = append(command, cData...)
	fmt.Printf("cmmand: %x\n", command)


	//begin := time.Now()

	if respLen, err = Sdsc_SendAPDUA(command, len(command), cHash, 0, 1); err != nil {
		fmt.Println("sdsc error, ret ", err)
		return nil, 0, err
	}

	//end := time.Now().Nanosecond()
	//fmt.Println("hash: ", time.Since(begin))
	cypher := make([]byte, respLen)
	for i = 0; i < respLen; i++ {
		cypher[i] = byte(cHash[i])
	}
	fmt.Printf("sdsc hash2: %x, len %d \n", cypher, respLen)
	return cypher, respLen, nil
}
*/
func Hash(buf []byte) ([]byte, error) {
	if len(buf) < 16 {
		return []byte{}, nil
	}
	bufLen := len(buf)
	if bufLen > cyhperHdrLen {
		bufLen = cyhperHdrLen
	}

	bufLen = bufLen - (bufLen % 16)
	//only head cyhperHdrLen(bytes) is used because of sdsc's performance limitation
	buf1 := buf[:bufLen]
	hash, _, err := SdscHash(1, buf1, len(buf1))
	if err != nil {
		fmt.Println("SdscHash fail", err)
		return []byte{}, err
	}
	return hash, nil
}
