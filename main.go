/*
 * @Author: SpenserCai
 * @Date: 2023-02-17 14:14:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 16:56:17
 * @Description: file content
 */
package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

func main() {
	// 从参数读取token
	if len(os.Args) > 2 {
		TELBOT_TOKEN = os.Args[1]
		// 将参数转换为int
		chatId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("chatId error: ", err)
			return
		}
		TELBOT_CHAT_ID = chatId
		InitBot()
	}

	// 获取微信进程
	process, err := GetWeChatProcess()
	if err != nil {
		fmt.Println("GetWeChatProcess error: ", err)
		return
	}
	WeChatDataObject.WeChatProcess = process

	// 获取微信句柄
	wechatProcessHandle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, process.ProcessID)
	if err != nil {
		fmt.Println("OpenProcess error: ", err)
		return
	}
	WeChatDataObject.WeChatHandle = wechatProcessHandle

	// 获取微信模块
	module, err := GetWeChatWinModule(process)
	if err != nil {
		fmt.Println("GetWeChatWinModule error: ", err)
		return
	}
	WeChatDataObject.WeChatWinModel = module

	err = GetWeChatInfo()
	if err != nil {
		fmt.Println("GetWeChatInfo error: ", err)
		return
	}
	RunCommand()

}
