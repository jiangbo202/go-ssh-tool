package utils

import (
    "fmt"
    "io"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
)

type SSHTerminal struct {
    Session *ssh.Session
    exitMsg string
    stdout  io.Reader
    stdin   io.Writer
    stderr  io.Reader
}

func (t *SSHTerminal) updateTerminalSize() {

    go func() {
        // SIGWINCH is sent to the process when the window size of the terminal has
        // changed.
        sigwinchCh := make(chan os.Signal, 1)
        signal.Notify(sigwinchCh, syscall.SIGWINCH)

        fd := int(os.Stdin.Fd())
        termWidth, termHeight, err := terminal.GetSize(fd)
        if err != nil {
            log.Printf("遇到错误: %s\n", err)
        }

        for {
            // The client updated the size of the local PTY. This change needs to occur
            // on the server side PTY as well.
            sigwinch := <-sigwinchCh
            if sigwinch == nil {
                return
            }
            currTermWidth, currTermHeight, err := terminal.GetSize(fd)
            if err != nil {
                log.Printf("遇到错误 Get Size Failed: %s\n", err)
            }
            // Terminal size has not changed, don't do anything.
            if currTermHeight == termHeight && currTermWidth == termWidth {
                continue
            }

            err = t.Session.WindowChange(currTermHeight, currTermWidth)
            if err != nil {
                log.Printf("遇到错误 Unable To Send Window Size Change: %s\n", err)
                continue
            }

            termWidth, termHeight = currTermWidth, currTermHeight

        }
    }()

}

func (t *SSHTerminal) interactiveSession() error {

    defer func() {
        if t.exitMsg == "" {
            fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
        } else {
            fmt.Fprintln(os.Stdout, t.exitMsg)
        }
    }()

    fd := int(os.Stdin.Fd())
    state, err := terminal.MakeRaw(fd)
    if err != nil {
        return err
    }
    defer func() {
        _ = terminal.Restore(fd, state)
    }()

    termWidth, termHeight, err := terminal.GetSize(fd)
    if err != nil {
        return err
    }

    termType := os.Getenv("TERM")
    if termType == "" {
        termType = "xterm-256color"
    }

    err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
    if err != nil {
        return err
    }

    t.updateTerminalSize()

    t.stdin, err = t.Session.StdinPipe()
    if err != nil {
        return err
    }
    t.stdout, err = t.Session.StdoutPipe()
    if err != nil {
        return err
    }
    t.stderr, _ = t.Session.StderrPipe()

    go func() {
        _, _ = io.Copy(os.Stderr, t.stderr)
    }()
    go func() {
        _, _ = io.Copy(os.Stdout, t.stdout)
    }()
    go func() {
        buf := make([]byte, 128)
        for {
            n, err := os.Stdin.Read(buf)
            if err != nil {
                log.Printf("遇到错误: %s\n", err)
                return
            }
            if n > 0 {
                _, err = t.stdin.Write(buf[:n])
                if err != nil {
                    log.Printf("遇到错误: %s\n", err)
                    t.exitMsg = err.Error()
                    return
                }
            }
        }
    }()

    err = t.Session.Shell()
    if err != nil {
        return err
    }
    err = t.Session.Wait()
    if err != nil {
        return err
    }
    return nil
}

func New(cli *ssh.Client) error {

    session, err := cli.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()

    s := SSHTerminal{
        Session: session,
    }

    return s.interactiveSession()
}

func (c Cli) NewTerminal() error {
    if c.client == nil {
        if err := c.Connect(); err != nil {
            return err
        }
    }

    err := New(c.client)
    if err != nil {
        return fmt.Errorf("new terminal client error: %w", err)
    }
    return nil
}
