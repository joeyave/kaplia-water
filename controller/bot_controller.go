package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joeyave/kaplia-water/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"strings"
)

type BotController struct {
	UserRepository *repository.UserRepository
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

func (c *BotController) RegisterUser(bot *gotgbot.Bot, ctx *ext.Context) error {

	user, err := c.UserRepository.FindOneByID(context.Background(), ctx.EffectiveUser.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			user, err = c.UserRepository.UpdateOne(context.Background(), &repository.User{
				ID: ctx.EffectiveUser.Id,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	ctx.Data["user"] = user
	user = ctx.Data["user"].(*repository.User)

	user.Name = strings.TrimSpace(fmt.Sprintf("%s %s", ctx.EffectiveUser.FirstName, ctx.EffectiveUser.LastName))

	return nil
}

func (c *BotController) UpdateUser(bot *gotgbot.Bot, ctx *ext.Context) error {

	user := ctx.Data["user"].(*repository.User)

	_, err := c.UserRepository.UpdateOne(context.Background(), user)
	return err
}
