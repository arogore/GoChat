package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "strings"
)

var port string = "1287"

/**
*  Connects to the host and prompts for a username
*  Starts waiting for messages and handles sending its' own
*/
func main() {
  host := "localhost"
  fmt.Println("Connecting to " + host + " on port " + port)

  conn, err := net.Dial("tcp", host+":"+port)
  if err != nil {
    fmt.Println("Error connecting" + err.Error())
    return
  }

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter a username: ")
  username, _ := reader.ReadString('\n')
  fmt.Println()
  username = strings.TrimSpace(username)
  conn.Write([]byte(string(username)))
  go waitMessages(conn)

  for {
    txt, _ := reader.ReadString('\n')
    txt = strings.TrimSpace(txt)
    if txt == "quit" {
      break
    }
    conn.Write([]byte(txt))
  }
  conn.Close()
}

/**
*  Waits for a message to be received and outputs it
*  @param conn - Connection we're waiting to receive from
*/
func waitMessages(conn net.Conn) {
  for {
    var buffer [512]byte
    len, _ := conn.Read(buffer[0:])
    fmt.Println("\n" + string(buffer[0:len]) + "\n")
  }
}
