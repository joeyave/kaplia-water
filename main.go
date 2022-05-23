package main

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joeyave/kaplia-water/controller"
	"html/template"
	"net/http"
	"os"
)

func main() {

	bot, err := gotgbot.NewBot(os.Getenv("BOT_TOKEN"), &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	bot.SetChatMenuButton(&gotgbot.SetChatMenuButtonOpts{
		MenuButton: gotgbot.MenuButtonWebApp{
			Text: "Меню",
			WebApp: gotgbot.WebAppInfo{
				Url: os.Getenv("HOST") + "/webapp/menu",
			},
		},
	})

	botController := controller.BotController{}

	webAppController := controller.WebAppController{
		Bot: bot,
	}
	shopController := controller.ShopController{
		Bot: bot,
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			Error:       nil,
			Panic:       nil,
			ErrorLog:    nil,
			MaxRoutines: 0,
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		if msg.ViaBot != nil && msg.ViaBot.Username == bot.Username {
			return false
		}
		return true
	}, botController.Start))

	router := gin.New()
	router.SetFuncMap(template.FuncMap{
		"json": func(s interface{}) string {
			jsonBytes, err := json.Marshal(s)
			if err != nil {
				return ""
			}
			return string(jsonBytes)
		},
		"price_to_str": func(price int) string {
			return fmt.Sprintf("%.2f", float32(price)/1000)
		},
	})

	router.LoadHTMLGlob("webapp/*.go.html")
	router.Static("/webapp/css", "./webapp/css")
	router.Static("/webapp/img", "./webapp/img")
	router.Static("/webapp/js", "./webapp/js")

	router.GET("/webapp/menu", webAppController.Menu)

	router.POST("/shop/api/makeOrder", shopController.MakeOrder)

	go func() {
		// Start receiving updates.
		err = updater.StartPolling(bot, &ext.PollingOpts{DropPendingUpdates: true})
		if err != nil {
			panic("failed to start polling: " + err.Error())
		}
		fmt.Printf("%s has been started...\n", bot.User.Username)

		// Idle, to keep updates coming in, and avoid bot stopping.
		updater.Idle()
	}()

	err = router.Run()
	if err != nil {
		panic("error starting Gin: " + err.Error())
	}
}
