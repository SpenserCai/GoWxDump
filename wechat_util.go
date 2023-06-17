package main

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lwow64ext
#include "wow64ext.h"
*/
import "C"

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
	"errors"
)

var WeChatExe = "WeChat.exe"
var WeChatWin  = "WeChatWin.dll"



// 获取微信进程对象，包含进程ID、进程句柄和Module列表
func GetWeChatProcess() (windows.ProcessEntry32, error) {
	var process windows.ProcessEntry32
	process.Size = uint32(unsafe.Sizeof(process))
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return process, err
	}
	defer windows.CloseHandle(snapshot)
	
	err = windows.Process32First(snapshot, &process)
	if err != nil {
	    fmt.Println(err);
		return process, err
	}

	if windows.UTF16ToString(process.ExeFile[:]) == WeChatExe {
		return process, nil
	}	  
	
	for {
		err = windows.Process32Next(snapshot, &process)
		if err != nil {
			return process, err
		}
//		fmt.Printf("pid:%d, ppid:%d, %d,  %d, exe:%s\n", process.ProcessID,  process.ParentProcessID,  process.DefaultHeapID   , unsafe.Sizeof(uintptr(0)), windows.UTF16ToString(process.ExeFile[:]) );
		if windows.UTF16ToString(process.ExeFile[:]) == WeChatExe {
			return process, nil
		}
	}
}

// 获取微信进程的WeChatWin.dll模块对象，包含模块基址、模块大小和模块路径()
func GetWeChatWinModule(process windows.ProcessEntry32) (uint64, string, error) {
	var module windows.ModuleEntry32
	module.Size = uint32(unsafe.Sizeof(module))
	
	modname,_ := windows.UTF16FromString(WeChatWin);
	var _fullname  [1024]uint16;
	_fullname[0] = 0; _fullname[1] = 0;

	addr := C.GetProcessModuleHandle64(C.int(WeChatDataObject.WeChatHandle), (*C.wchar_t)(&modname[0]), (*C.wchar_t)(&_fullname[0])); 
	fullname := windows.UTF16ToString(_fullname[:]);
	return uint64(addr), fullname, nil;

// bool ReadProcessMemory64(int hProcess, DWORD64 lpBaseAddress, void* lpBuffer, size_t nSize, size_t *lpNumberOfBytesRead);



	/*
	var hMods [1024]windows.Handle
	var cbNeeded uint32

	err := windows.EnumProcessModulesEx(WeChatDataObject.WeChatHandle, &hMods[0], uint32(unsafe.Sizeof(hMods)), &cbNeeded, windows.LIST_MODULES_ALL           ); 
	if  err != nil {
		fmt.Println("Failed to enumerate modules: ", err)
		return module, err
	}

	modCount := cbNeeded / uint32(unsafe.Sizeof(hMods[0]))
	for i := uint32(0); i < modCount; i++ {
		var modName [windows.MAX_PATH]uint16
		err = windows.GetModuleFileNameEx(WeChatDataObject.WeChatHandle, hMods[i], &modName[0], uint32(len(modName)))
		if err != nil {
			fmt.Println("Failed to get module file name: ", err)
			continue
		}
		if windows.UTF16ToString(modName[:]) == WeChatWin {
			fmt.Printf("Found module abc.dll at address: 0x%X\n", hMods[i])
		} else {
			fmt.Printf("module:%s, addr:0x%x\n",   windows.UTF16ToString(modName[:]), hMods[i]);
		}
	}
	return module, fmt.Errorf("geterror")
	*/
	
	/*	
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, process.ProcessID)
	if err != nil {
	    fmt.Println(err)
		return module, err
	}
	defer windows.CloseHandle(snapshot)
	
	err = windows.Module32First(snapshot, &module)
	if err != nil {
		return module, err
	}
	fmt.Println(windows.UTF16ToString(module.Module[:]))
	if windows.UTF16ToString(module.Module[:]) == WeChatWin {
		return module, nil
	}
	for {
		err = windows.Module32Next(snapshot, &module)
		if err != nil {
			return module, err
		}
		fmt.Println(windows.UTF16ToString(module.Module[:]))
		if windows.UTF16ToString(module.Module[:]) == WeChatWin {
			return module, nil
		}
	}
  */
}

// 通过模块获取版本号 c#代码为：string FileVersion = processModule.FileVersionInfo.FileVersion;转成go代码如下
func GetVersion(fullname string) (string, error) {
	image, imgErr := windows.LoadLibraryEx(fullname, 0, windows.LOAD_LIBRARY_AS_DATAFILE)
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
func GetWeChatData(process windows.Handle, offset uint64, nSize int) (string, error) {
	var buffer = make([]byte, nSize)
	err := C.ReadProcessMemory64(C.int(process), C.ulonglong(offset), (unsafe.Pointer)(&buffer[0]), C.uint(nSize), nil)
	if err == 0 {
		return "", errors.New("ReadProcessMemory64 failed")
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




func GetWeChatKey(process windows.Handle, offset uint64) (string, error) {
    var wow64 bool;
	e := windows.IsWow64Process(process, &wow64);
	if e != nil {
		return "", e
	}

	var buffer = make([]byte, 8)
	err := C.ReadProcessMemory64(C.int(process), C.ulonglong(offset), (unsafe.Pointer)(&buffer[0]), 8, nil)
	if err == 0 {
		return "", errors.New("ReadProcessMemory64 failed x")
	}
	var num = 32
	var buffer2 = make([]byte, num)
	offset2 := (uint64(buffer[3]) << 24) + (uint64(buffer[2]) << 16) + (uint64(buffer[1]) << 8) + (uint64(buffer[0]) << 0) 
	if !wow64 {
		offset2 += (uint64(buffer[7]) << 56) + (uint64(buffer[6]) << 48) + (uint64(buffer[5]) << 40) + (uint64(buffer[4]) << 32)
	}
	
	err = C.ReadProcessMemory64(C.int(process), C.ulonglong(offset2), (unsafe.Pointer)(&buffer2[0]), C.uint(num), nil)
	if err == 0 {
		return "", errors.New("ReadProcessMemory64 failed y")
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
	if err != nil {
		return "", err
	}
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
		return "", err
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
		if file.Name() == "All Users" || file.Name() == "Applet" || file.Name() == "WMPF" {
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
