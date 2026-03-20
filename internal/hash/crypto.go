package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileSHA256 computes the SHA256 hash of a file.
func FileSHA256(filePath string) (string, error) {
	return computeHash(filePath, sha256.New())
}

// FileSHA1 computes the SHA1 hash of a file.
func FileSHA1(filePath string) (string, error) {
	return computeHash(filePath, sha1.New())
}

// FileMD5 computes the MD5 hash of a file.
func FileMD5(filePath string) (string, error) {
	return computeHash(filePath, md5.New())
}

func computeHash(filePath string, hasher io.Writer) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// This is a bit hacky but works since all standard hash.Hash implementations
	// support Sum(nil). We cast to an interface that has Sum.
	type sumGetter interface {
		Sum(b []byte) []byte
	}
	return hex.EncodeToString(hasher.(sumGetter).Sum(nil)), nil
}

// StringSHA256 computes the SHA256 hash of a direct string (e.g., manifest content).
func StringSHA256(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}
