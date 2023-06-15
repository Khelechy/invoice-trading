package users

import (
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
	if result := h.DB.Create(&user); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&user)
}

func (h handler) GetIssuer(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	err := h.DB.Model(&models.User{}).Where("id = ? AND user_type = ?", id, "issuer").First(&user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h handler) GetInvestors(c *fiber.Ctx) error {
	var users []models.User

	err := h.DB.Model(&models.User{}).Where("user_type = ?", "investor").Find(&users).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&users)
}
