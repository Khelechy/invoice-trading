package models

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	refreshUserTable()

	newUser := User{
		Name:     "Kelechi",
		Balance:  5000,
		UserType: "investor",
	}

	want := "Kelechi"

	got, err := CreateUser(db, newUser)
	if err != nil {
		t.Errorf("Error creating user: %v\n", err)
	}

	if got.Name != want {
		t.Errorf("got %q wanted %q", got.Name, want)
	}
}

func TestGetIssuer(t *testing.T){
	refreshUserTable()
	seedInvestorandIssuer()

	want := "KCorp"
	got, err := GetIssuer(db, 2)

	if err != nil {
		t.Errorf("Error getting issuer: %v\n", err)
	}

	if got.Name != want {
		t.Errorf("got %q wanted %q", got.Name, want)
	}

	if got.UserType != "issuer"{
		t.Errorf("got %q wanted %q", got.UserType, "issuer")
	}
}

func TestGetInvestor(t *testing.T){
	refreshUserTable()
	seedInvestorandIssuer()

	want := "Kelechi"
	got, err := GetInvestor(db, 1)

	if err != nil {
		t.Errorf("Error get investor: %v\n", err)
	}

	if got.Name != want {
		t.Errorf("got %q wanted %q", got.Name, want)
	}

	if got.UserType != "investor"{
		t.Errorf("got %q wanted %q", got.UserType, "investor")
	}
}

func TestGetInvestors(t *testing.T){
	seedInvestors()

	got, err := GetInvestors(db)
	if err != nil {
		t.Errorf("Error getting investors: %v\n", err)
	}

	if len(got) < 1 {
		t.Errorf("")
	}

	if got[0].UserType != "investor"{
		t.Errorf("got %q wanted %q", got[0].UserType, "investor")
	}
}

