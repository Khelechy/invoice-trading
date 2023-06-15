package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/khelechy/invoice-trading/pkg/common/config"
	"github.com/khelechy/invoice-trading/pkg/common/db"
	"github.com/khelechy/invoice-trading/pkg/invoices"
	"github.com/khelechy/invoice-trading/pkg/users"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	app := fiber.New()
	db := db.Init(c.DBUrl)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString(c.Port)
	})

	invoices.RegisterRoutes(app, db)
	users.RegisterRoutes(app, db)

	app.Listen(c.Port)
}
