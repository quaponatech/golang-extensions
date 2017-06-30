package debug

import (
	"log"
	"os"
	"runtime"

	"github.com/quaponatech/golang-extensions/rtti"
)

//PrintStackTrace prints upto
// 2048 bytes of the actual runtime information stack
// to the given file descriptor e.g. os.Stdout
func PrintStackTrace(fd *os.File) int {
	if nil == fd {
		return 0
	}
	var buffer [2048]byte

	writtenBytes := runtime.Stack(buffer[:], false)
	//fd.WriteString("------------------------------STACK------------------------------\n")
	fd.Write(buffer[:writtenBytes])

	return writtenBytes
}

//PrettyStackTraceString prints upto
// 2048 bytes of the actual runtime information stack
// to the given file descriptor e.g. os.Stdout
func PrettyStackTraceString(size int) string {
	if 0 == size {
		return ""
	}
	buffer := make([]byte, size)

	writtenBytes := runtime.Stack(buffer[:], false)
	s := "------------------------------STACK------------------------------\n"
	s += string(buffer[:writtenBytes])
	s += "\n-------------------------------END-------------------------------"

	return s
}

// PrintCalledFunc prints result of GetFunctionName to log with called tag
func PrintCalledFunc(i interface{}) {
	log.Println(rtti.GetFunctionName(i) + " called!")
}
