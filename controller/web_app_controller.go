package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebAppController struct {
	Bot *gotgbot.Bot
}

type Item struct {
	ID          string
	Price       int
	PhotoURL    string
	Title       string
	Description string
}

func (h *WebAppController) Menu(ctx *gin.Context) {

	items := []*Item{
		{
			ID:          "qwe",
			Price:       69990,
			PhotoURL:    "./img/cafe/bottle.png",
			Title:       "Вода",
			Description: "Бутиль очищеної води 18,9 л.️",
		},
		{
			ID:          "ert",
			Price:       149990,
			PhotoURL:    "./img/cafe/pump.webp",
			Title:       "Помпа",
			Description: "Механічна помпа для бутилю.",
		},
		{
			ID:          "wer",
			Price:       799990,
			PhotoURL:    "./img/cafe/filter.png",
			Title:       "Фільтр",
			Description: "Фільтр з підігрівом води.",
		},
	}

	ctx.HTML(http.StatusOK, "index.go.html", gin.H{
		"items": items,
	})
}
