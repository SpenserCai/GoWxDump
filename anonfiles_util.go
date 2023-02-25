/*
 * @Author: SpenserCai
 * @Date: 2023-02-24 16:36:01
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-25 13:07:29
 * @Description: file content
 */
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type AnonFileError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

type AnonFileUploadData struct {
	File struct {
		Url struct {
			Full  string `json:"full"`
			Short string `json:"short"`
		} `json:"url"`
		Metadata struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Size struct {
				Bytes    int    `json:"bytes"`
				Readable string `json:"readable"`
			} `json:"size"`
		} `json:"metadata"`
	} `json:"file"`
}

// AnonFiles的基础返回结构 Status是bool类型的，Data是对象类型的
type AnonFilesBaseResponse struct {
	Status bool                   `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Error  AnonFileError          `json:"error"`
}

func AnonFilesUpload(filePath string) (string, error) {
	apiUrl := fmt.Sprintf("https://api.anonfiles.com/upload?token=%s", ANONFILES_TOKEN)
	// curl -F "file=@filePath" url 转成go代码并且使用代理
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// 创建一个缓冲区
	buf := new(bytes.Buffer)
	// 创建一个multipart writer
	writer := multipart.NewWriter(buf)
	// 创建一个form-data
	formFile, err := writer.CreateFormFile("file", strings.Split(filePath, "\\")[len(strings.Split(filePath, "\\"))-1])
	if err != nil {
		return "", err
	}
	// 将文件写入到form-data
	_, err = io.Copy(formFile, file)
	if err != nil {
		return "", err
	}
	// 关闭writer
	err = writer.Close()
	if err != nil {
		return "", err
	}
	// 创建一个http请求,并使用代理socks5://127.0.0.1:LOCAL_PROXY_PORT
	req, err := http.NewRequest("POST", apiUrl, buf)
	if err != nil {
		return "", err
	}
	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 创建一个http client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "socks5",
				Host:   "127.0.0.1:" + strconv.Itoa(LOCAL_PROXY_PORT),
			}),
		},
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 解析响应
	var response AnonFilesBaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	// 判断是否成功
	if response.Status {
		// 将data转成AnonFileUploadData
		var data AnonFileUploadData
		err = mapstructure.Decode(response.Data, &data)
		if err != nil {
			return "", err
		}
		// 返回短链接
		return data.File.Url.Short, nil
	} else {
		// 失败
		return "", errors.New(response.Error.Message)
	}

}
