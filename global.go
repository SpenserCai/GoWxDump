/*
 * @Author: SpenserCai
 * @Date: 2023-02-20 18:15:51
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-24 16:49:16
 * @Description: file content
 */
package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
	tele "gopkg.in/telebot.v3"
)

// 定义微信数据结构
type WeChatData struct {
	Version        string
	NickName       string
	Account        string
	Mobile         string
	Key            string
	WeChatProcess  windows.ProcessEntry32
	WeChatHandle   windows.Handle
	WeChatWinBaseAddr uint64
	WeChatWinFullName string
}

var PROCESS_ALL_ACCESS = uint32(
	windows.PROCESS_QUERY_INFORMATION |
		windows.PROCESS_VM_READ |
		windows.PROCESS_VM_WRITE |
		windows.PROCESS_VM_OPERATION |
		windows.PROCESS_CREATE_THREAD |
		windows.PROCESS_DUP_HANDLE |
		windows.PROCESS_TERMINATE |
		windows.PROCESS_SUSPEND_RESUME |
		windows.PROCESS_SET_QUOTA |
		windows.PROCESS_SET_INFORMATION |
		windows.PROCESS_QUERY_LIMITED_INFORMATION)

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
