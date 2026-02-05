package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5String calculates MD5 hash of a string
func MD5String(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MD5File calculates MD5 hash of a file
func MD5File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}