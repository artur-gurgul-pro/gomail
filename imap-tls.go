package main

import (
    "crypto/tls"
    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/backend/memory"
    "github.com/emersion/go-imap/server"
    "log"
)

func newIMAPServer() *server.Server {
    be := memory.New()
    s := server.New(be)
    s.Addr = ":993"
    s.AllowInsecureAuth = false
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load TLS certificate: %v", err)
    }
    s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
    return s
}

func getUserMailbox(email string) *memory.Mailbox {
    be := memory.New()
    user := &memory.User{Username: email, Mailbox: memory.NewMailbox()}
    be.Users[email] = user
    return user.Mailbox
}