/*
highload-wallet-api – API wrapper over high-load TON wallet smart contract

Copyright (C) 2021 Alexander Gapak

This file is part of highload-wallet-api.

highload-wallet-api is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

highload-wallet-api is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with highload-wallet-api.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"

	"highload-wallet-api/src/api"
	"highload-wallet-api/src/config"
	"highload-wallet-api/src/middlewares"
)

func main() {
	config.Configure()

	engine := html.New("./static", ".html")

	app := fiber.New(fiber.Config{
		Prefork: config.Cfg.Server.Prefork,
		Views:   engine,
	})
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Authorization, Content-Type, Origin",
		AllowCredentials: false,
		AllowOrigins: "*",
		AllowMethods: "*",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${path} ${method} ${status}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	router := app.Group("/")

	router.Get(
		"/",
		func(c *fiber.Ctx) error {
			return c.Render("HLform", fiber.Map{
				"Title": "Tranfser view",
			})
		},
	).Get(
		"/info",
		middlewares.New(config.Cfg),
		func(c *fiber.Ctx) error {
			walletinfo, err := os.ReadFile("contract/generated/wallet-info.txt")
			if err != nil {
				return c.SendString("error reading wallet-info.txt")
			}
			return c.SendString(string(walletinfo))
		},
	).Get(
		"/activate",
		middlewares.New(config.Cfg),
		func(c *fiber.Ctx) error {
			cmd := exec.Command("contract/activate-wallet.sh", "https://toncenter.com/api/v2/jsonRPC")
			stdout, err := cmd.Output()

			if err != nil {
				return c.SendString(err.Error())
			}

			return c.SendString(string(stdout))
		},
	).Post(
		"/transfer",
		middlewares.New(config.Cfg),
		api.Transfer,
	)

	log.Fatal(app.Listen(config.Cfg.Server.Host + ":" + config.Cfg.Server.Port))
}
