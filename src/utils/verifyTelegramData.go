package utils

import (
	"os"
	"time"

	telegramloginwidget "github.com/LipsarHQ/go-telegram-login-widget"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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
	auth := telegramloginwidget.AuthorizationData{
		AuthDate:  data.AuthDate,
		FirstName: data.FirstName,
		Hash:      data.Hash,
		ID:        data.ID,
		LastName:  data.LastName,
		PhotoURL:  data.PhotoURL,
		Username:  data.Username,
	}

	return auth.Check(os.Getenv("TG_TOKEN"))
}

type TelegramInitData struct {
	Data string `json:"data"`
}

func (data *TelegramInitData) VerifyTelegramInitData() error {
	return initdata.Validate(data.Data, os.Getenv("TG_TOKEN"), 24*time.Hour)
}

func (data *TelegramInitData) ParseTelegramInitData() (initdata.InitData, error) {
	return initdata.Parse(data.Data)
}
