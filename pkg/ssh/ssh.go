package ssh

import (
	"bytes"
	"fmt"
	"net"
	"errors"
	"time"
	"golang.org/x/crypto/ssh"
)


type SSHTunnel struct {
	Config  *ssh.ClientConfig
	Host    string
	Port	string
	Network string
	session *ssh.Session
	client  *ssh.Client
	Timeout time.Duration
}

func NewSSHTunnel(host, user, password string) (*SSHTunnel, error) {
	return makeSSHTunnel(host, user, password)
}

// 默认 150 秒超时
func makeSSHTunnel(host, user, password string) (*SSHTunnel, error) {
	config := ssh.ClientConfig{
		User:			 user,
		Auth:			 []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &SSHTunnel{
		Config: &config,
		Host: host,
		Port: "22",
		Network: "tcp",
		Timeout: 150 * time.Second,
	}, nil
}
func (s *SSHTunnel) SetTimeout(timeout time.Duration) {
	s.Timeout = timeout
}

func (s *SSHTunnel) Dial() error {
	conn, err := net.DialTimeout(s.Network, net.JoinHostPort(s.Host, s.Port), s.Timeout)
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	c, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(s.Host, s.Port), s.Config)
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Time{})
	s.client = ssh.NewClient(c, chans, reqs)
	return nil
}


func (s *SSHTunnel) RunCommond(cmd string) (string, string, int, error) {
	var err error
	s.session, err = s.client.NewSession()
	if err != nil {
		return "", "", 0, fmt.Errorf("Error creating session to %s@%s: '%v'", s.Config.User, s.Host, err)
	}
	//defer session.Close()

	// Run the command.
	code := 0
	var bout, berr bytes.Buffer
	s.session.Stdout, s.session.Stderr = &bout, &berr
	if err = s.session.Run(cmd); err != nil {
		// Check whether the command failed to run or didn't complete.
		if exiterr, ok := err.(*ssh.ExitError); ok {
			// If we got an ExitError and the exit code is nonzero, we'll
			// consider the SSH itself successful (just that the command run
			// errored on the host).
			if code = exiterr.ExitStatus(); code != 0 {
				err = nil
			}
		} else {
			// Some other kind of error happened (e.g. an IOError); consider the
			// SSH unsuccessful.
			err = fmt.Errorf("Failed running `%s` on %s@%s: '%v'", cmd, s.Config.User, s.Host, err)
		}
	}
	return bout.String(), berr.String(), code, err
}

func (s *SSHTunnel) Close() error {
	if s.session != nil {
		s.session.Close()
	}
	if s.client == nil {
		return errors.New("Cannot close tunnel. Tunnel was not opened.")
	}
	if err := s.client.Close(); err != nil {
		return err
	}
	return nil
}
