package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"net"
)

var (
	hostPrivateKeySigner ssh.Signer
)

func main() {
	app := cli.NewApp()
	app.Name = "sshd-go-sample"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{}
	app.Action = doMain
	app.Run(os.Args)
}

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Println(conn.RemoteAddr(), "authenticate with", key.Type())
	return nil, nil
}

func doMain(c *cli.Context) {
	log.Infof("doMain")

	keyPath := "./host_key"

	hostPrivateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	log.Infof("hostPrivateKey=%s", hostPrivateKey)

	hostPrivateKeySigner, err = ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		panic(err)
	}

	log.Infof("hostPrivateKeySigner=%s", hostPrivateKeySigner)

	config := ssh.ServerConfig{
		PublicKeyCallback: keyAuth,
	}

	config.AddHostKey(hostPrivateKeySigner)

	port := "2222"

	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		// From a standard TCP connection to an encrypted SSH connection
		sshConn, _, _, err := ssh.NewServerConn(conn, &config)
		if err != nil {
			panic(err)
		}

		log.Println("Connection from", sshConn.RemoteAddr())
		sshConn.Close()
	}
}
