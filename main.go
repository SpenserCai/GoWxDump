/*
 * @Author: SpenserCai
 * @Date: 2023-02-17 14:14:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-20 11:57:23
 * @Description: file content
 */
package main

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func GetWeChatInfo(wechatProcessHandle windows.Handle, module windows.ModuleEntry32) (WeChatData, error) {
	// 初始化WeChatData
	wechatData := WeChatData{}
	// 获取微信版本
	version, err := GetVersion(module)
	if err != nil {
		return wechatData, err
	}
	// 获取微信昵称
	nickName, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][0]), 100)
	if err != nil {
		return wechatData, err
	}
	// 获取微信账号
	account, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][1]), 100)
	if err != nil {
		return wechatData, err
	}
	// 获取微信手机号
	mobile, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][2]), 100)
	if err != nil {
		return wechatData, err
	}
	// 获取微信密钥
	key, err := GetWeChatKey(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[version][4]))
	if err != nil {
		return wechatData, err
	}
	// 设置微信数据
	wechatData.Version = version
	wechatData.NickName = nickName
	wechatData.Account = account
	wechatData.Mobile = mobile
	wechatData.Key = key
	return wechatData, nil
}

func main() {
	process, err := GetWeChatProcess()
	if err != nil {
		fmt.Println("GetWeChatProcess error: ", err)
		return
	}
	wechatProcessHandle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, process.ProcessID)
	if err != nil {
		fmt.Println("OpenProcess error: ", err)
		return
	}
	module, err := GetWeChatWinModule(process)
	if err != nil {
		fmt.Println("GetWeChatWinModule error: ", err)
		return
	}
	wechatData, err := GetWeChatInfo(wechatProcessHandle, module)
	if err != nil {
		fmt.Println("GetWeChatInfo error: ", err)
		return
	}
	fmt.Printf("WeChat Version: %s \n", wechatData.Version)
	fmt.Printf("WeChat NickName: %s \n", wechatData.NickName)
	fmt.Printf("WeChat Account: %s \n", wechatData.Account)
	fmt.Printf("WeChat Mobile: %s \n", wechatData.Mobile)
	fmt.Printf("WeChat Key: %s \n", wechatData.Key)
	fmt.Println("---------------------------------------------------------------------------------------------")
	// 获取用户数据目录
	wechatRoot, err := GetWeChatDir()
	if err != nil {
		fmt.Println("请手动设置微信消息目录")
		return
	}
	// 获取用户目录
	userDir, err := GetWeChatUserDir(wechatRoot)
	if err != nil {
		fmt.Println("GetWeChatUserDir error: ", err)
		return
	}
	for k, v := range userDir {
		fmt.Printf("[%s]:%s \n", k, v)
	}
	// 判断是否支持自动获取数据目录（version是否在SupportAutoGetDataVersionList列表中）
	if !IsSupportAutoGetData(wechatData.Version) {
		fmt.Println("不支持自动获取数据目录")
		return
	}
	// 获取用户数据目录
	dataDirName, err := GetWeChatData(wechatProcessHandle, module.ModBaseAddr+uintptr(OffSetMap[wechatData.Version][5]), 100)
	if err != nil {
		fmt.Println("GetWeChatDataDir error: ", err)
		return
	}
	// 获取用户数据目录，拼接成绝对路径
	dataDir := filepath.Join(wechatRoot, dataDirName)
	fmt.Println("WeChat DataDir: ", dataDir)

}
