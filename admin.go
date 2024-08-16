package main

import (
    "encoding/json"
    "log"
    "net/http"
    "net/smtp"
)

type EmailRequest struct {
    From    string
    To      string
    Subject string
    Body    string
}

func main() {
    http.HandleFunc("/add_user", addUserHandler)
    http.HandleFunc("/remove_user", removeUserHandler)
    http.HandleFunc("/list_users", listUsersHandler)
    http.HandleFunc("/send_email", sendEmailHandler)
    http.HandleFunc("/get_emails", getEmailsHandler)
    log.Println("Admin panel running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    addUser(user.Username, user.Password, user.Email)
    w.WriteHeader(http.StatusCreated)
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    removeUser(user.Username)
    w.WriteHeader(http.StatusOK)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    users := listUsers()
    if err := json.NewEncoder(w).Encode(users); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    var emailReq EmailRequest
    if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := sendEmail(emailReq.From, emailReq.To, emailReq.Subject, emailReq.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func sendEmail(from, to, subject, body string) error {
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, "your-email-password", "smtp.example.com")
    return smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, []byte(msg))
}

func getEmailsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Email parameter is required", http.StatusBadRequest)
        return
    }

    emails := getUserEmails(email)
    if err := json.NewEncoder(w).Encode(emails); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}