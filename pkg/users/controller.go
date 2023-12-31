package users

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/khelechy/invoice-trading/pkg/common/models"

	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := app.Group("/users")
	routes.Post("/", h.CreateUser)
	routes.Get("/investors", h.GetInvestors)
	routes.Get("/issuers/:id", h.GetIssuer)
}

func (h handler) CreateUser(c *fiber.Ctx) error {
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// insert new db entry
	newUser, err := models.CreateUser(h.DB, user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&newUser)
}

func (h handler) GetIssuer(c *fiber.Ctx) error {
	id := c.Params("id")

	userId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	user, err := models.GetIssuer(h.DB, uint(userId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h handler) GetInvestors(c *fiber.Ctx) error {
	users, err := models.GetInvestors(h.DB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&users)
}
