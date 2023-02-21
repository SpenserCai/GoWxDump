package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetWeChatInfo() error {
	// 获取微信版本
	version, err := GetVersion(WeChatDataObject.WeChatWinModel)
	if err != nil {
		return err
	}
	// 获取微信昵称
	nickName, err := GetWeChatData(WeChatDataObject.WeChatHandle, WeChatDataObject.WeChatWinModel.ModBaseAddr+uintptr(OffSetMap[version][0]), 100)
	if err != nil {
		return err
	}
	// 获取微信账号
	account, err := GetWeChatData(WeChatDataObject.WeChatHandle, WeChatDataObject.WeChatWinModel.ModBaseAddr+uintptr(OffSetMap[version][1]), 100)
	if err != nil {
		return err
	}
	// 获取微信手机号
	mobile, err := GetWeChatData(WeChatDataObject.WeChatHandle, WeChatDataObject.WeChatWinModel.ModBaseAddr+uintptr(OffSetMap[version][2]), 100)
	if err != nil {
		return err
	}
	// 获取微信密钥
	key, err := GetWeChatKey(WeChatDataObject.WeChatHandle, WeChatDataObject.WeChatWinModel.ModBaseAddr+uintptr(OffSetMap[version][4]))
	if err != nil {
		return err
	}
	// 设置微信数据
	WeChatDataObject.Version = version
	WeChatDataObject.NickName = nickName
	WeChatDataObject.Account = account
	WeChatDataObject.Mobile = mobile
	WeChatDataObject.Key = key
	return nil
}

func ShowInfoCmd() {
	fmt.Printf("WeChat Version: %s \n", WeChatDataObject.Version)
	fmt.Printf("WeChat NickName: %s \n", WeChatDataObject.NickName)
	fmt.Printf("WeChat Account: %s \n", WeChatDataObject.Account)
	fmt.Printf("WeChat Mobile: %s \n", WeChatDataObject.Mobile)
	fmt.Printf("WeChat Key: %s \n", WeChatDataObject.Key)
}

func DecryptCmd() {
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
	dataDir := ""
	if IsSupportAutoGetData(WeChatDataObject.Version) {
		// 获取用户数据目录
		dataDirName, err := GetWeChatData(WeChatDataObject.WeChatHandle, WeChatDataObject.WeChatWinModel.ModBaseAddr+uintptr(OffSetMap[WeChatDataObject.Version][5]), 100)
		if err != nil {
			fmt.Println("GetWeChatDataDir error: ", err)
			return
		}
		// 获取用户数据目录，拼接成绝对路径
		dataDir = filepath.Join(wechatRoot, dataDirName)
	}

	// 判断目录是否存在如果不存，要求用户从userDir中选择一个目录
	_, err = os.Stat(dataDir)
	if err != nil {
		fmt.Println("物资自动识别，请从下面选择一个目录")
		for k, v := range userDir {
			fmt.Printf("[%s]:%s \n", k, v)
		}
		var input string
		// 提示输入
		fmt.Print("请选择上述id中的一个:")
		fmt.Scanln(&input)
		// 判断输入是否合法
		if _, ok := userDir[input]; !ok {
			fmt.Println("输入错误")
			return
		}
		dataDir = userDir[input]
	}
	fmt.Println("WeChat DataDir: ", dataDir)
	// 复制聊天记录文件到缓存目录dataDir + \Msg\Multi
	err = CopyMsgDb(filepath.Join(dataDir, "Msg", "Multi"))
	if err != nil {
		fmt.Println("CopyMsgDb error: ", err)
		return
	}
	err = CopyMsgDb(filepath.Join(dataDir, "Msg"))
	if err != nil {
		fmt.Println("CopyMicroMsgDb error: ", err)
		return
	}
	// 解密tmp目录下的所有.db文件，解密后的文件放在decrypted目录下
	err = DecryptDb(WeChatDataObject.Key)
	if err != nil {
		fmt.Println("DecryptDb error: ", err)
		return
	}
	// 清理缓存目录
	err = os.RemoveAll(CurrentPath + "\\tmp")
	if err != nil {
		fmt.Println("RemoveAll error: ", err)
		return
	}
}
