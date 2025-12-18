package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 简易 SFTP 服务端，使用内存生成的 Ed25519 host key 与用户名/密码认证
// 用法示例：
//
//	go run ./test/sftpserver -addr :2022 -user test -pass 123456
func main() {
	addr := flag.String("addr", ":2022", "监听地址，如 :2022")
	user := flag.String("user", "test", "登录用户名")
	pass := flag.String("pass", "123456", "登录密码")
	flag.Parse()

	// 生成临时 Ed25519 主机密钥
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("generate host key failed: %v", err)
	}
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		log.Fatalf("new signer failed: %v", err)
	}

	// SSH Server 配置
	conf := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if conn.User() == *user && string(password) == *pass {
				return nil, nil
			}
			return nil, fmt.Errorf("username/password rejected for %s", conn.User())
		},
	}
	conf.AddHostKey(signer)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s failed: %v", *addr, err)
	}
	log.Printf("SFTP server listening on %s (user=%s)\n", *addr, *user)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			sshConn, chans, reqs, err := ssh.NewServerConn(c, conf)
			if err != nil {
				log.Printf("handshake failed: %v", err)
				return
			}
			defer sshConn.Close()
			go ssh.DiscardRequests(reqs)

			for newChan := range chans {
				if newChan.ChannelType() != "session" {
					newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}
				ch, reqs, err := newChan.Accept()
				if err != nil {
					log.Printf("channel accept failed: %v", err)
					continue
				}

				go func(ch ssh.Channel, in <-chan *ssh.Request) {
					defer ch.Close()
					for req := range in {
						// 仅处理 subsystem=sftp
						if req.Type == "subsystem" {
							req.Reply(true, nil)
							server, err := sftp.NewServer(ch)
							if err != nil {
								log.Printf("sftp new server failed: %v", err)
								return
							}
							if err := server.Serve(); err == io.EOF {
								_ = server.Close()
								return
							} else if err != nil {
								log.Printf("sftp serve failed: %v", err)
							}
							return
						}
						req.Reply(false, nil)
					}
				}(ch, reqs)
			}
		}(conn)
	}
}
