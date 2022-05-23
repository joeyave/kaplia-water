package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"os"
)

type BotController struct {
}

func (c *BotController) Start(bot *gotgbot.Bot, ctx *ext.Context) error {

	markup := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "Зробити замовлення",
					WebApp: &gotgbot.WebAppInfo{
						Url: os.Getenv("HOST") + "/webapp/menu",
					},
				},
			},
		},
	}

	_, err := ctx.EffectiveChat.SendMessage(bot, "Почнемо 💧\n\nТисни на кнопку нижче, щоб зробити замовлення!", &gotgbot.SendMessageOpts{
		ReplyMarkup: markup,
	})
	if err != nil {
		return err
	}

	return nil
}
