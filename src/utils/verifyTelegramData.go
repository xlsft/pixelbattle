package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

type TelegramData struct {
	AuthDate  int64  `json:"auth_date"`
	FirstName string `json:"first_name"`
	Hash      string `json:"hash"`
	ID        int64  `json:"id"`
	LastName  string `json:"last_name,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	Username  string `json:"username,omitempty"`
}

func (data *TelegramData) VerifyTelegramData() bool {
	var (
		token       = []byte(os.Getenv("TG_TOKEN"))
		secret      = sha256.Sum256(token)
		checkString = data.String()
		authCode    = computeHmac256([]byte(checkString), secret[:])
		hexAuthCode = hex.EncodeToString(authCode)
	)

	if hexAuthCode != data.Hash {
		return false
	}

	return true
}

func (c *TelegramData) String() string {
	s := fmt.Sprint("auth_date=", c.AuthDate, "\nfirst_name=", c.FirstName, "\nid=", c.ID)

	if c.LastName != "" {
		s = fmt.Sprint(s, "\nlast_name=", c.LastName)
	}

	if c.PhotoURL != "" {
		s = fmt.Sprint(s, "\nphoto_url=", c.PhotoURL)
	}

	if c.Username != "" {
		s = fmt.Sprint(s, "\nusername=", c.Username)
	}

	return s
}

func computeHmac256(msg []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return h.Sum(nil)
}
