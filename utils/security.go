package utils

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}

// 从文件中读取RSA key
func RSAReadKeyFromFile(filename string) []byte {
	f, err := os.Open(filename)
	var b []byte

	if err != nil {
		return b
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b = make([]byte, fileInfo.Size())
	f.Read(b)
	return b
}

// 通过文件加载私钥对象
func LoadPrivateKey(privateKeyContent string, isPath bool) (privateKey *rsa.PrivateKey, err error) {
	var privateKeyBytes []byte
	//如果传入是私钥路径直接读取，如果传入是私钥内容则格式化成正确的私钥格式
	if isPath {
		privateKeyBytes, err = ioutil.ReadFile(privateKeyContent)
		if err != nil {
			return nil, fmt.Errorf("读取私钥证书失败 file err:%s", err.Error())
		}
	} else {
		//var publicHeader = "-----BEGIN RSA PRIVATE KEY-----\n"
		//var publicTail = "-----END RSA PRIVATE KEY-----\n"
		var publicHeader = "-----BEGIN PRIVATE KEY-----\n"
		var publicTail = "-----END PRIVATE KEY-----\n"
		var temp string
		split(privateKeyContent, &temp)
		privateKeyBytes = []byte(publicHeader + temp + publicTail)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("解码私钥失败")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败 err:%s", err.Error())
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("非法私钥文件，检查私钥文件")
	}
	return privateKey, nil
}

func LoadPublicKey(publicKeyContent string, isPath bool) (publicKey *rsa.PublicKey, err error) {
	var publicKeyBytes []byte
	//如果传入是私钥路径直接读取，如果传入是私钥内容则格式化成正确的私钥格式
	if isPath {
		publicKeyBytes, err = ioutil.ReadFile(publicKeyContent)
		if err != nil {
			return nil, fmt.Errorf("读取公钥证书失败 file err:%s", err.Error())
		}
	} else {
		if strings.Contains(publicKeyContent, "BEGIN PUBLIC KEY") {
			publicKeyBytes = []byte(publicKeyContent)
		} else {
			var publicHeader = "-----BEGIN PUBLIC KEY-----\n"
			var publicTail = "-----END PUBLIC KEY-----\n"
			var temp string
			split(publicKeyContent, &temp)
			publicKeyBytes = []byte(publicHeader + temp + publicTail)
		}

	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("解码公钥失败")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败 err:%s", err.Error())
	}
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("非法公钥文件，检查公钥文件")
	}
	return publicKey, nil
}

// 64个字符换行
func split(key string, temp *string) {
	if len(key) <= 64 {
		*temp = *temp + key + "\n"
	}
	for i := 0; i < len(key); i++ {
		if (i+1)%64 == 0 {
			*temp = *temp + key[:i+1] + "\n"
			key = key[i+1:]
			split(key, temp)
			break
		}
	}
}

// RSAEncrypt RSA加密
func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	// 解析公钥
	block, _ := pem.Decode(publicBytes)

	if block == nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确")
	}

	// 使用X509将解码之后的数据 解析出来
	// x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 将数据加密为base64格式
	return []byte(EncodeStr2Base64(string(res))), nil
}

// RSADecrypt RSA解密
func RSADecrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	// 将base64数据解析
	data := []byte(DecodeStrFromBase64(string(base64Data)))
	// 解析私钥
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确")
	}
	// 还原数据
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	return res, nil
}

// originalData 签名前的原始数据
// privateKey RSA 私钥
func SignRsa(originalData string, privateKey *rsa.PrivateKey) (string, error) {

	h := sha256.New()
	h.Write([]byte(originalData))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(signature), err
}

// originalData 签名前的原始数据
// signData Base64 格式的签名串
// pubKey 公钥（需与加密时使用的私钥相对应）
// 返回 true 代表验签通过，反之为不通过
func VerySignWithRsa(originalData, signData string, publicKey *rsa.PublicKey) (bool, error) {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return false, err
	}

	// sha256 加密方式，必须与 下面的 crypto.SHA256 对应
	// 例如使用 sha1 加密，此处应是 sha1.New()，对应 crypto.SHA1
	hash := sha256.New()
	hash.Write([]byte(originalData))
	verifyErr := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), sign)
	return verifyErr == nil, nil
}

// Aes/ECB模式的加密方法，PKCS7填充方式
func AesEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		return nil, errors.New("plaintext empty")
	}

	blockSize := block.BlockSize()
	origData := PKCS7Padding(src, blockSize)

	mode := NewECBEncrypter(block)

	crypted := make([]byte, len(origData))

	mode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Aes/ECB模式的解密方法，PKCS7填充方式
func AesDecrypt(src, key []byte) ([]byte, error) {
	Block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 确保数据块是 AES 所需的块大小的倍数
	if len(src)%Block.BlockSize() != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	if len(src) == 0 {
		return nil, errors.New("plaintext empty")
	}
	mode := NewECBDecrypter(Block)
	ciphertext := src
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext, nil
}

// ECB模式结构体
type ecb struct {
	b         cipher.Block
	blockSize int
}

// 实例化ECB对象
func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

// ECB加密类
type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}
}

// ECB解密类
type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}
}

// PKCS7填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7去除
func PKCS7UnPadding(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

// 零点填充
func ZerosPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}

// 零点去除
func ZerosUnPadding(ciphertext []byte) []byte {
	return bytes.TrimFunc(ciphertext, func(r rune) bool {
		return r == rune(0)
	})
}
