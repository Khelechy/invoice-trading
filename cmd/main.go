package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/khelechy/invoice-trading/pkg/common/config"
	"github.com/khelechy/invoice-trading/pkg/common/db"
	"github.com/khelechy/invoice-trading/pkg/common/models"
	"github.com/khelechy/invoice-trading/pkg/invoices"
	"github.com/khelechy/invoice-trading/pkg/users"
)

func main() {
	c, err := config.LoadConfig()
	var m sync.Mutex

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	app := fiber.New()
	db := db.Init(c.DBUrl)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString(c.Port)
	})

	//Setup job queue through channel
	jobChan := make(chan models.Bid)

	invoices.RegisterRoutes(app, db, jobChan)
	users.RegisterRoutes(app, db)

	go worker(db, jobChan, &m)

	app.Listen(c.Port)
}

func worker(db *gorm.DB, jobChan <-chan models.Bid, m *sync.Mutex) {
	for job := range jobChan {
		invoices.ProcessBid(db, job, m)
	}
}
