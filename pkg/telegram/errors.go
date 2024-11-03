package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errUnautorized               = errors.New("user is not autorized")
	errInvalidURL                = errors.New("url is invalid")
	errUnableToSaveApiKey        = errors.New("unable to save api_key")
	errUnableToSaveSpreadsheetID = errors.New("unable to save spreadsheet_id")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "Aniqlanmagan xatolik sodir bo'ldi.")
	switch {
	case errors.Is(err, errUnautorized):
		msg.Text = "Siz hali avtorizatsiyadan o'tmagansiz. /auth kamandasi orqali avtorizatsiyadan o'tib oling"
		b.bot.Send(msg)
	case errors.Is(err, errInvalidURL):
		msg.Text = "Mavjud bo'lmagan URL"
		b.bot.Send(msg)
	case errors.Is(err, errUnableToSaveApiKey):
		msg.Text = "Api Keyni saqlashni iloji bo'lmadi"
		b.bot.Send(msg)
	case errors.Is(err, errUnableToSaveSpreadsheetID):
		msg.Text = "SpreadsheetIDni saqlashni iloji bo'lmadi"
		b.bot.Send(msg)
	}
}
