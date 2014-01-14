// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/user"
	"strings"

	"code.google.com/p/go.crypto/ssh"

	"launchpad.net/juju-core/utils"
)

// GoCryptoClient is an implementation of Client that
// uses the embedded go.crypto/ssh SSH client.
//
// GoCryptoClient is intentionally limited in the
// functionality that it enables, as it is currently
// intended to be used only for non-interactive command
// execution.
type GoCryptoClient struct {
	signers []ssh.Signer
}

// NewGoCryptoClient creates a new GoCryptoClient.
//
// If no signers are specified, NewGoCryptoClient will
// use the private key generated by LoadClientKeys.
func NewGoCryptoClient(signers ...ssh.Signer) (*GoCryptoClient, error) {
	return &GoCryptoClient{signers: signers}, nil
}

// Command implements Client.Command.
func (c *GoCryptoClient) Command(host string, command []string, options *Options) *Cmd {
	shellCommand := utils.CommandString(command...)
	signers := c.signers
	if len(signers) == 0 {
		signers = privateKeys()
	}
	return &Cmd{impl: &goCryptoCommand{
		signers: signers,
		host:    host,
		command: shellCommand,
	}}
}

// Copy implements Client.Copy.
//
// Copy is currently unimplemented, and will always return an error.
func (c *GoCryptoClient) Copy(source, dest string, options *Options) error {
	return fmt.Errorf("Copy is not implemented")
}

type goCryptoCommand struct {
	signers []ssh.Signer
	host    string
	command string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	conn    *ssh.ClientConn
	sess    *ssh.Session
}

func (c *goCryptoCommand) ensureSession() (*ssh.Session, error) {
	if c.sess != nil {
		return c.sess, nil
	}
	if len(c.signers) == 0 {
		return nil, fmt.Errorf("no private keys available")
	}
	username, host := splitUserHost(c.host)
	if username == "" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("getting current user: %v", err)
		}
		username = currentUser.Username
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(keyring{c.signers}),
		},
	}
	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}
	sess, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, err
	}
	c.conn = conn
	c.sess = sess
	c.sess.Stdin = c.stdin
	c.sess.Stdout = c.stdout
	c.sess.Stderr = c.stderr
	return sess, nil
}

func (c *goCryptoCommand) Start() error {
	sess, err := c.ensureSession()
	if err != nil {
		return err
	}
	if c.command == "" {
		return sess.Shell()
	}
	return sess.Start(c.command)
}

func (c *goCryptoCommand) Close() error {
	if c.sess == nil {
		return nil
	}
	err0 := c.sess.Close()
	err1 := c.conn.Close()
	if err0 == nil {
		err0 = err1
	}
	c.sess = nil
	c.conn = nil
	return err0
}

func (c *goCryptoCommand) Wait() error {
	if c.sess == nil {
		return fmt.Errorf("Command has not been started")
	}
	err := c.sess.Wait()
	c.Close()
	return err
}

func (c *goCryptoCommand) Kill() error {
	if c.sess == nil {
		return fmt.Errorf("Command has not been started")
	}
	return c.sess.Signal(ssh.SIGKILL)
}

func (c *goCryptoCommand) SetStdio(stdin io.Reader, stdout, stderr io.Writer) {
	c.stdin = stdin
	c.stdout = stdout
	c.stderr = stderr
}

func (c *goCryptoCommand) StdinPipe() (io.WriteCloser, io.Reader, error) {
	sess, err := c.ensureSession()
	if err != nil {
		return nil, nil, err
	}
	wc, err := sess.StdinPipe()
	return wc, sess.Stdin, err
}

func (c *goCryptoCommand) StdoutPipe() (io.ReadCloser, io.Writer, error) {
	sess, err := c.ensureSession()
	if err != nil {
		return nil, nil, err
	}
	wc, err := sess.StdoutPipe()
	return ioutil.NopCloser(wc), sess.Stdout, err
}

func (c *goCryptoCommand) StderrPipe() (io.ReadCloser, io.Writer, error) {
	sess, err := c.ensureSession()
	if err != nil {
		return nil, nil, err
	}
	wc, err := sess.StderrPipe()
	return ioutil.NopCloser(wc), sess.Stderr, err
}

// keyring implements ssh.ClientKeyring
type keyring struct {
	signers []ssh.Signer
}

func (k keyring) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.signers) {
		// nil key marks the end of the keyring; must not return an error.
		return nil, nil
	}
	return k.signers[i].PublicKey(), nil
}

func (k keyring) Sign(i int, rand io.Reader, data []byte) ([]byte, error) {
	if i < 0 || i >= len(k.signers) {
		return nil, fmt.Errorf("no key at position %d", i)
	}
	return k.signers[i].Sign(rand, data)
}

func splitUserHost(s string) (user, host string) {
	userHost := strings.SplitN(s, "@", 2)
	if len(userHost) == 2 {
		return userHost[0], userHost[1]
	}
	return "", userHost[0]
}
