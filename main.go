package main

import (
    "log"
    "net"
    "net/http"
)

const (
    smtpHost = "localhost"
    smtpPort = "2525"
    imapHost = "localhost"
    imapPort = "143"
    httpPort = "8080"
)

func main() {
    go startSMTPServer()
    go startIMAPServer()
    go startAdminPanel()
    select {} // block forever
}

func startSMTPServer() {
    addr := net.JoinHostPort(smtpHost, smtpPort)
    l, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Failed to bind to port %s: %v", smtpPort, err)
    }
    defer l.Close()
    log.Printf("SMTP server listening on %s", addr)

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }

        go handleSMTPConnection(conn)
    }
}

func startIMAPServer() {
    addr := net.JoinHostPort(imapHost, imapPort)
    s := newIMAPServer()
    l, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Failed to bind to port %s: %v", imapPort, err)
    }
    defer l.Close()
    log.Printf("IMAP server listening on %s", addr)

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }

        go s.Serve(conn)
    }
}

func startAdminPanel() {
    http.HandleFunc("/add_user", addUserHandler)
    http.HandleFunc("/remove_user", removeUserHandler)
    http.HandleFunc("/list_users", listUsersHandler)
    http.HandleFunc("/send_email", sendEmailHandler)
    http.HandleFunc("/get_emails", getEmailsHandler)
    log.Println("Admin panel running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}