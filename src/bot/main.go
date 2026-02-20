package bot

import (
	"context"
	"fmt"
	"os"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func InitializeTelegramBot() {
	bot, err := tg.NewBot(os.Getenv("TG_TOKEN"), tg.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	for update := range updates {
		if update.Message != nil {
			chat := update.Message.Chat.ChatID()
			if update.Message.Text == "/start" {
				keyboard := tu.InlineKeyboard(
					tu.InlineKeyboardRow(
						tu.InlineKeyboardButton("В бой!").WithWebApp(&tg.WebAppInfo{
							URL: os.Getenv("NUXT_PUBLIC_BASE_URL"),
						}),
					),
				)
				response, err := bot.SendMessage(
					context.Background(),
					tu.Message(
						chat,
						"[\n](https://pixelbattle.xlsft.ru/og_image.png)*xlsft`s pixelbattle*. Один холст и бесконечная битва за каждый пиксель. Присоединяйся!",
					).WithReplyMarkup(keyboard).WithParseMode("Markdown").WithLinkPreviewOptions(&tg.LinkPreviewOptions{
						PreferLargeMedia: true,
						ShowAboveText:    true,
					}),
				)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(response)
			}
		}

	}
}
