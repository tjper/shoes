// Package shoes is the initializer for the shoes package.
package shoes

import (
	"log"
	"os"
)

var LogErr *log.Logger

func init() {
	// Setup Logging
	file := "/home/james/go/log/shoesapi.txt"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	LogErr = log.New(f, "Error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

func main() {
	// Initialize ShoesApp
	a := App{}
	a.DbClient = DbClient{Db: Db()}

	a.Start()
}
