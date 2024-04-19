package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/Rorical/SMTPForward/pb"
	"github.com/Rorical/SMTPForward/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strings"
)

type Config struct {
	URI         string `json:"uri"`
	TLSCertPath string `json:"tls_cert_path"`
	Token       string `json:"token"`
	SMTPListen  string `json:"smtp_listen"`
}

type SMTPForwardClient struct {
	client pb.SMTPForwardClient
}

func (c *SMTPForwardClient) SendEmail(data string, from string, recipients []string) error {
	ctx := context.TODO()
	res, err := c.client.SendSMTP(ctx, &pb.SMTPData{
		Data:       data,
		From:       from,
		Recipients: recipients,
	})
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New("failed to send email")
	}
	return nil
}

func NewSMTPForward(cfg *Config) (*SMTPForwardClient, error) {
	cred, err := credentials.NewClientTLSFromFile(cfg.TLSCertPath, "")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(cfg.URI, grpc.WithTransportCredentials(cred), grpc.WithPerRPCCredentials(security.NewTokenCredential(cfg.Token)))
	if err != nil {
		return nil, err
	}
	return &SMTPForwardClient{
		client: pb.NewSMTPForwardClient(conn),
	}, nil
}

type SMTPData struct {
	data       string
	from       string
	recipients []string
}

func send(conn net.Conn, msg string) error {
	_, err := fmt.Fprintf(conn, "%s\r\n", msg)
	return err
}

func HandleSMTPSession(conn net.Conn) (*SMTPData, error) {
	defer conn.Close()
	session := &SMTPData{
		recipients: make([]string, 0),
	}
	scanner := bufio.NewScanner(conn)
	err := send(conn, "220 Welcome to the SMTPForward Agent")
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "HELO") {
			err := send(conn, "250-mail")
			if err != nil {
				return nil, err
			}
			err = send(conn, "250 Hello, NiHao")
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(text, "EHLO") {
			err = send(conn, "500 Command not recognized")
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(text, "MAIL FROM:") {
			session.from = text[10:]
			err := send(conn, "250 OK")
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(text, "RCPT TO:") {
			session.recipients = append(session.recipients, text[8:])
			err := send(conn, "250 OK")
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(text, "DATA") {
			err := send(conn, "354")
			if err != nil {
				return nil, err
			}
			session.data = ""
			for scanner.Scan() {
				dataLine := scanner.Text()
				if dataLine == "." {
					break
				}
				session.data += dataLine + "\n"
			}
			err = send(conn, "250 Ok")
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(text, "QUIT") {
			err := send(conn, "221 Bye")
			if err != nil {
				return session, err
			}
			return session, nil
		} else {
			err := send(conn, "500 Syntax error, command unrecognized")
			if err != nil {
				return nil, err
			}
		}
	}
	return session, scanner.Err()
}

func Listen(cfg *Config) error {
	listener, err := net.Listen("tcp", cfg.SMTPListen)
	if err != nil {
		return err
	}
	defer listener.Close()

	mailForward, err := NewSMTPForward(cfg)
	if err != nil {
		return err
	}

	fmt.Printf("Server is listening on %s\n", cfg.SMTPListen)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			data, err := HandleSMTPSession(conn)
			if err != nil {
				panic(err)
			}
			err = mailForward.SendEmail(data.data, data.from, data.recipients)
			if err != nil {
				panic(err)
			}
		}()
	}
}
