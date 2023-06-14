package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khelechy/invoice-trading/pkg/common/models"

	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func (h handler) GetIssuer(c *fiber.Ctx) error {
	var id string
	var user models.User

	err := h.DB.Model(&models.Invoice{}).Where("id = ? AND usertype", id, "issuer").First(&user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h handler) GetInvestors(c *fiber.Ctx) error {
	var users []models.User

	err := h.DB.Model(&models.Invoice{}).Where("usertype = ?", "investor").Find(&users).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&users)
}
