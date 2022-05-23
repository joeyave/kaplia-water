package controller

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gin-gonic/gin"
	"github.com/joeyave/kaplia-water/repository"
	"net/http"
)

type WebAppController struct {
	Bot               *gotgbot.Bot
	ProductRepository *repository.ProductRepository
}

func (h *WebAppController) Menu(ctx *gin.Context) {

	products, err := h.ProductRepository.FindAll(context.Background())
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.HTML(http.StatusOK, "index.go.html", gin.H{
		"items": products,
	})
}
