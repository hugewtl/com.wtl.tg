package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"unsafe"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main2() {
	//连接数据库
	runAsTerminal("cd /home/antif; ls -l")
}

func runAsTerminal(cmd string) {
	session, _ := ConnectSSH()
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	//解决无效句柄的问题
	fd = int(os.Stdout.Fd())

	termWidth, termHeight, err1 := terminal.GetSize(fd)
	if err1 != nil {
		panic(err1)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err2 := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err2 != nil {
		log.Fatal(err2)
	}

	session.Run(cmd)
}

//连接远程机器
func ConnectSSH() (*ssh.Session, error) {

	var input, input1, password string

	fmt.Print("请输入远程服务器的ip: ")
	fmt.Scanf("%s\n", &input)
	host := input

	fmt.Print("请输入远程服务器的用户名: ")
	fmt.Scanf("%s\n", &input1)
	user := input1

	fmt.Print("请输入远程服务器的密码: ")
	passwd, err2 := gopass.GetPasswd()

	port := 22
	if err2 != nil {
		fmt.Printf("%s", err2)
	} else {
		password = bytes2str(passwd)
	}
	var (
		// auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	// auth = make([]ssh.AuthMethod, 0)
	// auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
