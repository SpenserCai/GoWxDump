/*
 * @Author: SpenserCai
 * @Date: 2023-02-20 18:15:51
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 15:21:47
 * @Description: file content
 */
package main

import (
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
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

// TelBot对象
var TelBot *tele.Bot
