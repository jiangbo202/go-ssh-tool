package utils

import (
    "fmt"
    "github.com/pkg/sftp"
    "github.com/pquerna/otp/totp"
    "golang.org/x/crypto/ssh"
    "io"
    "net"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type Cli struct {
    user         string
    pwd          string
    addr         string
    googleSecret string
    client       *ssh.Client
}

func NewCli(user, pwd, addr, googleSecret string) Cli {
    return Cli{
        user:         user,
        pwd:          pwd,
        addr:         addr,
        googleSecret: googleSecret,
    }
}

// Connect 连接远程服务器
func (c *Cli) Connect() error {
    var authList []ssh.AuthMethod
    // 这里实测仅需使用密码登录时，使用ssh.Password；密码+2fa登录时，变成交互式KeyboardInteractive(这里的密码也变成交互了，而不再是ssh.Password)
    if c.googleSecret != "" {
        authList = append(authList, ssh.KeyboardInteractive(keyboardInteractivePassword(c.googleSecret, c.pwd)))
    } else {
        authList = append(authList, ssh.Password(c.pwd))
    }

    config := &ssh.ClientConfig{
        User:            c.user,
        Auth:            authList,
        HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
    }
    client, err := ssh.Dial("tcp", c.addr, config)
    if nil != err {
        return fmt.Errorf("connect server error: %w", err)
    }
    c.client = client
    return nil
}

// Run 运行命令
func (c Cli) Run(shell string) (string, error) {
    if c.client == nil {
        if err := c.Connect(); err != nil {
            return "", err
        }
    }

    session, err := c.client.NewSession()
    if err != nil {
        return "", fmt.Errorf("create new session error: %w", err)
    }
    defer session.Close()

    buf, err := session.CombinedOutput(shell)
    return string(buf), err
}

func IsFile(f string) bool {
    fi, e := os.Stat(f)
    if e != nil {
        return false
    }
    return !fi.IsDir()

}

func PathExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    return false
}

// Scp 下载
func (c Cli) Scp(srcFileName, targetFileName string) (int64, error) {
    if c.client == nil {
        if err := c.Connect(); err != nil {
            return 0, err
        }
    }
    if targetFileName == "" {
        // 与原名相同
        targetFileName = "./" + filepath.Base(srcFileName)
    }
    if strings.HasSuffix(targetFileName, "/") { // 如果以/结尾则为目录
        targetFileName += filepath.Base(srcFileName)
    }

    if PathExists(targetFileName) && !IsFile(targetFileName) { // 如果是一个目录
        targetFileName += "/" + filepath.Base(srcFileName)
    }

    fmt.Println("targetFileName:", targetFileName)
    fmt.Println("srcFileName:", srcFileName)
    sftpClient, err := sftp.NewClient(c.client)
    if err != nil {
        return 0, fmt.Errorf("new sftp client error: %w", err)
    }
    defer sftpClient.Close()

    source, err := sftpClient.Open(srcFileName)
    if err != nil {
        return 0, fmt.Errorf("sftp client open file error: %w", err)
    }
    defer source.Close()

    target, err := os.OpenFile(targetFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return 0, fmt.Errorf("open local file error: %w", err)
    }
    defer target.Close()

    n, err := io.Copy(target, source)
    if err != nil {
        return 0, fmt.Errorf("copy file error: %w", err)
    }
    return n, nil
}

// Scp2 上传
func (c Cli) Scp2(srcFileName, targetFileName string) (int64, error) {
    if c.client == nil {
        if err := c.Connect(); err != nil {
            return 0, err
        }
    }
    if targetFileName == "" {
        // 与原名相同
        targetFileName = "./" + filepath.Base(srcFileName)
    }
    if strings.HasSuffix(targetFileName, "/") { // 如果以/结尾则为目录
        targetFileName += filepath.Base(srcFileName)
    }
    if PathExists(targetFileName) && !IsFile(targetFileName) { // 如果是一个目录
        targetFileName += "/" + filepath.Base(srcFileName)
    }

    sftpClient, err := sftp.NewClient(c.client)
    if err != nil {
        return 0, fmt.Errorf("new sftp client error: %w", err)
    }
    defer sftpClient.Close()

    fmt.Println("targetFileName:", targetFileName)
    fmt.Println("srcFileName:", srcFileName)
    source, err := sftpClient.Create(targetFileName)
    if err != nil {
        return 0, fmt.Errorf("sftp client create file error: %w", err)
    }
    defer source.Close()

    target, err := os.Open(srcFileName)
    if err != nil {
        return 0, fmt.Errorf("open local file error: %w", err)
    }
    defer target.Close()

    n, err := io.Copy(source, target)
    if err != nil {
        return 0, fmt.Errorf("copy file error: %w", err)
    }
    return n, nil
}

func keyboardInteractivePassword(secret, passwd string) ssh.KeyboardInteractiveChallenge {
    return func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
        fmt.Println("questions:", questions)
        if len(questions) > 0 {
            if questions[0] == "Verification code: " {
                token, err := totp.GenerateCode(secret, time.Now().UTC())
                fmt.Println("token:", token)
                return []string{token}, err
            } else if questions[0] == "Password: " {
                return []string{passwd}, err
            }
        }

        return []string{}, err
    }
}
