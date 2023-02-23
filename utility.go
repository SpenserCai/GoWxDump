/*
 * @Author: SpenserCai
 * @Date: 2023-02-21 10:44:10
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 16:12:29
 * @Description: file content
 */
package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func CopyFile(src, dst string) error {
	// 判断源文件是否存在
	_, err := os.Stat(src)
	if err != nil {
		return err
	}
	// 读取源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// 拷贝文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// 获取公网ip
func GetPublicIp() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
