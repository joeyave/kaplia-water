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
			ID:          "1",
			Price:       69990,
			PhotoURL:    "./img/cafe/bottle.png",
			Title:       "Вода",
			Description: "Капля Water™️",
		},
	}
	ctx.HTML(http.StatusOK, "index.go.html", gin.H{
		"items": items,
	})
}
