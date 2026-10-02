package auth

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
)

func GenerateToken() (string, string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil { return "", "", err }
    plain := hex.EncodeToString(b)
    return plain, HashToken(plain), nil
}

func HashToken(token string) string {
    sum := sha256.Sum256([]byte(token))
    return hex.EncodeToString(sum[:])
}