package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebAppController struct {
	Bot *gotgbot.Bot
}

func (h *WebAppController) Menu(ctx *gin.Context) {

	//fmt.Println(ctx.Request.URL.String())
	//
	//hex := ctx.Query("bandId")
	//bandID, err := primitive.ObjectIDFromHex(hex)
	//if err != nil {
	//	return
	//}
	//
	//band, err := h.BandService.FindOneByID(bandID)
	//if err != nil {
	//	return
	//}
	//
	//event := &entity.Event{
	//	Time:   time.Now(),
	//	BandID: bandID,
	//	Band:   band,
	//}
	//eventNames, err := h.EventService.GetMostFrequentEventNames(bandID, 4)
	//if err != nil {
	//	return
	//}
	//
	ctx.HTML(http.StatusOK, "index.go.html", gin.H{})
}
