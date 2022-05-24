package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/gin-gonic/gin"
	"github.com/joeyave/kaplia-water/controller"
	"github.com/joeyave/kaplia-water/repository"
	"github.com/joeyave/kaplia-water/state"
	"github.com/joeyave/kaplia-water/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"html/template"
	"net/http"
	"os"
	"time"
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

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating mongo client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer mongoClient.Disconnect(ctx)
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("error pinging mongo")
	}

	userRepository := &repository.UserRepository{
		MongoClient: mongoClient,
	}
	productRepository := &repository.ProductRepository{
		MongoClient: mongoClient,
	}

	botController := &controller.BotController{
		UserRepository: userRepository,
	}

	webAppController := &controller.WebAppController{
		Bot:               bot,
		ProductRepository: productRepository,
	}
	shopController := &controller.ShopController{
		Bot:               bot,
		UserRepository:    userRepository,
		ProductRepository: productRepository,
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

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, botController.RegisterUser), 0)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		if msg.ViaBot != nil && msg.ViaBot.Username == bot.Username {
			return false
		}
		return true
	}, botController.Start), 1)

	dispatcher.AddHandlerToGroup(handlers.NewCallback(util.CallbackState(state.ConfirmOrder_ChooseTime), botController.ConfirmOrder_ChooseTime), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(util.CallbackState(state.ConfirmOrder), botController.ConfirmOrder), 1)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, botController.UpdateUser), 2)

	router := gin.New()
	router.SetFuncMap(template.FuncMap{
		"hex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
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
