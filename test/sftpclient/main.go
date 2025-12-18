package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 简易 SFTP 客户端：支持上传、下载、列目录
// 用法示例：
//
//	go run ./test/sftpclient -addr 127.0.0.1:2022 -user test -pass 123456 -ls /
//	go run ./test/sftpclient -addr 127.0.0.1:2022 -user test -pass 123456 -get /remote/file -o local.file
//	go run ./test/sftpclient -addr 127.0.0.1:2022 -user test -pass 123456 -put local.file -to /remote/file
func main() {
	addr := flag.String("addr", "127.0.0.1:2022", "服务端地址")
	user := flag.String("user", "test", "用户名")
	pass := flag.String("pass", "123456", "密码")
	ls := flag.String("ls", "", "列出远程目录")
	get := flag.String("get", "", "下载远程文件路径")
	output := flag.String("o", "", "下载到的本地文件路径")
	put := flag.String("put", "", "上传本地文件路径")
	to := flag.String("to", "", "上传远程文件路径")
	flag.Parse()

	sshConf := &ssh.ClientConfig{
		User:            *user,
		Auth:            []ssh.AuthMethod{ssh.Password(*pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", *addr, sshConf)
	if err != nil {
		log.Fatalf("ssh dial failed: %v", err)
	}
	defer client.Close()

	s, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalf("new sftp client failed: %v", err)
	}
	defer s.Close()

	if *ls != "" {
		entries, err := s.ReadDir(*ls)
		if err != nil {
			log.Fatalf("readdir failed: %v", err)
		}
		for _, e := range entries {
			fmt.Printf("%s\t%12d\t%v\n", e.Name(), e.Size(), e.ModTime())
		}
		return
	}

	if *get != "" {
		remote, err := s.Open(*get)
		if err != nil {
			log.Fatalf("open remote failed: %v", err)
		}
		defer remote.Close()
		localPath := *output
		if localPath == "" {
			localPath = filepath.Base(*get)
		}
		local, err := os.Create(localPath)
		if err != nil {
			log.Fatalf("create local failed: %v", err)
		}
		defer local.Close()
		if _, err := io.Copy(local, remote); err != nil {
			log.Fatalf("download failed: %v", err)
		}
		fmt.Printf("downloaded %s -> %s\n", *get, localPath)
		return
	}

	if *put != "" {
		if *to == "" {
			log.Fatalf("-to 远程路径必填")
		}
		local, err := os.Open(*put)
		if err != nil {
			log.Fatalf("open local failed: %v", err)
		}
		defer local.Close()
		remote, err := s.Create(*to)
		if err != nil {
			log.Fatalf("create remote failed: %v", err)
		}
		defer remote.Close()
		if _, err := io.Copy(remote, local); err != nil {
			log.Fatalf("upload failed: %v", err)
		}
		fmt.Printf("uploaded %s -> %s\n", *put, *to)
		return
	}

	flag.Usage()
}
