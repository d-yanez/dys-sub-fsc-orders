package ports

import "context"

type TelegramMessage struct {
	Text           string
	PhotoURL       string
	ParseMode      string
	DisablePreview bool
}

type TelegramNotifier interface {
	Send(ctx context.Context, msg TelegramMessage) error
}
