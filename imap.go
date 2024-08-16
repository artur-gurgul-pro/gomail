package main

import (
    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/backend/memory"
    "github.com/emersion/go-imap/server"
    "log"
)

func newIMAPServer() *server.Server {
    be := memory.New()
    s := server.New(be)
    s.Addr = ":143"
    s.AllowInsecureAuth = true
    return s
}

func getUserMailbox(email string) *memory.Mailbox {
    be := memory.New()
    user := &memory.User{Username: email, Mailbox: memory.NewMailbox()}
    be.Users[email] = user
    return user.Mailbox
}