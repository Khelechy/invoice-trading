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
	invoice.Reference = "somerandomstring"

	// insert new db entry
	if result := h.DB.Create(&invoice); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&invoice)
}

func (h handler) GetInvoice(c *fiber.Ctx) error {

	var id int
	var invoice models.Invoice

	err := h.DB.Model(&models.Invoice{}).Where("id = ?", id).Preload("Bids.Investor").First(&invoice).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&invoice)
}

func (h handler) PlaceBid(c *fiber.Ctx) error {
	body := models.PlaceBid{}

	// parse body, attach to AddBookRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	//Validate sufficient funds for investor

	var investor models.User
	var invoice models.Invoice

	err := h.DB.Model(&models.User{}).Where("id = ?", body.InvestorId).First(&investor).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if investor.Balance < body.Amount {
		return fiber.NewError(fiber.StatusBadRequest, "You have insufficient funds")
	}

	//Validate Invoice exists
	err = h.DB.Model(&models.Invoice{}).Where("id = ?", body.InvoiceId).First(&invoice).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if invoice.Status == "locked" || invoice.Status == "approved" {
		return fiber.NewError(fiber.StatusBadRequest, "Sorry you cant trade this invoice")
	}

	//Push bid request to queue

	//TODO
	//Create Bid Model
	//Push to queue
	//In processing
	//Lock thread
	//Read invoice again from db
	//

	return c.Status(fiber.StatusOK).JSON("Bid placed successfully")

}

func ProcessBid(db *gorm.DB, bidRequest models.Bid) {
	//lock thread
	var invoice models.Invoice
	var investor models.User

	//revalidate invoice status
	err := db.Model(&models.Invoice{}).Where("id = ?", bidRequest.InvoiceId).First(&invoice).Error
	if err != nil {
		return
	}

	if invoice.Status == "locked" || invoice.Status == "approved" {
		return
	}

	//revalidate balance
	err = db.Model(&models.User{}).Where("id = ?", bidRequest.InvestorId).First(&invoice).Error
	if err != nil {
		return
	}

	if investor.Balance < bidRequest.Amount {
		return
	}

	//Update investor balance

	investor.Balance -= bidRequest.Amount
	db.Save(investor)

	//Create bid record
	result := db.Create(bidRequest)

	if result.Error.Error() != "" {
		//perform a rollback
	}

	//Update invoice record
	invoice.AmountBided += bidRequest.Amount
	invoice.Bids = append(invoice.Bids, bidRequest)

	if invoice.AmountBided >= invoice.Amount {

		// Lock invoice
		invoice.Status = "locked"
	}

	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&invoice)

	//unlock thread

}
