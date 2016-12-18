package main

import (
"fmt"
"net"
"strconv"
"strings"
)

var last10 []string
var connections []ConnectedUser

type ConnectedUser struct {
conn  net.Conn
uName string
}

/**
*	Starts listening for connections and sends them to a handler
*/
func main() {
fmt.Println("Listening on port 1287")
ln, err := net.Listen("tcp", ":1287")
if err != nil {
  fmt.Print("Error in the connection")
  return
}
for {
  conn, err := ln.Accept()
  if err != nil {
    fmt.Print("Error in client connection")
    } else {
      connUser := addUser(conn)
      fmt.Println("Connection amount: " + strconv.Itoa(len(connections)) + " - Connected user: " + connUser.uName)
      go handleConnection(connUser)
    }
  }
}

/**
*	Handles a user's connection
*	@param user - The user we want to handle
*/
func handleConnection(user *ConnectedUser) {
  for {
    var b [512]byte
    n, err := user.conn.Read(b[0:])
    if err != nil {
      fmt.Println("Error in reading message: " + err.Error() + " from " + user.uName)
      go broadcast(user, " has disconnected.")
      user.conn.Close()
      for i, u := range connections {
        if u == *user {
          connections = append(connections[:i], connections[i+1:]...)
          break
        }
      }
      break
    } else {
      fmt.Println(user.uName + ": " + string(b[:n]))
      go broadcast(user, string(b[:n]))
    }
  }
}

/**
*	Broadcasts a message to all users except oUser
* 	@param oUser - Original sender of the message
*	@param msg - Message we're sending
*/
func broadcast(oUser *ConnectedUser, msg string) {
  for _, user := range connections {
    if *oUser != user {
      writeMessage(&user, oUser.uName+": "+msg)
    }
  }
}

/**
*	Adds a user to the connection slice
*	@param conn - User's connection
*	@return connUser - New ConnectedUser
*/
func addUser(conn net.Conn) *ConnectedUser {
  connUser := ConnectedUser{conn, *retrieveUsername(conn)}
  connections = append(connections, connUser)
  return &connUser
}

/**
*	Retrieves a username from the client
*	@param conn - User's connection
*	@return uname - User's username
*/
func retrieveUsername(conn net.Conn) *string {
  var b [512]byte
  n, _ := conn.Read(b[0:])
  uname := strings.TrimSpace(string(b[0:n]))
  return &uname
}

/**
*	Writes a message to a user
*	@param user - The user we want to write a message to
*	@param msg - The message we want to send
*/
func writeMessage(user *ConnectedUser, msg string) {
  user.conn.Write([]byte(msg))
}
