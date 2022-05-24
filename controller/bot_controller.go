package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joeyave/kaplia-water/repository"
	"github.com/joeyave/kaplia-water/state"
	"github.com/joeyave/kaplia-water/util"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"strconv"
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

func (c *BotController) ConfirmOrder_ChooseTime(bot *gotgbot.Bot, ctx *ext.Context) error {

	payload := util.ParseCallbackPayload(ctx.CallbackQuery.Data)

	markup := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Перша половина дня ☀️", CallbackData: util.CallbackData(state.ConfirmOrder, payload+":am")},
			},
			{
				{Text: "Друга половина дня 🌙", CallbackData: util.CallbackData(state.ConfirmOrder, payload+":pm")},
			},
			{
				{Text: "Зв'яжемося додатково 📞", CallbackData: util.CallbackData(state.ConfirmOrder, payload+":phone")},
			},
			//{
			//	{Text: "Назад ↩️", CallbackData: util.CallbackData(state.ConfirmOrder, payload+":phone")},
			//},
		},
	}

	text := ctx.CallbackQuery.Message.Text + "\n\nОбери час, коли замовлення буде доставлено:"

	_, _, err := ctx.EffectiveMessage.EditText(bot, text, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: markup,
		Entities:    ctx.CallbackQuery.Message.Entities,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *BotController) ConfirmOrder(bot *gotgbot.Bot, ctx *ext.Context) error {

	payload := util.ParseCallbackPayload(ctx.CallbackQuery.Data)
	split := strings.Split(payload, ":")

	userID, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return err
	}
	messageID, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return err
	}

	textForClient := "✅ Вітаю! Ваше замовлення підтверджене."
	switch split[2] {
	case "am":
		textForClient += "\n\nОчікуйте на доставку у <b>першій половині</b> дня."
	case "pm":
		textForClient += "\n\nОчікуйте на доставку у <b>другій половині</b> дня."
	default:
		textForClient += "\n\nЗ вами <b>зв'яжуться додатково</b> щодо часу доставки."
	}

	_, err = bot.SendMessage(userID, textForClient, &gotgbot.SendMessageOpts{
		ReplyToMessageId: messageID,
		ParseMode:        "HTML",
	})
	if err != nil {
		return err
	}

	textForAdmin := strings.ReplaceAll(ctx.EffectiveMessage.OriginalHTML(),
		"Обери час, коли замовлення буде доставлено:",
		"✅ Замовлення підтверджене. Клієнт отримав повідомлення.")
	switch split[2] {
	case "am":
		textForAdmin += "\n\nДоставка у <b>першій половині</b> дня."
	case "pm":
		textForAdmin += "\n\nДоставка у <b>другій половині</b> дня."
	default:
		textForAdmin += "\n\nТреба <b>зв'язатися додатково</b> щодо часу доставки."
	}

	_, _, err = ctx.EffectiveMessage.EditText(bot, textForAdmin, &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return err
	}

	return nil
}
