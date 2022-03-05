package util

import (
	"github.com/aofei/mimesniffer"
	"io/fs"
)

func GetTypeFromFile(f fs.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	return mimesniffer.Sniff(buffer), nil
}
