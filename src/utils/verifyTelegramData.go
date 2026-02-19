package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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

func (data *TelegramData) VerifyTelegramData() error {

	check := []string{
		"auth_date=" + strconv.FormatInt(data.AuthDate, 10),
		"first_name=" + data.FirstName,
		"id=" + strconv.FormatInt(data.ID, 10),
	}

	if data.LastName != "" {
		check = append(check, "last_name="+data.LastName)
	}
	if data.Username != "" {
		check = append(check, "username="+data.Username)
	}
	if data.PhotoURL != "" {
		check = append(check, "photo_url="+data.PhotoURL)
	}

	sort.Strings(check)

	fmt.Println(check)
	fmt.Println(os.Getenv("TG_TOKEN"))

	secret := sha256.Sum256([]byte(os.Getenv("TG_TOKEN")))

	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(strings.Join(check, "\n")))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	fmt.Println(calculatedHash, data.Hash)
	if !hmac.Equal([]byte(calculatedHash), []byte(data.Hash)) {
		return errors.New("Hash is invalid")
	}
	if time.Now().Unix()-data.AuthDate > 86400 {
		return errors.New("data is outdated")
	}

	return nil
}
