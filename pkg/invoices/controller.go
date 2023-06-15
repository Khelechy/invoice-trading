package invoices

import (
	"log"
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

	invoice.Amount = body.Amount
	invoice.IssuerId = body.IssuerId
	invoice.Reference = "somerandomstring"

	var issId string = strconv.FormatUint(uint64(invoice.IssuerId), 10)
	_, err := models.GetIssuer(h.DB, issId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// insert new db entry
	id, err := models.CreateInvoice(h.DB, invoice)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	invoice.ID = id

	return c.Status(fiber.StatusCreated).JSON(&invoice)
}

func (h handler) GetInvoice(c *fiber.Ctx) error {

	id := c.Params("id")

	invoice, err := models.GetInvoice(h.DB, id)
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

	var invId string = strconv.FormatUint(uint64(body.InvestorId), 10)
	investor, err := models.GetInvestor(h.DB, invId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if investor.Balance < body.Amount {
		return fiber.NewError(fiber.StatusBadRequest, "You have insufficient funds")
	}

	//Validate Invoice exists
	var invoiceId string = strconv.FormatUint(uint64(body.InvoiceId), 10)
	invoice, err := models.GetInvoice(h.DB, invoiceId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if invoice.Status == "locked" || invoice.Status == "approved" {
		return fiber.NewError(fiber.StatusBadRequest, "Sorry you cant trade this invoice")
	}

	//Trim bidrequest if amount is above available invoice bid-able

	totalAmount := invoice.Amount
	amountBided := invoice.AmountBided
	availableToBid := totalAmount - amountBided
	amountToBid := body.Amount

	if amountToBid > availableToBid {
		amountToBid = availableToBid
	}

	//Build bid
	var bidRequest models.Bid
	bidRequest.InvestorId = investor.ID
	bidRequest.InvoiceId = invoice.ID
	bidRequest.Amount = amountToBid

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

	if queryValue == "" {
		return fiber.NewError(fiber.StatusNotFound, "No action added")
	}

	invoice, err := models.GetInvoice(h.DB, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if invoice.Status == "" {
		return fiber.NewError(fiber.StatusNotFound, "Trade has not yet been completed")
	}

	if invoice.Status == "locked" {

		var updateFlag string

		updateFlag = queryValue

		//For approval
		if updateFlag == "approve" {

			// Get Issuer's account and impact balance

			var issId string = strconv.FormatUint(uint64(invoice.IssuerId), 10)
			issuer, err := models.GetIssuer(h.DB, issId)
			if err != nil {
				return fiber.NewError(fiber.StatusNotFound, err.Error())
			}

			issuer.Balance += invoice.AmountBided
			h.DB.Save(&issuer)

			invoice.Status = "approved"
			h.DB.Save(&invoice)

		} else if updateFlag == "reject" {

			//Get all investors involved and roll back balance
			if len(invoice.Bids) <= 0 {
				return fiber.NewError(fiber.StatusNotFound, "No investor participated") // shouldnt happen
			}

			for i := 0; i < len(invoice.Bids); i++ {

				var invId string = strconv.FormatUint(uint64(invoice.Bids[i].Investor.ID), 10)
				investor, _ := models.GetInvestor(h.DB, invId)
				investor.Balance += invoice.Bids[i].Amount
				h.DB.Save(&investor)
			}

			invoice.Status = "rejected"
			h.DB.Save(&invoice)
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON("Trade has already been closed")
}

func processBid(db *gorm.DB, bidRequest models.Bid) {
	//lock thread

	//revalidate invoice status
	var invoiceId string = strconv.FormatUint(uint64(bidRequest.InvoiceId), 10)
	invoice, err := models.GetInvoice(db, invoiceId)
	if err != nil {
		return
	}

	if invoice.Status == "locked" || invoice.Status == "approved" {
		return
	}

	//revalidate balance
	var invId string = strconv.FormatUint(uint64(bidRequest.InvestorId), 10)
	investor, err := models.GetInvestor(db, invId)
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
	db.Save(&investor)

	//Update invoice record
	invoice.AmountBided += bidRequest.Amount
	invoice.Bids = append(invoice.Bids, bidRequest)

	if invoice.AmountBided >= invoice.Amount {
		// Lock invoice
		invoice.Status = "locked"
	}

	db.Save(&invoice)

	//unlock thread

}
