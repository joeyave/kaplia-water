package util

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"strconv"
	"strings"
)

func IetfToIsoLangCode(languageCode string) string {
	switch languageCode {
	default:
		return "ru_RU"
	}
}

func CallbackData(state int, payload string) string {
	callbackData := strconv.Itoa(state) + ":" + payload
	if len(callbackData) > 64 {
		panic(fmt.Sprintf("size of callback_data is bigger than 64: size=%d, callback_data=%s", len(callbackData), callbackData))
	}
	return callbackData
}

func CallbackState(state int) filters.CallbackQuery {
	return func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasPrefix(cq.Data, strconv.Itoa(state)+":")
	}
}

func ParseCallbackPayload(data string) string {
	parsedData := strings.Split(data, ":")
	return strings.Join(parsedData[1:], ":")
}

const CallbackCacheURL = "https://t.me/callbackCache"
