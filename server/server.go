package server

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	pb2 "github.com/Rorical/SMTPForward/pb"
	smtp "github.com/Rorical/go-smtp"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"
	"time"
)

// http server for sending smtp request

type Sender struct {
	Domain       string
	DkimSelector string
	DkimKey      *rsa.PrivateKey
}

func NewSender(domain, selector string, keyPath string) *Sender {
	fileContent, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(fileContent)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return &Sender{
		Domain:       domain,
		DkimSelector: selector,
		DkimKey:      privateKey.(*rsa.PrivateKey),
	}
}

func rebuildEmail(message *mail.Message) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	for key, values := range message.Header {
		for _, value := range values {
			_, err = buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = buf.WriteString("\r\n")
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadFrom(message.Body)
	return &buf, err
}

func (s *Sender) trySend(msg *bytes.Buffer, sender, recipient, host string) error {
	cli, err := smtp.Dial(host+":25", s.Domain)
	if err != nil {
		return err
	}

	defer cli.Close()

	if err = cli.Mail(sender, nil); err != nil {
		return err
	}

	if err = cli.Rcpt(recipient, nil); err != nil {
		return err
	}

	w, err := cli.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(msg.Bytes()); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	err = cli.Quit()
	return err
}

func (s *Sender) SendMail(message *mail.Message, sender, recipient string) error {
	if message.Header.Get("Date") == "" {
		message.Header["Date"] = []string{time.Now().Format(time.RFC1123Z)}
	}
	if message.Header.Get("Message-Id") == "" {
		message.Header["Message-Id"] = []string{fmt.Sprintf("<%s@%s>", uuid.New().String(), s.Domain)}
	}

	unsignedMsg, err := rebuildEmail(message)
	if err != nil {
		return err
	}
	options := &dkim.SignOptions{
		Domain:   s.Domain,
		Selector: s.DkimSelector,
		Signer:   s.DkimKey,
	}
	var signedMsg bytes.Buffer
	if err = dkim.Sign(&signedMsg, unsignedMsg, options); err != nil {
		return err
	}

	recipientAddr := strings.Split(recipient, "@")[1]
	mxs, err := net.LookupMX(recipientAddr)
	if err != nil {
		return err
	}
	for _, mx := range mxs {
		err = s.trySend(&signedMsg, sender, recipient, mx.Host)
		if err == nil {
			return nil
		} else {
			log.Println(err)
		}
	}
	return err
}

type SMTPForwardServer struct {
	pb2.SMTPForwardServer
	Token  string
	Sender *Sender
}

type Config struct {
	Domain       string `json:"domain"`
	DkimSelector string `json:"dkim_selector"`
	DkimKeyPath  string `json:"dkim_key_path"`
	Address      string `json:"address"`
	TLSCertPath  string `json:"tls_cert_path"`
	TLSKeyPath   string `json:"tls_key_path"`
	Token        string `json:"token"`
}

func NewSMTPForwardServer(cfg *Config) *SMTPForwardServer {
	sender := NewSender(cfg.Domain, cfg.DkimSelector, cfg.DkimKeyPath)
	return &SMTPForwardServer{Sender: sender, Token: cfg.Token}
}

func (s *SMTPForwardServer) SendSMTP(ctx context.Context, data *pb2.SMTPData) (*pb2.SMTPResult, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Must Provide Token")
	}

	if len(md["authorization"]) == 0 || md["authorization"][0] != s.Token {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid Token")
	}

	sender := strings.Trim(strings.TrimSpace(data.GetFrom()), "<>")

	message, err := mail.ReadMessage(strings.NewReader(data.GetData()))
	if err != nil {
		return &pb2.SMTPResult{Success: false}, err
	}

	for _, recipient := range data.GetRecipients() {
		err = s.Sender.SendMail(message, sender, strings.Trim(strings.TrimSpace(recipient), "<>"))
		if err != nil {
			return &pb2.SMTPResult{Success: false}, err
		}
	}
	return &pb2.SMTPResult{Success: true}, nil
}

func Serve(cfg *Config) error {
	creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertPath, cfg.TLSKeyPath)

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	defer lis.Close()

	s := grpc.NewServer(grpc.Creds(creds))
	pb2.RegisterSMTPForwardServer(s, NewSMTPForwardServer(cfg))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
