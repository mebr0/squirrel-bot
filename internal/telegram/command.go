package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case "game":
		return b.handleGameCommand(message)
	}

	return nil
}

func (b *Bot) handleGameCommand(message *tgbotapi.Message) error {
	player := squirrel.NewPlayer(message.Chat.ID, message.Chat.UserName, message.Chat.FirstName, message.Chat.LastName)

	b.game = squirrel.NewGameWithBots(player)

	b.game.StartFirstRound()

	ui := b.drawGame(0)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      message.Chat.ID,
			ReplyMarkup: b.inlineKeyboard(0),
		},
		Text:                  ui,
		ParseMode:             "MarkdownV2",
		DisableWebPagePreview: false,
	}

	m, err := b.bot.Send(msg)

	if err != nil {
		return err
	}

	b.game.UpdateChats([4]int{m.MessageID, -1, -1, -1})

	return nil
}
