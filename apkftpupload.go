package main

import (
	"fmt"
	ftp  "github.com/jlaffaye/ftp"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var serverlist = [2]string{"xx.xx.xx.xx:21", "xx.xx.xx.xx:21"}
var quit = make(chan int, len(serverlist))

func main() {
	fmt.Println("upload begin!")
	var dir string
	if len(os.Args) != 2 {
		//log.Fatal("Usage:" + filepath.Base(os.Args[0]) + " log_dir ")
		//os.Exit(1)
		dir = "D:\\apkupload"
	}else
	{
		dir = os.Args[1]
	}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if m, _ := regexp.MatchString("^\\S+\\.apk$",f.Name()); m{
			fmt.Println(f.Name())
			fmt.Println(path)
			domainName := "/";
			fmt.Println(time.Now())
			for i := 0; i < len(serverlist) ;i++ {
				fmt.Printf("%s ", serverlist[i])
				go ftpUploadFile(serverlist[i], "apkUpload", "xxxx", path, domainName,f.Name())
			}
			for i := 0; i < len(serverlist) ;i++   {
				<- quit
			}
			//ftpUploadFile("xx.xx.xx.xx:21", "apkUpload", "xxxx", path, domainName,f.Name())
		}
		return nil
	})
	fmt.Println("upload end!")
}

func ftpUploadFile(ftpserver, ftpuser, pw, localFile, remoteSavePath, saveName string) {
	ftp, err := ftp.Connect(ftpserver)
	if err != nil {
		fmt.Println(err)
	}
	err = ftp.Login(ftpuser, pw)
	if err != nil {
		fmt.Println(err)
	}
	//注意是 pub/log，不能带“/”开头
	//ftp.ChangeDir("pub/log")
	dir, err := ftp.CurrentDir()
	fmt.Println(dir)
	ftp.MakeDir(remoteSavePath)
	ftp.ChangeDir(remoteSavePath)
	dir, _ = ftp.CurrentDir()
	fmt.Println(dir)
	file, err := os.Open(localFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	err = ftp.Stor(saveName, file)
	if err != nil {
		fmt.Println(err)
	}
	ftp.Logout()
	ftp.Quit()
	fmt.Println("success upload file:", localFile)
	quit <- 0
}

