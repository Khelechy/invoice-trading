package invoices

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khelechy/invoice-trading/pkg/common/models"

	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func (h handler) CreateInvoice(c *fiber.Ctx) error {
	body := models.CreateInvoiceDto{}

	// parse body, attach to AddBookRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var invoice models.Invoice

	invoice.Amount = body.Amount
	invoice.IssuerId = body.IssuerId

	// insert new db entry
	if result := h.DB.Create(&invoice); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&invoice)
}

func (h handler) GetInvoice(c *fiber.Ctx) error {

	var id string
	var invoice models.Invoice

	err := h.DB.Model(&models.Invoice{}).Where("id = ?", id).Preload("Bids.Investor").First(&invoice).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&invoice)
}
