package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"musaitHrMgBotGo/pkg/repository"
	"musaitHrMgBotGo/pkg/services"
)

const (
	commandStart    = "start"
	commandAuth     = "auth"
	commandGet      = "get"
	commandHelp     = "help"
	replyHelpMsg    = "Bot haqida quyidagi malumotlarni keltirib o'tamiz.\n Bot quyidagi kamandalardan iborat.\n /start\n /auth\n /get \n /help\n Noqulayliklar uchun developerga murojat qiling @mnasilloyev."
	replyWelcomeMsg = "Assalomu alekum botimizga xush kelibsiz. Bu bot orqali siz google spreadsheetdagi malumotlaringizni osonlik bilan olishingiz mumkin. Buning uchun avval avtorizatsiyadan o'tishingiz kerak bo'ladi va quyidagi kamandani bosing /auth"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandAuth:
		return b.handleAuthCommand(message)
	case commandGet:
		return b.handleGetCommand(message)
	case commandHelp:
		return b.handleHelpCommand(message)
	default:
		return b.handleUnknownCommand(message)

	}
}
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Siz menga kamanda bering")
	_, err := b.bot.Send(msg)
	return err
}

// Functions of Commands
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.Text = replyWelcomeMsg
	_, err := b.bot.Send(msg)
	return err
}
func (b *Bot) handleAuthCommand(message *tgbotapi.Message) error {
	userID := message.Chat.ID

	text := "Iltimos menga 'API_KEY'ingizni yuboring"
	msg := tgbotapi.NewMessage(userID, text)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	userState[userID] = "awaiting_api_key"
	return nil
}

func (b *Bot) handleGetCommand(message *tgbotapi.Message) error {
	userID := message.Chat.ID
	spreadsheetId, err := b.userRepository.Get(userID, repository.SpreadSheetId)
	if err != nil {
		fmt.Printf("error get spreadsheet id: %v\n", err)
	}
	apiKey, err := b.userRepository.Get(userID, repository.ApiKey)
	if err != nil {
		fmt.Printf("error get api key: %v\n", err)
	}
	url := services.NewURL(spreadsheetId, "Sheet1", "A1", "K50", apiKey)
	data, err := services.GetSheetData(url)
	if err != nil {
		return err
	}
	markup, err := b.createInlineKeyboard(data)
	if err != nil {
		fmt.Printf("error create inline keyboard: %v\n", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Xodimlar ro'yxati:")
	msg.ReplyMarkup = markup
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyHelpMsg)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return err
}
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, "Siz mavjud bo'lmagan kamandani berdingiz")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleAuthProcess(update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userState[userID]

	switch state {
	case "awaiting_api_key":
		apiKey := update.Message.Text
		msg := tgbotapi.NewMessage(userID, "API kalit saqlandi. Endi Spreadsheet ID ni kiriting")
		_, err := b.bot.Send(msg)
		if err != nil {
			fmt.Printf("Api_key Message yuborilmadi %w", err)
		}
		if err = b.userRepository.Save(userID, apiKey, repository.ApiKey); err != nil {
			fmt.Printf("Api_key Saqlanmadi %w", err)
		}
		userState[userID] = "awaiting_spreadsheet_id"
	case "awaiting_spreadsheet_id":
		spreadsheetID := update.Message.Text
		msg := tgbotapi.NewMessage(userID, "Spreadsheet ID qabul qilindi va saqlandi!")
		_, err := b.bot.Send(msg)
		if err != nil {
			fmt.Printf("spreadsheetID Message yuborilmadi %w", err)
		}
		if err = b.userRepository.Save(userID, spreadsheetID, repository.SpreadSheetId); err != nil {
			fmt.Printf("spreadsheetID Saqlanmadi %w", err)
		}
		delete(userState, userID)

	}
}

func (b *Bot) createInlineKeyboard(input []map[string]string) (tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, item := range input {
		name := item["Ism familiya"]
		id := item["#"]
		btn := tgbotapi.NewInlineKeyboardButtonData(name, id)
		row = append(row, btn)

		if len(row) == 3 {
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		keyboard = append(keyboard, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...), nil
}

func (b *Bot) handleCallbackQuery(update tgbotapi.Update) error {
	userID := update.CallbackQuery.Message.Chat.ID
	spreadsheetId, err := b.userRepository.Get(userID, repository.SpreadSheetId)
	if err != nil {
		fmt.Printf("error get spreadsheet id: %v\n", err)
		return err
	}
	apiKey, err := b.userRepository.Get(userID, repository.ApiKey)
	if err != nil {
		fmt.Printf("error get api key: %v\n", err)
		return err
	}
	url := services.NewURL(spreadsheetId, "Sheet1", "A1", "K50", apiKey)
	data, err := services.GetSheetData(url)
	if err != nil {
		return err
	}
	callbackQuery := update.CallbackQuery
	text, err := b.handleCallBackQueryMessage(callbackQuery, data)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
func (b *Bot) handleCallBackQueryMessage(query *tgbotapi.CallbackQuery, data []map[string]string) (string, error) {
	for _, item := range data {
		if query.Data == item["#"] {
			id := item["#"]
			name := item["Ism familiya"]
			fixedSalary := item["Fixed salary"]
			workingPercentage := item["Ishlagan %"]
			fixByTime := item["Vaqtbay Fixed"]
			avans := item["Avans"]
			debt := item["Qarz"]
			lunch := item["Tushlik"]
			KPI := item["KPI"]
			final := item["Final"]
			msg := fmt.Sprintf("ID: %s\nIsm familiya: %s\nFixed salary: %s\nIshlagan: %s\nVaqtbay Fixed: %s\nAvans: %s\nQarz: %s\nTushlik: %s\nKPI: %s\nFinal: %s\n", id, name, fixedSalary, workingPercentage, fixByTime, avans, debt, lunch, KPI, final)
			return msg, nil
		}
	}
	return "", nil
}
