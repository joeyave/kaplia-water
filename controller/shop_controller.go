package controller

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gin-gonic/gin"
	"strings"
)

type ShopController struct {
	Bot *gotgbot.Bot
}

type OrderedItem struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type MakeOrderData struct {
	OrderedItems []*OrderedItem `json:"order_data"`

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
	b.WriteString("Підсумок замовлення:\n\n")
	for _, item := range data.OrderedItems {
		fmt.Fprintf(b, " %s x%d\n", item.ID, item.Count)
	}

	fmt.Fprintf(b, "\nДата: %s\n", data.Date)
	fmt.Fprintf(b, "Номер: %s\n", data.Phone)
	fmt.Fprintf(b, "Адреса: %s\n", data.Address)
	if data.Comment != "" {
		fmt.Fprintf(b, "Коментар: %s\n", data.Comment)
	}

	_, err = c.Bot.AnswerWebAppQuery(data.Auth.QueryID, gotgbot.InlineQueryResultArticle{
		Id:    data.Auth.QueryID,
		Title: data.Auth.QueryID,
		InputMessageContent: gotgbot.InputTextMessageContent{
			MessageText: b.String(),
			ParseMode:   "HTML",
		},
	})
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	_, err = c.Bot.SendMessage(data.Auth.User.ID, "Ваше замовлення прийняте!\n\nЧекайте на повідомлення або дзвінок з підтвердженням.", nil)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.Status(200)
}
