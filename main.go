/*
 * @Author: SpenserCai
 * @Date: 2023-02-17 14:14:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-23 15:39:19
 * @Description: file content
 */
package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func main() {
	// 从参数读取token
	if len(os.Args) > 1 {
		TELBOT_TOKEN = os.Args[1]
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
