package keygen

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

const ACString = `{"K":"%s", "N":"%s", "O":"%s", "DI":"%s", "T":%d}`

type RequestCode struct {
	K  string `json:"K"`
	P  string `json:"P"`
	DI string `json:"DI"`
}

type ActivationCode struct {
	K  string `json:"K"`
	N  string `json:"N"`
	O  string `json:"O"`
	DI string `json:"DI"`
	T  int    `json:"T"`
}

type Keygen struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	requestCode    RequestCode
	activationCode ActivationCode
}

func loadPublicKey(filepath string) *rsa.PublicKey {
	block, _ := pem.Decode(loadFile(filepath))
	if block == nil {
		fmt.Println("加载公钥失败")
		os.Exit(1)
	}
	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("加载公钥失败", err)
		os.Exit(1)
	}
	pub := public.(*rsa.PublicKey)
	return pub
}

func loadPrivateKey(filepath string) *rsa.PrivateKey {
	block, _ := pem.Decode(loadFile(filepath))
	if block == nil {
		fmt.Println("加载私钥失败")
		os.Exit(1)
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("加载私钥失败", err)
		os.Exit(1)
	}
	pri := private.(*rsa.PrivateKey)
	return pri
}

func loadFile(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(filepath, "文件不存在")
		os.Exit(1)
	}
	defer file.Close()

	key := make([]byte, 2048)
	num, err := file.Read(key)

	return key[:num]
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func privateEncrypt(pri *rsa.PrivateKey, data []byte) ([]byte, error) {
	signData, err := rsa.SignPKCS1v15(nil, pri, crypto.Hash(0), data)
	if err != nil {
		return nil, err
	}
	return signData, nil
}

func NewKeygen(publicFile, privateFile, name, organize string, time int) *Keygen {
	return &Keygen{
		privateKey:  loadPrivateKey(privateFile),
		publicKey:   loadPublicKey(publicFile),
		requestCode: RequestCode{},
		activationCode: ActivationCode{
			K:  "",
			DI: "",
			N:  name,
			O:  organize,
			T:  time,
		},
	}
}

func (k *Keygen) GetActivationCode(code string) {
	k.decrypt(code)
	k.activationCode.DI = k.requestCode.DI
	k.activationCode.K = k.requestCode.K

	b := fmt.Sprintf(ACString, k.activationCode.K, k.activationCode.N, k.activationCode.O, k.activationCode.DI, k.activationCode.T)

	partLen := k.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(b), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bts, err := privateEncrypt(k.privateKey, chunk)
		if err != nil {
			return
		}
		buffer.Write(bts)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("***************** 激 活 码 ******************")
	fmt.Print( base64.StdEncoding.EncodeToString(buffer.Bytes()), "\r\n")

	fmt.Println("********************************************")
	fmt.Println("*           Navicat Keygen End             *")
	fmt.Println("*              version 1.0                 *")
	fmt.Println("********************************************")
}

func (k *Keygen) decrypt(code string) {
	data, _ := base64.StdEncoding.DecodeString(code)
	b, err := rsa.DecryptPKCS1v15(rand.Reader, k.privateKey, data)
	if err != nil {
		fmt.Println("解密失败")
		os.Exit(1)
	}
	_ = json.Unmarshal(b, &k.requestCode)
	fmt.Println("***************** 解 密 成 功 ****************")
	fmt.Print(string(b), "\r\n")
	fmt.Println("***************** 解 密 成 功 ****************")
}
