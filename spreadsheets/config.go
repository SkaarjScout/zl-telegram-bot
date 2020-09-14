package spreadsheets

type Config struct {
	SpreadsheetId   string `yaml:"SpreadsheetId"`
	CredentialsJson string `yaml:"CredentialsJson"`
	RefreshToken    string `yaml:"RefreshToken"`
}
