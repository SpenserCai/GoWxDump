/*
 * @Author: SpenserCai
 * @Date: 2023-02-20 18:15:51
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-21 10:45:14
 * @Description: file content
 */
package main

import (
	"os"
	"path/filepath"
)

// 初始化全局的微信信息对象
var WeChatDataObject = WeChatData{}

// 获取程序内运行的目录
var CurrentPath = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}()
