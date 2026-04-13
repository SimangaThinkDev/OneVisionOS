package monitor

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Checksum computes the SHA256 hash of a file
func Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyIntegrity checks if the file matches the expected hash
func VerifyIntegrity(filePath string, expectedHash string) (bool, error) {
	actualHash, err := Checksum(filePath)
	if err != nil {
		return false, err
	}
	return actualHash == expectedHash, nil
}
