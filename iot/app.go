package iot

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/xkamail/huberlink-platform/iot/account"
	"github.com/xkamail/huberlink-platform/iot/auth"
	"github.com/xkamail/huberlink-platform/pkg/api"
	"github.com/xkamail/huberlink-platform/pkg/config"
	"github.com/xkamail/huberlink-platform/pkg/discord"
	"github.com/xkamail/huberlink-platform/pkg/uierr"
)

func Handlers() http.Handler {
	cfg := config.Load()

	discordClient := discord.NewClient(cfg.DiscordClientID, cfg.DiscordClientSecret, cfg.DiscordRedirectURI)

	router := chi.NewRouter()
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		api.WriteError(w, uierr.NotFound("api entry not found"))
	})
	// auth
	{
		router.Post("/auth/sign-in", h(func(ctx context.Context, r *http.Request) (any, error) {

			code := r.URL.Query().Get("code")
			return auth.SignInWithDiscord(ctx, discordClient, code)
		}))
		router.Post("/auth/refresh-token", h(func(ctx context.Context, r *http.Request) (any, error) {

			code := r.URL.Query().Get("refreshToken")
			return auth.InvokeRefreshToken(ctx, code)
		}))
		router.With(auth.SignInMiddleware).Get("/auth/me", h(func(ctx context.Context, r *http.Request) (any, error) {

			return account.FromContext(ctx)
		}))
	}
	// user
	{
		router.Get("/user/me", nil)
	}
	// home
	{
		// list my home
		router.Get("/home", nil)
		// create home

		// join home
	}
	// device
	{
		router.Get("/devices/all", nil)
		router.Post("/devices", nil)
		router.Get("/devices/{id}", nil)
		router.Delete("/devices/{id}", nil)
		router.Patch("/devices/{id}", nil)
	}
	return router
}

func h[T any](fn func(ctx context.Context, r *http.Request) (T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := fn(r.Context(), r)
		if err != nil {
			api.WriteError(w, err)
			return
		}
		api.Write(w, res)
	}
}
