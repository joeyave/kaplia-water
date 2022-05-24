package controller

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gin-gonic/gin"
	"github.com/joeyave/kaplia-water/repository"
	"github.com/joeyave/kaplia-water/state"
	"github.com/joeyave/kaplia-water/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type ShopController struct {
	Bot               *gotgbot.Bot
	UserRepository    *repository.UserRepository
	ProductRepository *repository.ProductRepository
}

type OrderedProduct struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type MakeOrderData struct {
	OrderedProducts []*OrderedProduct `json:"order_data"`

	Date    string `json:"date"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Comment string `json:"comment"`

	Auth struct {
		QueryID string `json:"query_id"`
		User    struct {
			ID           int64  `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
		} `json:"user"`
		AuthDate string `json:"auth_date"`
		Hash     string `json:"hash"`
	} `json:"_auth"`
}

func (c *ShopController) MakeOrder(ctx *gin.Context) {

	var data *MakeOrderData
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	b := &strings.Builder{}

	totalPrice := float32(0)
	for _, orderedProduct := range data.OrderedProducts {
		ID, err := primitive.ObjectIDFromHex(orderedProduct.ID)
		if err != nil {
			return
		}
		product, err := c.ProductRepository.FindOneByID(context.Background(), ID)
		if err != nil {
			return
		}

		price := float32(product.Price*orderedProduct.Count) / 1000
		totalPrice += price

		fmt.Fprintf(b, "%s x%d - <b>₴%.2f</b>\n", product.Title, orderedProduct.Count, price)
	}

	fmt.Fprintf(b, "\nВсього - <b>₴%.2f</b>\n", totalPrice)

	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	fmt.Fprintf(b, "\nДата: %s\n", date.Format("02.01.2006"))
	fmt.Fprintf(b, "Номер: +%s\n", data.Phone)
	fmt.Fprintf(b, "Адреса: %s\n", data.Address)
	if data.Comment != "" {
		fmt.Fprintf(b, "Коментар: %s\n", data.Comment)
	}

	summaryText := fmt.Sprintf("Підсумок замовлення:\n\n%s", b.String())
	_, err = c.Bot.AnswerWebAppQuery(data.Auth.QueryID, gotgbot.InlineQueryResultArticle{
		Id:    data.Auth.QueryID,
		Title: data.Auth.QueryID,
		InputMessageContent: gotgbot.InputTextMessageContent{
			MessageText: summaryText,
			ParseMode:   "HTML",
		},
	})
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	msg, err := c.Bot.SendMessage(data.Auth.User.ID, "⏰ Ваше замовлення прийняте!\n\nОчікуйте на повідомлення або дзвінок з підтвердженням.", nil)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	textForAdmins := fmt.Sprintf("Вітаю! Ви маєте нове замовлення.\n\n%s", b.String())
	markupForAdmins := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Відхилити 🚫", CallbackData: util.CallbackData(state.DeclineOrder, fmt.Sprintf("%d:%d", data.Auth.User.ID, msg.MessageId))},
				{Text: "Підтвердити ✅", CallbackData: util.CallbackData(state.ConfirmOrder_ChooseTime, fmt.Sprintf("%d:%d", data.Auth.User.ID, msg.MessageId))},
			},
		},
	}

	admins, err := c.UserRepository.FindManyByRole(context.Background(), repository.AdminRole)
	if err != nil {
		return
	}
	for _, admin := range admins {
		_, err := c.Bot.SendMessage(admin.ID, textForAdmins, &gotgbot.SendMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: markupForAdmins,
		})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		time.Sleep(1 * time.Second)
	}

	ctx.Status(200)
}
