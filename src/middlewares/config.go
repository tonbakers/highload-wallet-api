package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"highload-wallet-api/src/api"
	"highload-wallet-api/src/config"
)

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Token        string
}

func New(config config.FileConfig) fiber.Handler {
	cfg := Config{
		Filter: func(c *fiber.Ctx) bool { return config.Server.AllowToken },
		Token:  config.Server.Token,
	}

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			fmt.Println("Authorization middleware was skipped.")
		}
		fmt.Println("Authorization middleware is enabled.")

		check := func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(
					api.Apierrs.ErrorInvalidAuthHeader.Code,
				).JSON(api.Apierrs.ErrorInvalidAuthHeader)
			}
			if !strings.Contains(authHeader, "Bearer") {
				return c.Status(
					api.Apierrs.ErrorInvalidAuthHeader.Code,
				).JSON(api.Apierrs.ErrorInvalidAuthHeader)
			}
			values := strings.Split(authHeader, " ")
			token := strings.TrimSpace(values[1])
			if token != cfg.Token {
				return c.Status(
					api.Apierrs.ErrorUnauthorized.Code,
				).JSON(api.Apierrs.ErrorUnauthorized)
			}
			return nil
		}
		return check(c)
	}
}
