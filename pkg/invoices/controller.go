package invoices

import (
	"log"

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

	routes := app.Group("/invoices")

	routes.Post("/", h.CreateInvoice)
	routes.Get("/:id/update", h.UpdateTrade)
	routes.Get("/:id", h.GetInvoice)
	routes.Post("/bid", h.PlaceBid)
	
	
}

func (h handler) CreateInvoice(c *fiber.Ctx) error {
	body := models.CreateInvoiceDto{}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var invoice models.Invoice
	var issuer models.User

	invoice.Amount = body.Amount
	invoice.IssuerId = body.IssuerId
	invoice.Reference = "somerandomstring"

	err := h.DB.Model(&models.User{}).Where("id = ? AND user_type = ?", body.IssuerId, "issuer").First(&issuer).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// insert new db entry
	if result := h.DB.Create(&invoice); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&invoice)
}

func (h handler) GetInvoice(c *fiber.Ctx) error {

	id := c.Params("id")
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

	//Build bid
	var bidRequest models.Bid
	bidRequest.InvestorId = investor.ID
	bidRequest.InvoiceId = invoice.ID
	bidRequest.Amount = body.Amount

	//Push bid request to queue

	//TODO
	//Create Bid Model
	//Push to queue
	//In processing
	//Lock thread
	//Read invoice again from db
	//
	processBid(h.DB, bidRequest)

	return c.Status(fiber.StatusOK).JSON("Bid placed successfully")

}

func (h handler) UpdateTrade(c *fiber.Ctx) error {
	id := c.Params("id")

	queryValue := c.Query("action")

	if queryValue == ""{
		return fiber.NewError(fiber.StatusNotFound, "No action added")
	}

	var invoice models.Invoice

	err := h.DB.Model(&models.Invoice{}).Where("id = ?", id).Preload("Bids.Investor").First(&invoice).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if invoice.Status == ""{
		return fiber.NewError(fiber.StatusNotFound, "Trade has not yet been completed")
	}

	var updateFlag string

	updateFlag = queryValue

	//For approval 
	if updateFlag == "approve"{

		// Get Issuer's account and impact balance 
		var issuer models.User

		err := h.DB.Model(&models.User{}).Where("id = ? AND user_type = ?", invoice.IssuerId, "issuer").First(&issuer).Error
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		issuer.Balance += invoice.AmountBided
		h.DB.Save(issuer)

		invoice.Status = "approved"
		h.DB.Save(invoice)

	}else if updateFlag == "reject"{

		//Get all investors involved and roll back balance
		if len(invoice.Bids) <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "No investor participated") // shouldnt happen
		}

		for i := 0; i < len(invoice.Bids); i++ {
			
			var investor models.User
			err = h.DB.Model(&models.User{}).Where("id = ?", invoice.Bids[i].Investor.ID).First(&investor).Error
			investor.Balance += invoice.Bids[i].Amount
			h.DB.Save(investor)
		}


		invoice.Status = "rejected"
		h.DB.Save(invoice)
	}

	return c.Status(fiber.StatusOK).JSON("Done")
}

func processBid(db *gorm.DB, bidRequest models.Bid) {
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
	err = db.Model(&models.User{}).Where("id = ?", bidRequest.InvestorId).First(&investor).Error
	if err != nil {
		return
	}

	if investor.Balance < bidRequest.Amount {
		return
	}

	//Create bid record
	if result := db.Create(&bidRequest); result.Error != nil {
		log.Fatalln(result.Error)
	}


	//Update investor balance

	investor.Balance -= bidRequest.Amount
	db.Save(investor)

	

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
