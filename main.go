package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	sshKeys := []string{
		"/home/manishbishwakarma/.ssh/id_rsa",
		"/home/manishbishwakarma/.ssh/id_ed25519",
	}

	serverAdd := "localhost:22" 

	var sshClient *ssh.Client
	var dialErr error

	hostKeyCallback, err := knownhosts.New("/home/manishbishwakarma/.ssh/known_hosts")
	if err != nil {
		log.Fatalf("Could not load known_hosts file: %v", err)
	}

	for _, key := range sshKeys {
		keyBytes, err := os.ReadFile(key)
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", key, err)
			continue
		}

		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			fmt.Printf("Failed to parse key %s: %v\n", key, err)
			continue
		}

		sshConf := &ssh.ClientConfig{
			User: "manishbishwakarma", 
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: hostKeyCallback,
			Timeout:         4 * time.Second,
		}

		sshClient, dialErr = ssh.Dial("tcp", serverAdd, sshConf)
		if dialErr != nil {
			fmt.Printf("Key login failed using %s\n", key)
			continue
		}
		fmt.Printf("Successfully logged in using %s!\n", key)
		break
	}

	if sshClient == nil {
		log.Fatalf("All keys didn't work! Last error: %v", dialErr)
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	fmt.Println("Running command: cat /etc/os-release")
	output, sessErr := session.CombinedOutput("cat /etc/os-release")
	if sessErr != nil {
		log.Fatalf("Failed to run command: %v", sessErr)
	}

	fmt.Println("--- SERVER OUTPUT ---")
	fmt.Println(string(output))
	fmt.Println("---------------------")

	select {}
}