## Invoice Trading

#### To run this project

- You must have Golang and Postgres installed on your machine
- Clone the repo to your pc 
- Create a postgres database and update DB_URL the dev.env config file in the project ``./pkg/common/config/envs

-- example: postgres://username:password@localhost:5432/dbname

- Change directory into the cmd folder and run the main file

-- example:
 cd cmd
 go run main.go
 
 ### Simulating a successful Invoice trading flow

 Call the following endpoints 
 Postman collection: 

### Create Two users, Investor and Issuer 

##### Create User
 http://localhost:3000/users `POST`

```json
{
    "name":"Kelechi",
    "balance": 1000,
    "user_type":"investor"
}
```

```json
{
    "name":"AREX Corp",
    "balance": 0,
    "user_type":"issuer"
}
```

** "user_type" takes two values "issuer" and "investor" **

##### Create Invoice 
 http://127.0.0.1:3000/invoices/

 ```json
 {
    "amount": 500,
    "issuer_id": 2
}
```

#### User management endpoints

##### Get All Investors and Balance

 http://127.0.0.1:3000/users/investors

 ```json
 [
    {
        "ID": 1,
        "CreatedAt": "2023-06-15T22:38:19.606949+01:00",
        "UpdatedAt": "2023-06-15T23:55:13.101905+01:00",
        "DeletedAt": null,
        "name": "kelechi",
        "balance": 500,
        "user_type": "investor"
    },
    {
        "ID": 4,
        "CreatedAt": "2023-06-16T18:54:38.683635+01:00",
        "UpdatedAt": "2023-06-16T18:54:38.683635+01:00",
        "DeletedAt": null,
        "name": "john",
        "balance": 1000,
        "user_type": "investor"
    }
]
 ```

 #### Get Back an Issuer
  http://127.0.0.1:3000/users/issuers/2

  ```json
  {
    "ID": 2,
    "CreatedAt": "2023-06-15T22:38:31.794711+01:00",
    "UpdatedAt": "2023-06-15T23:55:36.953562+01:00",
    "DeletedAt": null,
    "name": "kcorp",
    "balance": 6500,
    "user_type": "issuer"
}
  ```


##### Place Bid 
 http://127.0.0.1:3000/invoices/bid
 
 ```json
 {
    "amount":1000,
    "investor_id":4,
    "invoice_id": 6
}
```

#### Fetch Back Invoice,status and associated bids
 http://127.0.0.1:3000/invoices/6

 `response` 
 ```json
 {
    "ID": 5,
    "CreatedAt": "2023-06-15T23:54:54.189733+01:00",
    "UpdatedAt": "2023-06-15T23:55:36.954197+01:00",
    "DeletedAt": null,
    "amount": 500,
    "amount_bided": 500,
    "issuer_id": 2,
    "reference": "somerandomstring",
    "status": "approved",
    "bids": [
        {
            "ID": 8,
            "CreatedAt": "2023-06-15T23:55:13.099163+01:00",
            "UpdatedAt": "2023-06-15T23:55:13.099163+01:00",
            "DeletedAt": null,
            "amount": 500,
            "invoice_id": 5,
            "investor_id": 1,
            "investor": {
                "ID": 1,
                "CreatedAt": "2023-06-15T22:38:19.606949+01:00",
                "UpdatedAt": "2023-06-15T23:55:13.101905+01:00",
                "DeletedAt": null,
                "name": "kelechi",
                "balance": 500,
                "user_type": "investor"
            }
        }
    ]
}
```

##### Update invoice trade
 http://127.0.0.1:3000/invoices/6/update?action=approve
** You can set the action to either "approve" or "reject" **

