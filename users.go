package main

import (
    "sync"
)

type User struct {
    Username string
    Password string
    Email    string
}

var (
    users     = make(map[string]User)
    usersLock sync.RWMutex
)

func addUser(username, password, email string) {
    usersLock.Lock()
    defer usersLock.Unlock()
    users[username] = User{Username: username, Password: password, Email: email}
}

func removeUser(username string) {
    usersLock.Lock()
    defer usersLock.Unlock()
    delete(users, username)
}

func listUsers() []User {
    usersLock.RLock()
    defer usersLock.RUnlock()
    var userList []User
    for _, user := range users {
        userList = append(userList, user)
    }
    return userList
}

func getUserEmails(email string) []Email {
    usersLock.RLock()
    defer usersLock.RUnlock()
    return mailStore[email]
}