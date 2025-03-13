package telegram

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/env"
	"ambic/internal/infra/file"
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramIf interface {
	SendPartnerRegistrationMessage(data dto.PartnerRegistrationTelegramMessage) error
}

type Telegram struct {
	bot      *bot.Bot
	template string
	chatId   int64
}

func NewTelegram(env *env.Env) TelegramIf {
	bot, err := bot.New(env.TelegramBotToken)

	if err != nil {
		log.Fatalf("Gagal membuat bot: %v", err)
	}

	template := "internal/infra/templates/message"
	chatId := env.TelegramChatID

	return &Telegram{bot, template, chatId}
}

func (t *Telegram) SendPartnerRegistrationMessage(data dto.PartnerRegistrationTelegramMessage) error {
	body, err := file.ReadHTML(t.template, "partner_registration")
	if err != nil {
		return err
	}

	message := fmt.Sprintf(body, data.UserID, data.UserName, data.UserUsername, data.UserEmail, data.UserPhone, data.UserAddress, data.UserGender, data.UserRegisteredAt,
		data.PartnerID, data.BusinessType, data.BusinessName, data.BusinessAddress, data.BusinessCity, data.BusinessGmaps, data.BusinessInstagram, data.BusinessInstagram, data.BusinessPhoto, data.BusinessRegisteredAt)

	_, err = t.bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID:    t.chatId,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})

	if err != nil {
		return fmt.Errorf("gagal mengirim pesan: %w", err)
	}

	return nil
}
