//go:build !linux

package threads

func getFileDescriptorsCount() int {
	return -1 // not supported
}
