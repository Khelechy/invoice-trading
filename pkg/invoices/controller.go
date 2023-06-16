package invoices

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/khelechy/invoice-trading/pkg/common/models"

	"gorm.io/gorm"
)

type handler struct {
	DB      *gorm.DB
	JobChan chan models.Bid
}

func RegisterRoutes(app *fiber.App, db *gorm.DB, jobChan chan models.Bid) {
	h := &handler{
		DB:      db,
		JobChan: jobChan,
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

	_, err := models.GetIssuer(h.DB, invoice.IssuerId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// insert new db entry
	newInvoice, err := models.CreateInvoice(h.DB, invoice)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(&newInvoice)
}

func (h handler) GetInvoice(c *fiber.Ctx) error {

	id := c.Params("id")
	invoiceId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	invoice, err := models.GetInvoice(h.DB, uint(invoiceId))
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

	investor, err := models.GetInvestor(h.DB, body.InvestorId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if investor.Balance < body.Amount {
		return fiber.NewError(fiber.StatusBadRequest, "You have insufficient funds")
	}

	//Validate Invoice exists
	invoice, err := models.GetInvoice(h.DB, body.InvoiceId)
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
	h.JobChan <- bidRequest

	return c.Status(fiber.StatusOK).JSON("Bid placed successfully")

}

func (h handler) UpdateTrade(c *fiber.Ctx) error {
	id := c.Params("id")

	invoiceId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	queryValue := c.Query("action")

	if queryValue == "" {
		return fiber.NewError(fiber.StatusNotFound, "No action added")
	}

	invoice, err := models.GetInvoice(h.DB, uint(invoiceId))
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
			issuer, err := models.GetIssuer(h.DB, invoice.IssuerId)
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
				investor, _ := models.GetInvestor(h.DB, invoice.Bids[i].Investor.ID)
				investor.Balance += invoice.Bids[i].Amount
				h.DB.Save(&investor)
			}

			invoice.Status = "rejected"
			h.DB.Save(&invoice)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Trade has already been closed")
	}

	return c.Status(fiber.StatusBadRequest).JSON("Trade has been updated")
}

func ProcessBid(db *gorm.DB, bidRequest models.Bid) {
	//lock thread

	//revalidate invoice status
	invoice, err := models.GetInvoice(db, bidRequest.InvoiceId)
	if err != nil {
		return
	}

	if invoice.Status == "locked" || invoice.Status == "approved" {
		return
	}

	//revalidate balance
	investor, err := models.GetInvestor(db, bidRequest.InvestorId)
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
