// openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

// go mod init go-email-server

// go get github.com/emersion/go-imap
// go get github.com/emersion/go-imap/backend/memory
// go get github.com/emersion/go-imap/server

package main

import (
    "crypto/tls"
    "log"
    "net"
    "net/http"
)

const (
    smtpHost = "localhost"
    smtpPort = "2525"
    imapHost = "localhost"
    imapPort = "993"
    httpPort = "8443"
    certFile = "cert.pem"
    keyFile  = "key.pem"
)

func main() {
    go startSMTPServer()
    go startIMAPServer()
    go startAdminPanel()
    select {} // block forever
}

func startSMTPServer() {
    addr := net.JoinHostPort(smtpHost, smtpPort)
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load TLS certificate: %v", err)
    }
    config := &tls.Config{Certificates: []tls.Certificate{cert}}

    l, err := tls.Listen("tcp", addr, config)
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
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to load TLS certificate: %v", err)
    }
    config := &tls.Config{Certificates: []tls.Certificate{cert}}

    s := newIMAPServer()
    l, err := tls.Listen("tcp", addr, config)
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
    log.Println("Admin panel running on https://localhost:8443")
    log.Fatal(http.ListenAndServeTLS(":"+httpPort, certFile, keyFile, nil))
}