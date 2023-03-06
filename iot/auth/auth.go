package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/xkamail/huberlink-platform/iot/account"
	"github.com/xkamail/huberlink-platform/pkg/api"
	"github.com/xkamail/huberlink-platform/pkg/config"
	"github.com/xkamail/huberlink-platform/pkg/discord"
	"github.com/xkamail/huberlink-platform/pkg/passhash"
	"github.com/xkamail/huberlink-platform/pkg/pgctx"
	"github.com/xkamail/huberlink-platform/pkg/rand"
	"github.com/xkamail/huberlink-platform/pkg/snowid"
	"github.com/xkamail/huberlink-platform/pkg/uierr"
)

// 7 Days
const _refreshTokenLifetime = 24 * time.Hour * 7

type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type SignInParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *SignInParam) Valid() error {
	p.Username = strings.TrimSpace(p.Username)
	if p.Username == "" {
		return uierr.Invalid("username", "username is required")
	}
	if p.Password == "" {
		return uierr.Invalid("password", "password is required")
	}
	if len(p.Password) < 6 {
		return uierr.Invalid("password", "password must be at least 6 characters")
	}
	if len(p.Username) < 3 {
		return uierr.Invalid("username", "username must be at least 3 characters")
	}
	return nil
}

func SignIn(ctx context.Context, p *SignInParam) (*TokenResponse, error) {
	acc, err := account.FindByUsername(ctx, p.Username)
	if errors.Is(err, account.ErrNotFound) {
		return nil, ErrUsernameAndPasswordIncorrect
	}
	if err != nil {
		return nil, err
	}
	if !passhash.CheckPasswordHash(p.Password, acc.Password) {
		return nil, ErrUsernameAndPasswordIncorrect
	}

	tx, err := pgctx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	cfg := config.Load()
	userID := acc.ID.Int()
	jwtToken, err := jwtGenerate(userID, time.Hour*3, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := createRefreshToken(ctx, tx, userID)
	if err != nil {

		return nil, err
	}

	// update latest discord profile
	_, err = tx.Exec(ctx, `update users set updated_at = $1 where id = $2`,
		time.Now(),
		userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &TokenResponse{
		jwtToken,
		refreshToken,
	}, nil
}

type SignInWithDiscordParam struct {
	Code string `json:"code" validate:"required"`
}

func (p *SignInWithDiscordParam) Valid() error {
	p.Code = strings.TrimSpace(p.Code)
	if p.Code == "" {
		return uierr.Invalid("code", "code is required")
	}
	return nil
}

func SignInWithDiscord(ctx context.Context, discord discord.Client, p *SignInWithDiscordParam) (*TokenResponse, error) {
	accessToken, err := discord.GetAccessToken(ctx, p.Code)
	if err != nil {
		return nil, err
	}
	profile, err := discord.GetProfile(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	var (
		userID int64
	)
	err = pgctx.QueryRow(ctx, `
			select id from users where (discord_id = $1 or email = $2)`,
		profile.ID,
		profile.Email,
	).Scan(
		&userID,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	tx, err := pgctx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// create a new account
	if userID == 0 {
		newID, err := account.Create(ctx, &account.User{
			Username:  profile.Username,
			Email:     profile.Email,
			Password:  "",
			DiscordID: profile.ID,
		})
		if err != nil {
			return nil, err
		}

		userID = newID.Int()

	}
	cfg := config.Load()
	jwtToken, err := jwtGenerate(userID, time.Hour*3, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := createRefreshToken(ctx, tx, userID)
	if err != nil {

		return nil, err
	}

	// update latest discord profile
	_, err = tx.Exec(ctx, `
			update users 
			set updated_at = $1,
			    email = $3,
			    avatar_url = $4
			where id = $2`,
		time.Now(),
		userID,
		profile.Email,
		profile.AvatarURL(),
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &TokenResponse{
		jwtToken,
		refreshToken,
	}, nil
}

func InvokeRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	var (
		userID    int64
		expiredAt time.Time
	)
	err := pgctx.QueryRow(ctx, `select user_id, expired_at from users_refresh_tokens where token = $1`,
		refreshToken,
	).Scan(
		&userID,
		&expiredAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	// do generate jwt
	if expiredAt.Before(time.Now()) {
		return nil, ErrRefreshTokenExpired
	}
	cfg := config.Load()
	jwtToken, err := jwtGenerate(userID, time.Hour*3, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	tx, err := pgctx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `delete from users_refresh_tokens where token = $1`, refreshToken)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := createRefreshToken(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `update users set updated_at = $1 where id = $2`, time.Now(), userID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &TokenResponse{
		jwtToken,
		newRefreshToken,
	}, nil
}

func createRefreshToken(ctx context.Context, tx pgx.Tx, userID int64) (string, error) {
	refreshToken, err := rand.String(300)
	if err != nil {
		return "", err
	}
	now := time.Now()
	_, err = tx.Exec(ctx, `
		insert into users_refresh_tokens 
		    (id,user_id,token,expired_at,issued_at,created_at) 
		values ($1,$2,$3,$4,$5,$6)`,
		snowid.Gen(),
		userID,
		refreshToken,
		now.Add(_refreshTokenLifetime),
		now,
		now,
	)
	return refreshToken, err
}

func SignInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		if len(raw) == 0 {
			api.WriteError(w, ErrNoJWTToken)
			return
		}
		if len(raw) < len("Bearer ") {
			api.WriteError(w, ErrInvalidJWTSchema)
			return
		}
		var jwtToken string
		if len(raw) > len("Bearer ") {
			jwtToken = raw[len("Bearer "):]
		}
		claims, err := jwtVerify(jwtToken)
		if err != nil {
			api.WriteError(w, err)
			return
		}
		acc, err := account.Find(r.Context(), claims.UserID)
		if err != nil {
			api.WriteError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(account.NewContext(r.Context(), acc)))
	})
}
