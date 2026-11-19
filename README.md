# Go SSH Client

A lightweight, production-grade SSH client written from scratch in Go. This project demonstrates how to establish a secure, encrypted connection to a remote Linux machine, authenticate using private keys, and execute commands over an SSH session.

## 📖 About The Project

This project was built as a learning journey into network protocols, cryptography, and the Go programming language. Instead of relying on the system's built-in `ssh` command, this tool uses Go's official `golang.org/x/crypto/ssh` library to handle the underlying SSH transport layer, key exchanges, and authentication.

### Key Features
* **Multiple Key Support:** Automatically iterates through available SSH private keys (RSA and Ed25519) in your `~/.ssh` directory.
* **Production-Grade Security:** Uses `knownhosts` to verify the server's fingerprint, protecting against Man-in-the-Middle (MITM) attacks. 
* **Session Management:** Creates a dedicated SSH session channel to execute remote commands and capture the output.
* **Strict Error Handling:** Gracefully handles missing keys, failed authentications, and execution errors.

---

## 🔒 Security Warning
*Insecure Host Key Verification (like `ssh.InsecureIgnoreHostKey()`) is a common anti-pattern used in tutorials.* 
This project implements **secure host key verification** by reading your local `~/.ssh/known_hosts` file. The connection will fail if the server's key is not found or does not match.

---

## 🚀 Getting Started

### Prerequisites
* **Go:** Version 1.18 or higher installed on your machine.
* **SSH Server:** A running SSH daemon on the target machine (for local testing, `sshd` must be running on your machine).
* **SSH Keys:** A valid private/public key pair.

### Setup Instructions (Local Testing)
To test this project locally on a Fedora/Linux machine, follow these steps:

**1. Start the SSH Server**
```bash
sudo systemctl start sshd
sudo systemctl enable sshd
```

