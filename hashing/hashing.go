package hashing

import (
	"encoding/base64"
	"fmt"
	"io"

	"crypto/hmac"
	"crypto/sha256"
)

// BuildPasswordHash ...
func BuildPasswordHash(username, password, salt1, salt2 string) (string, error) {
	if "" == username || "" == password {
		return "", fmt.Errorf("Invalid username or password")
	} else if "" == salt1 || "" == salt2 {
		// Empty salts are okay.
	}

	hashedpassword := sha256.New()
	io.WriteString(hashedpassword, password)

	var pwsha = fmt.Sprintf("%x", hashedpassword.Sum(nil))

	io.WriteString(hashedpassword, salt1)
	io.WriteString(hashedpassword, username)
	io.WriteString(hashedpassword, salt2)
	io.WriteString(hashedpassword, pwsha)

	return fmt.Sprintf("%x", hashedpassword.Sum(nil)), nil
}

// BuildHashWithOneSaltAndHmac ...
func BuildHashWithOneSaltAndHmac(info, salt string) (string, error) {
	if "" == info {
		return "", fmt.Errorf("Invalid information")
	} else if "" == salt {
		// Empty salts are okay.
	}

	key := []byte(salt)
	//log.Printf("Bytes of Salt: %x", key)

	h := hmac.New(sha256.New, key)
	//log.Printf("HashedSecret (Key added): %x", h.Sum(nil))

	h.Write([]byte(info))
	//log.Printf("HashedSecret (Info added): %x", h.Sum(nil))

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
