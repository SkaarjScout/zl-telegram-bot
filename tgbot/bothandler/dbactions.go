package bothandler

import (
	"context"
	"database/sql"
	"fmt"
)

func (bot *Bot) checkUserTable(ctx context.Context, conn *sql.Conn) error {
	_, err := conn.ExecContext(ctx, fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
					id int PRIMARY KEY NOT NULL
		)`, bot.botConfig.UserTableName))
	return err
}

func (bot *Bot) addUser(ctx context.Context, userId int) error {
	conn, err := bot.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("error on connection get: %w", err)
	}
	defer conn.Close()
	if err = bot.checkUserTable(ctx, conn); err != nil {
		return fmt.Errorf("error on table check: %w", err)
	}
	_, err = conn.ExecContext(ctx, fmt.Sprintf(
		`INSERT INTO %s (id) values ($1)
		ON CONFLICT (id) DO NOTHING`, bot.botConfig.UserTableName), userId)
	if err != nil {
		return fmt.Errorf("error on user add to database: %w", err)
	}
	return nil
}