**2. Generate an SSH Key (if you don't have one)**
```bash
ssh-keygen -t ed25519 -C "local-test"
```

**3. Authorize Your Key**
Copy your public key to the authorized keys file so the server trusts you.
```bash
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

**4. Trust the Local Server's Host Key**
Since this program uses strict host key verification, you must add your local server to your `known_hosts` file.
```bash
ssh-keygen -R localhost
ssh-keyscan -H localhost >> ~/.ssh/known_hosts
```

---

## 💻 Usage

1. Clone or download the project files to your local machine.
2. Navigate to the project directory:
   ```bash
   cd ssh-golang
   ```
3. Run the program:
   ```bash
   go run main.go
   ```

### Expected Output
If everything is configured correctly, you will see a successful login message followed by the output of the `cat /etc/os-release` command from your local machine.

<img width="1299" height="774" alt="Screenshot From 2026-11-19 22-06-15" src="https://github.com/user-attachments/assets/f51170a8-6635-4d1d-86b1-7761c68d62fb" />

---

## ⚙️ How It Works (Under the Hood)

1. **Host Key Verification:** The program uses `knownhosts.New()` to load your trusted hosts. It will reject any server whose fingerprint doesn't match the file.
2. **Key Parsing:** It reads the raw bytes of your private keys (`os.ReadFile`) and turns them into cryptographic `Signer` objects (`ssh.ParsePrivateKey`).
3. **The Dial:** It attempts to connect to the server using `ssh.Dial()`. If the first key fails, it loops and tries the next one.
4. **The Session:** Once the encrypted client connection (`sshClient`) is established, it opens a `session`. 
5. **Execution:** Finally, it runs a command using `session.CombinedOutput()`, which captures both standard output and standard error, and prints the result to the screen.

---

## 🛠️ Customization

You can easily modify the `main.go` file to connect to a different server or run a different command:

* **Change the target:** Modify `serverAdd := "localhost:22"` to your remote server's IP (e.g., `"191.168.1.50:22"`).
* **Change the user:** Modify `User: "your-username"` to the remote username.
* **Change the command:** Modify `session.CombinedOutput("cat /etc/os-release")` to run any bash command you want (e.g., `"whoami"` or `"ls -la"`).

---

## 📚 Key Learnings & Concepts

Building this project from scratch was a deep dive into network protocols, cryptography, and Go's ecosystem. Here is a breakdown of the core concepts I mastered during this journey:

### 1. Go Module Management (`go mod`)
*   **The Dependency Trap:** I learned that modern Go requires a `go.mod` file to manage external dependencies. You cannot simply run `go get` in an empty directory; the compiler needs to know the root of the project.
*   **The Module Proxy:** I hit network timeouts trying to reach the default Google proxy (`proxy.golang.org`). I learned how to configure the `GOPROXY` environment variable (using a mirror like `goproxy.cn`) to successfully download the `golang.org/x/crypto` library.
*   **Dependency Tracking:** I observed how `go get` automatically updates the `go.mod` and `go.sum` files, ensuring reproducible builds by locking the exact version of the SSH library used.

### 2. Asymmetric Cryptography: Keys vs. Signers
*   **Public vs. Private Keys:** I learned that SSH uses asymmetric encryption. The **Private Key** (kept secret on the client) and the **Public Key** (placed in the server's `authorized_keys` file) form a mathematically linked pair.
*   **The Role of the Signer:** I discovered that the SSH library cannot use a raw text file directly. The `os.ReadFile()` function reads the key as raw bytes, which must be parsed using `ssh.ParsePrivateKey()`. This transforms the raw data into an `ssh.Signer`—a cryptographic object capable of signing authentication challenges sent by the server to prove my identity without ever transmitting the private key itself.

### 3. Host Key Verification & Preventing MITM Attacks
*   **The Insecure Trap:** I noticed many tutorials use `ssh.InsecureIgnoreHostKey()`. I learned that this is a critical security flaw. It blindly trusts any server, making the client vulnerable to **Man-in-the-Middle (MITM) attacks**, where an attacker could intercept the connection and steal credentials.
*   **Trusting the Fingerprint:** I implemented production-grade security using the `knownhosts` package. By reading the local `~/.ssh/known_hosts` file, the client verifies the server's public key fingerprint. If the server's identity doesn't match the saved fingerprint, the connection is instantly dropped.

### 4. The SSH Architecture: Clients vs. Sessions
*   **The Tunnel (Client):** I learned that `ssh.Dial()` establishes the main encrypted transport layer—the secure tunnel between the client and server. This is a single, persistent connection.
*   **The Channels (Sessions):** SSH is designed for multiplexing. Once the `sshClient` tunnel is established, I can open multiple `Sessions` within it. A session is essentially a dedicated channel for a specific task (like opening a shell, running a command, or transferring a file).
*   **Resource Cleanup:** I used Go's `defer` keyword to ensure that both the session and the client connection are gracefully closed when the program exits, preventing resource leaks on both the client and server side.

### 5. Remote Execution & Output Handling
*   **Executing Commands:** I learned how to use `session.CombinedOutput()` to send a bash command to the remote server and wait for it to finish.
*   **Handling `[]byte`:** I learned that network data is returned as raw bytes (`[]byte`). I had to convert this back to a human-readable string using `string(output)` to print the server's OS information to the local terminal.
*   **Graceful Error Handling:** I implemented a loop that gracefully handles missing keys (e.g., if an `id_rsa` file doesn't exist, it skips it and tries `id_ed25519`). I also used `log.Fatalf` to cleanly exit the program if the connection or command execution ultimately fails. 

### 6. Linux System Administration (Bonus)
*   To test the client, I had to configure my local Fedora machine as a target. I learned how to start and enable the SSH daemon (`systemctl start/enable sshd`), generate modern Ed25519 keys (`ssh-keygen`), and enforce strict file permissions (`chmod 700/600`) to satisfy OpenSSH's security requirements.
