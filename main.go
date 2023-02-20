/*
 * @Author: SpenserCai
 * @Date: 2023-02-17 14:14:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-02-20 16:48:28
 * @Description: file content
 */
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// 获取程序内运行的目录
var CurrentPath = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}()

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

func CopyFile(src, dst string) error {
	// 判断源文件是否存在
	_, err := os.Stat(src)
	if err != nil {
		return err
	}
	// 读取源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// 拷贝文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func CopyMsgDb(dataDir string) error {
	// 判断目录是否存在
	_, err := os.Stat(dataDir)
	if err != nil {
		return err
	}
	// 判断运行目录是否存在tmp目录没有则创建
	_, err = os.Stat(CurrentPath + "\\tmp")
	if err != nil {
		err = os.Mkdir(CurrentPath+"\\tmp", os.ModePerm)
		if err != nil {
			return err
		}
	}
	// 正则匹配，将所有MSG数字.db文件拷贝到tmp目录，不扫描子目录
	err = filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ok, _ := filepath.Match("MSG*.db", info.Name()); ok {
			err = CopyFile(path, CurrentPath+"\\tmp\\"+info.Name())
			if err != nil {
				return err
			}
		}
		// 复制MicroMsg.db到tmp目录
		if ok, _ := filepath.Match("MicroMsg.db", info.Name()); ok {
			err = CopyFile(path, CurrentPath+"\\tmp\\"+info.Name())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 如果不存在decrypted目录则创建
	_, err = os.Stat(CurrentPath + "\\decrypted")
	if err != nil {
		err = os.Mkdir(CurrentPath+"\\decrypted", os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func DecryptDb(key string) error {
	// 判断tmp目录是否存在
	_, err := os.Stat(CurrentPath + "\\tmp")
	if err != nil {
		return err
	}
	// 判断decrypted目录是否存在
	_, err = os.Stat(CurrentPath + "\\decrypted")
	if err != nil {
		return err
	}
	// 正则匹配，将所有MSG数字.db文件解密到decrypted目录，不扫描子目录
	err = filepath.Walk(CurrentPath+"\\tmp", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ok, _ := filepath.Match("*.db", info.Name()); ok {
			err = Decrypt(key, path, CurrentPath+"\\decrypted\\"+info.Name())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
	// for k, v := range userDir {
	//     fmt.Printf("[%s]:%s \n", k, v)
	// }
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
	// 判断目录是否存在如果不存，要求用户从userDir中选择一个目录
	_, err = os.Stat(dataDir)
	if err != nil {
		fmt.Println("数据目录不存在，请从下面选择一个目录")
		for k, v := range userDir {
			fmt.Printf("[%s]:%s \n", k, v)
		}
		var input string
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
	err = DecryptDb(wechatData.Key)
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
