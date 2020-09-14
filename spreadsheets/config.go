package spreadsheets

type Config struct {
	SpreadsheetId     string `yaml:"SpreadsheetId"`
	CredentialsJson   string `yaml:"CredentialsJson"`
	OauthRefreshToken string `yaml:"OauthRefreshToken"`
}
