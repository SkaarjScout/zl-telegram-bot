package spreadsheets

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

type Client struct {
	service *sheets.Service
	config  Config
}

func NewClient(clientConfig Config) Client {

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(clientConfig.CredentialsJson), sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Panicf("Unable to parse app credentials to config: %v", err)
	}
	srv, err := sheets.NewService(context.Background(), option.WithTokenSource(
		config.TokenSource(context.Background(), &oauth2.Token{RefreshToken: clientConfig.RefreshToken})),
	)
	if err != nil {
		log.Panicf("Unable to create spreadsheets service: %v", err)
	}
	return Client{srv, clientConfig}

}

func (client *Client) FindRow(nickname string) ([]interface{}, error) {
	resp, err := client.service.Spreadsheets.Values.Get(client.config.SpreadsheetId, "A:A").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get column from sheet: %w", err)
	}
	for index, row := range resp.Values {
		if row[0] == nickname {
			resp, err := client.service.Spreadsheets.Values.Get(
				client.config.SpreadsheetId, fmt.Sprintf("%d:%d", index+1, index+1)).Do()
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve row from sheet: %w", err)
			}
			return resp.Values[0], nil
		}
	}
	return nil, nil
}
