package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SQLITE_FILE_HEADER = "SQLite format 3\x00"
	IV_SIZE            = 16
	HMAC_SHA1_SIZE     = 20
	KEY_SIZE           = 32
	DEFAULT_PAGESIZE   = 4096
	DEFAULT_ITER       = 64000
)

func Decrypt(key string, filePath string, decryptedPath string) error {
	password, err := hex.DecodeString(strings.Replace(key, " ", "", -1))
	// fmt.Println(password)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	blist, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	salt := blist[:16]
	byteKey := pbkdf2.Key(password, salt, DEFAULT_ITER, KEY_SIZE, sha1.New)
	first := blist[16:DEFAULT_PAGESIZE]
	mac_salt := make([]byte, 16)
	for i := 0; i < 16; i++ {
		mac_salt[i] = salt[i] ^ 58
	}
	mac_key := pbkdf2.Key(byteKey, mac_salt, 2, KEY_SIZE, sha1.New)
	hash_mac := hmac.New(sha1.New, mac_key)
	hash_mac.Write(first[:len(first)-32])
	hash_mac.Write([]byte{1, 0, 0, 0})
	if bytes.Equal(hash_mac.Sum(nil), first[len(first)-32:len(first)-12]) {
		fmt.Println("Decryption Success")
	} else {
		fmt.Println("Password Error")
	}

	// 将python代码：blist = [blist[i:i + DEFAULT_PAGESIZE] for i in range(DEFAULT_PAGESIZE, len(blist), DEFAULT_PAGESIZE)] 转成go语言
	newblist := make([][]byte, 0)
	for i := DEFAULT_PAGESIZE; i < len(blist); i += DEFAULT_PAGESIZE {
		newblist = append(newblist, blist[i:i+DEFAULT_PAGESIZE])
	}

	// 将文件写入decryptePath
	deFile, err := os.OpenFile(decryptedPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer deFile.Close()
	deFile.Write([]byte(SQLITE_FILE_HEADER))
	t, err := aes.NewCipher(byteKey)
	if err != nil {
		return err
	}
	iv := first[len(first)-48 : len(first)-32]
	blockMode := cipher.NewCBCDecrypter(t, iv)
	decrypted := make([]byte, len(first)-48)
	blockMode.CryptBlocks(decrypted, first[:len(first)-48])
	deFile.Write(decrypted)
	deFile.Write(first[len(first)-48:])

	for _, i := range newblist {
		t, err := aes.NewCipher(byteKey)
		if err != nil {
			return err
		}
		blockMode := cipher.NewCBCDecrypter(t, i[len(i)-48:len(i)-32])
		decrypted := make([]byte, len(i)-48)
		blockMode.CryptBlocks(decrypted, i[:len(i)-48])
		deFile.Write(decrypted)
		deFile.Write(i[len(i)-48:])
	}
	return nil
}
