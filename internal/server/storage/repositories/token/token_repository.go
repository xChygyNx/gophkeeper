package token

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/model"
	customErrors "github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

const lengthToken = 32

type TokenRepository interface {
	Create(user *model.User) (string, error)
}

type Token struct {
	db *database.DB
}

func New(db *database.DB) *Token {
	return &Token{
		db: db,
	}
}

func (t *Token) Create(userID int64, lifetime time.Duration) (*model.Token, error) {
	token := &model.Token{}
	accessToken := encryption.GenerateAccessToken(lengthToken)
	currentTime := time.Now()

	if err := t.db.Pool.QueryRow(
		"INSERT INTO access_token (access_token, user_id, created_at, end_date_at) VALUES ($1, $2, $3, $4) "+
			"RETURNING access_token, user_id, created_at, end_date_at",
		accessToken,
		userID,
		currentTime,
		currentTime.Add(lifetime),
	).Scan(&token.AccessToken, &token.UserID, &token.CreatedAt, &token.EndDateAt); err != nil {
		return nil, fmt.Errorf("error in create token record in DB: %w", err)
	}
	return token, nil
}

func (t *Token) Validate(endDate time.Time) bool {
	valid := time.Now().After(endDate)
	return valid
}

func (t *Token) GetEndDateToken(accessToken string) (time.Time, error) {
	var end time.Time
	if err := t.db.Pool.QueryRow("SELECT end_date_at FROM access_token where access_token = $1",
		accessToken,
	).Scan(&end); err != nil {
		return end, fmt.Errorf("error scan end date of token from DB: %w", err)
	}
	return end, nil
}

func (t *Token) Block(accessToken string) (string, error) {
	var token string
	if err := t.db.Pool.QueryRow("UPDATE access_token SET end_date_at = $1 "+
		"where access_token = $2 RETURNING access_token",
		time.Now(),
		accessToken,
	).Scan(&token); err != nil {
		return token, fmt.Errorf("error in block access token: %w", err)
	}
	return token, nil
}

func (t *Token) BlockAllTokenUser(userID int64) (string, error) {
	var token string
	if err := t.db.Pool.QueryRow("UPDATE access_token SET end_date_at = $1 "+
		"where user_id = $2 RETURNING access_token",
		time.Now(),
		userID,
	).Scan(&token); err != nil {
		return token, fmt.Errorf("error in block all access tokens for UserID %d: %w", userID, err)
	}
	return token, nil
}

func (t *Token) GetList(userID int64) ([]model.Token, error) {
	tokens := make([]model.Token, 0)
	rows, err := t.db.Pool.Query("SELECT access_token, user_id, created_at, end_date_at FROM access_token "+
		"where user_id = $1",
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tokens, customErrors.ErrRecordNotFound
		} else {
			return tokens, fmt.Errorf("error in get list access tokens for UserID %d from DB: %w", userID, err)
		}
	}
	defer func() {
		err = rows.Close()
	}()
	for rows.Next() {
		token := model.Token{}
		err = rows.Scan(&token.AccessToken, &token.UserID, &token.CreatedAt, &token.EndDateAt)
		if err != nil {
			return tokens, fmt.Errorf("error in scan data about access token"+
				" for UserID %d from DB: %w", userID, err)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
