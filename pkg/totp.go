package pkg

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// Code .
func Code(secret string) string {
	// Generate the time based on 30 second time periods.
	nowRaw := time.Now()
	seconds := 0
	if nowRaw.Second() > 30 {
		seconds = 30
	}
	now := time.Date(nowRaw.Year(), nowRaw.Month(), nowRaw.Day(), nowRaw.Hour(), nowRaw.Minute(), seconds, 0, nowRaw.Location())

	// Generate the challenge.
	t0 := int64(now.Unix() / 30)
	challenge, err := computeCode(secret, t0)
	if err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("%06d", challenge)
}

// ComputeCode returns the 6 digit code based on the secret and the time.
func computeCode(secret string, value int64) (int, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return 0, err
	}

	hash := hmac.New(sha1.New, key)
	err = binary.Write(hash, binary.BigEndian, value)
	if err != nil {
		return 0, err
	}
	h := hash.Sum(nil)

	offset := h[19] & 0x0f

	truncated := binary.BigEndian.Uint32(h[offset : offset+4])

	truncated &= 0x7fffffff
	code := truncated % 1000000

	return int(code), nil
}

// UUID generates a UUID for use as an ID.
func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
