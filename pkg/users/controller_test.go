package users

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/khelechy/invoice-trading/pkg/common/models"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func setUpDatabase() {

	psqldb, err := gorm.Open(postgres.Open("postgres://postgres:root@localhost:5432/invoice_trading_db_test"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db = psqldb
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		// First test case
		{
			description:  "Create Users",
			route:        "/users/",
			expectedCode: 200,
		},
	}

	newUser := models.User{
		Name:     "Kelechi",
		Balance:  5000,
		UserType: "investor",
	}

	requestBody, _ := json.Marshal(&newUser)

	var c *fiber.Ctx
	c.Request().SetBody(requestBody)

	// Define Fiber app.
	app := fiber.New()

	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("POST", test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

}

func Test_handler_CreateUser(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}

	newUser := models.User{
		Name:     "Kelechi",
		Balance:  5000,
		UserType: "investor",
	}

	requestBody, _ := json.Marshal(newUser)

	app := fiber.New()
    c := app.AcquireCtx(&fasthttp.RequestCtx{})
	c.Request().SetBody(requestBody)
	c.Request().Header.Add("Content-Type", "application/json")


	h := handler{
		DB: db,
	}

	argss := args{
		c: c,
	}

	tests := []struct {
		name    string
		h       handler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Create user",
			h:    h,
			args: argss,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.CreateUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("handler.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
