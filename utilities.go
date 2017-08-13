package utilities

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// ParseJTWClaims parses a JWT token for value accessing
func ParseJTWClaims(tokenString string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWT_INVALID", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

type ByteReader struct {
	Reader   *bytes.Reader
	Position uint64
}

func (br *ByteReader) ReadNBytes(n uint64) ([]byte, error) {
	bslice := []byte{}

	for i := uint64(0); i < n; i++ {
		b, err := br.Reader.ReadByte()
		if err != nil {
			return []byte{}, err
		}
		bslice = append(bslice, b)
	}

	br.Position += n

	return bslice, nil
}

func (br *ByteReader) ReadByte() (byte, error) {
	br.Position++
	return br.Reader.ReadByte()
}

func (br *ByteReader) Finished() bool {
	return int(br.Position) >= br.Reader.Len()
}
