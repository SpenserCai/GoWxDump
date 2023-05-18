package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"golang.org/x/sys/windows"
)

// 获取微信进程对象，包含进程ID、进程句柄和Module列表
func GetWeChatProcess() (windows.ProcessEntry32, error) {
	var process windows.ProcessEntry32
	process.Size = uint32(unsafe.Sizeof(process))
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return process, err
	}
	defer windows.CloseHandle(snapshot)
	for {
		err = windows.Process32Next(snapshot, &process)
		if err != nil {
			return process, err
		}
		if windows.UTF16ToString(process.ExeFile[:]) == "WeChat.exe" {
			return process, nil
		}
	}
}

// 获取微信进程的WeChatWin.dll模块对象，包含模块基址、模块大小和模块路径()
func GetWeChatWinModule(process windows.ProcessEntry32) (windows.ModuleEntry32, error) {
	var module windows.ModuleEntry32
	module.Size = uint32(unsafe.Sizeof(module))
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, process.ProcessID)
	if err != nil {
		return module, err
	}
	defer windows.CloseHandle(snapshot)
	for {
		err = windows.Module32Next(snapshot, &module)
		if err != nil {
			return module, err
		}
		if windows.UTF16ToString(module.Module[:]) == "WeChatWin.dll" {
			return module, nil
		}
	}
}

// 通过模块获取版本号 c#代码为：string FileVersion = processModule.FileVersionInfo.FileVersion;转成go代码如下
func GetVersion(module windows.ModuleEntry32) (string, error) {
	image, imgErr := windows.LoadLibraryEx(windows.UTF16ToString(module.ExePath[:]), 0, windows.LOAD_LIBRARY_AS_DATAFILE)
	if imgErr != nil {
		return "", fmt.Errorf("LoadLibraryEx error: %v", imgErr)
	}
	resInfo, infoErr := windows.FindResource(image, windows.ResourceID(1), windows.RT_VERSION)
	if infoErr != nil {
		return "", fmt.Errorf("FindResource error: %v", infoErr)
	}
	resData, dataErr := windows.LoadResourceData(image, resInfo)
	if dataErr != nil {
		return "", fmt.Errorf("LoadResourceData error: %v", dataErr)
	}
	var info *windows.VS_FIXEDFILEINFO
	size := uint32(unsafe.Sizeof(*info))
	err := windows.VerQueryValue(unsafe.Pointer(&resData[0]), `\`, unsafe.Pointer(&info), &size)
	if err != nil {
		return "", fmt.Errorf("VerQueryValue error: %v", err)
	}
	// 从低位到高位，分别为主版本号、次版本号、修订号、编译号
	version := fmt.Sprintf("%d.%d.%d.%d", info.FileVersionMS>>16, info.FileVersionMS&0xffff, info.FileVersionLS>>16, info.FileVersionLS&0xffff)
	return version, nil
}

// 获取微信数据：入参为微信进程句柄，偏移地址，返回值为昵称和错误信息
func GetWeChatData(process windows.Handle, offset uintptr, nSize int) (string, error) {
	var buffer = make([]byte, nSize)
	err := windows.ReadProcessMemory(process, offset, &buffer[0], uintptr(nSize), nil)
	if err != nil {
		return "", err
	}
	// 声明一个字节数组，暂时为空
	var textBytes []byte = nil
	for _, v := range buffer {
		if v == 0 {
			break
		}
		textBytes = append(textBytes, v)
	}
	// 返回utf8编码的字符串
	return string(textBytes), nil
}

// 获取微信key：入参为微信进程句柄，偏移地址，返回值为key和错误信息
func GetWeChatKey(process windows.Handle, offset uintptr) (string, error) {
	var buffer = make([]byte, 4)
	err := windows.ReadProcessMemory(process, offset, &buffer[0], 4, nil)
	if err != nil {
		return "", err
	}
	var num = 32
	var buffer2 = make([]byte, num)
	// c# 代码(IntPtr)(((int)array[3] << 24) + ((int)array[2] << 16) + ((int)array[1] << 8) + (int)array[0]);转成go代码如下
	offset2 := uintptr((int(buffer[3]) << 24) + (int(buffer[2]) << 16) + (int(buffer[1]) << 8) + int(buffer[0]))
	err = windows.ReadProcessMemory(process, offset2, &buffer2[0], uintptr(num), nil)
	if err != nil {
		return "", err
	}
	// 将byte数组转成hex字符串，并转成大写
	key := hex.EncodeToString(buffer2)
	key = strings.ToUpper(key)
	return key, nil

}

func GetWeChatFromRegistry() (string, error) {
	// 打开注册表的微信路径：HKEY_CURRENT_USER\Software\Tencent\WeChat\FileSavePath
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Tencent\WeChat`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()
	//获取key的值
	value, _, err := key.GetStringValue("FileSavePath")
	if err != nil {
		return "", err
	}
	return value, nil

}

func GetWeChatDir() (string, error) {
	msgDir := ""
	wDir, err := GetWeChatFromRegistry()
	// 如果wDir为MyDocument:
	if wDir == "MyDocument:" {
		// 获取%USERPROFILE%/Documents目录
		profile, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		// 获取微信消息目录
		msgDir = filepath.Join(profile, "Documents", "WeChat Files")
	} else {
		// 获取微信消息目录
		msgDir = filepath.Join(wDir, "WeChat Files")
	}
	// 判断目录是否存在
	_, err = os.Stat(msgDir)
	if err != nil {
		// 手动输入目录
		fmt.Print("输入WeChat Files所在目录路径（如d:\\documents）：")
		fmt.Scanln(&msgDir)
		msgDir = filepath.Join(msgDir, "WeChat Files")
		// 判断目录是否存在
		_, err = os.Stat(msgDir)
		if err != nil {
			return "", err
		}
	}
	return msgDir, nil
}

func IsSupportAutoGetData(version string) bool {
	// 判断version是否在支持的版本列表中
	for _, v := range SupportAutoGetDataVersionList {
		if version == v {
			return true
		}
	}
	return false
}

// 获取微信消息目录下的所有用户目录，排除All Users目录和Applet目录，返回一个map，key用户id，value用户目录
func GetWeChatUserDir(wechatRoot string) (map[string]string, error) {
	userDir := make(map[string]string)
	// 获取微信消息目录下的所有用户目录
	files, err := ioutil.ReadDir(wechatRoot)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// 排除All Users目录和Applet目录
		if file.Name() == "All Users" || file.Name() == "Applet" {
			continue
		}
		userDir[file.Name()] = filepath.Join(wechatRoot, file.Name())
	}
	return userDir, nil
}

// 解密微信数据库
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

// 复制微信的数据文件
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
