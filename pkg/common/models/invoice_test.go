package models

import (
	"testing"
)

func TestCreateInvoice(t *testing.T) {
	refreshUserTable()
	refreshInvoiceTable()

	_, issuer, _ := seedInvestorandIssuer()

	newinvoice := Invoice{
		Amount: 2000,
		IssuerId: issuer.ID,
		Reference: "somerandomstring",
	}

	want := "somerandomstring"

	got, err := CreateInvoice(db, newinvoice)
	if err != nil {
		t.Errorf("Error creating invoice: %v\n", err)
	}

	if got.Reference != want {
		t.Errorf("got %q wanted %q", got.Reference, want)
	}
}

func TestGetInvoice(t *testing.T){
	refreshInvoiceTable()
	seedInvoice()

	want := "somerandomstring"
	got, err := GetInvoice(db, 1)

	if err != nil {
		t.Errorf("Error getting invoice: %v\n", err)
	}

	if got.Reference != want {
		t.Errorf("got %q wanted %q", got.Reference, want)
	}

}

