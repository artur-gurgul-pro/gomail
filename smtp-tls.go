package main

import (
    "bufio"
    "crypto/tls"
    "encoding/base64"
    "log"
    "net"
    "strings"
)

func handleSMTPConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    writeResponse(writer, "220 Welcome to Go SMTP Server")

    var email Email
    var authenticated bool
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Failed to read from connection: %v", err)
            return
        }
        line = strings.TrimSpace(line)
        log.Printf("Received: %s", line)

        if strings.HasPrefix(line, "EHLO") {
            writeResponse(writer, "250-Hello")
            writeResponse(writer, "250-STARTTLS")
            writeResponse(writer, "250 AUTH LOGIN")
        } else if strings.HasPrefix(line, "STARTTLS") {
            writeResponse(writer, "220 Ready to start TLS")
            tlsConn := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{}})
            handleSMTPConnection(tlsConn)
            return
        } else if strings.HasPrefix(line, "AUTH LOGIN") {
            writeResponse(writer, "334 VXNlcm5hbWU6") // "Username:" in base64
            username, _ := reader.ReadString('\n')
            decodedUsername, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(username))
            writeResponse(writer, "334 UGFzc3dvcmQ6") // "Password:" in base64
            password, _ := reader.ReadString('\n')
            decodedPassword, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(password))
            if authenticateUser(string(decodedUsername), string(decodedPassword)) {
                authenticated = true
                writeResponse(writer, "235 Authentication successful")
            } else {
                writeResponse(writer, "535 Authentication failed")
                return
            }
        } else if strings.HasPrefix(line, "MAIL FROM:") && authenticated {
            email.From = strings.TrimPrefix(line, "MAIL FROM:")
            writeResponse(writer, "250 OK")
        } else if strings.HasPrefix(line, "RCPT TO:") && authenticated {
            email.To = strings.TrimPrefix(line, "RCPT TO:")
            writeResponse(writer, "250 OK")
        } else if strings.HasPrefix(line, "DATA") && authenticated {
            writeResponse(writer, "354 Start mail input; end with <CRLF>.<CRLF>")
            var dataLines []string
            for {
                dataLine, err := reader.ReadString('\n')
                if err != nil {
                    log.Printf("Failed to read data from connection: %v", err)
                    return
                }
                dataLine = strings.TrimSpace(dataLine)
                if dataLine == "." {
                    break
                }
                dataLines = append(dataLines, dataLine)
            }
            email.Body = strings.Join(dataLines, "\n")
            mailStore[email.To] = append(mailStore[email.To], email)
            writeResponse(writer, "250 OK")
        } else if strings.HasPrefix(line, "QUIT") {
            writeResponse(writer, "221 Bye")
            return
        } else {
            writeResponse(writer, "500 Unrecognized command")
        }
    }
}

func writeResponse(writer *bufio.Writer, response string) {
    writer.WriteString(response + "\r\n")
    writer.Flush()
}

func authenticateUser(username, password string) bool {
    // Implement your authentication logic here
    // For now, let's assume a simple check
    return username == "user" && password == "pass"
}