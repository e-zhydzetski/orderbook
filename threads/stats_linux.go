package threads

import (
	"fmt"
	"log"
	"os"
)

func getFileDescriptorsCount() int {
	fd, err := os.ReadDir(fmt.Sprintf("/proc/%d/fd", os.Getpid()))
	if err != nil {
		log.Println(err)
		return -1
	}
	return len(fd)
}
