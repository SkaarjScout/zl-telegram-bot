package tgbot

type Config struct {
	TelegramBotToken string `yaml:"TelegramBotToken"`
	DebugEnabled     bool   `yaml:"DebugEnabled"`
	UpdateTimeout    int    `yaml:"UpdateTimeout"`
}
