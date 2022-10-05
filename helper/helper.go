package helper

import (
	"fmt"
	"github.com/joho/godotenv"
	"math/rand"
	"net"
	"os"
	"time"
)

func ReadEnvFile(path string) (map[string]string, error) {
	return godotenv.Read(path)
}

// FreePort asks the kernel for a free open port that is ready to use.
func FreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RandomizedName(name string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := 4
	var suffix string
	for i := 0; i < digits; i++ {
		suffix += string(rune(rnd.Intn(10) + 48))
	}

	return fmt.Sprintf("%s-%s", name, suffix)
}
