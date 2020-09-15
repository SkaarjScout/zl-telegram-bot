package bothandler

type Config struct {
	TelegramBotToken string `yaml:"TelegramBotToken"`
	DebugEnabled     bool   `yaml:"DebugEnabled"`
	UpdateTimeout    int    `yaml:"UpdateTimeout"`
	WorkerCount      int    `yaml:"WorkerCount"`
	UserTableName    string `yaml:"UserTableName"`
}
