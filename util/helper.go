package util

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func SplitKeyboardToColumns(k [][]gotgbot.KeyboardButton, colNum int) [][]gotgbot.KeyboardButton {

	var newK [][]gotgbot.KeyboardButton
	var i int

	for _, row := range k {
		for _, button := range row {
			if i == colNum {
				i = 0
			}

			if i == 0 {
				newK = append(newK, []gotgbot.KeyboardButton{button})
			} else if i < colNum {
				newK[len(newK)-1] = append(newK[len(newK)-1], button)
			}
			i++
		}
	}

	return newK
}

func SplitInlineKeyboardToColumns(k [][]gotgbot.InlineKeyboardButton, colNum int) [][]gotgbot.InlineKeyboardButton {

	var newK [][]gotgbot.InlineKeyboardButton
	var i int

	for _, row := range k {
		for _, button := range row {
			if i == colNum {
				i = 0
			}

			if i == 0 {
				newK = append(newK, []gotgbot.InlineKeyboardButton{button})
			} else if i < colNum {
				newK[len(newK)-1] = append(newK[len(newK)-1], button)
			}
			i++
		}
	}

	return newK
}
