
# coinbit-wallet
This project is simple deposit and get wallet balance application to fulfill Stockbit online test for Backend Engineer - Golang (Blockchain Platform) position. The application is consists of two endpoints, POST request to deposit certain amount to a wallet and GET request to get details of wallet including wallet_id, balance and a flag whether the wallet has ever done one or more deposits with amounts more than 10,000 within a single 2-minute window. 

## Tech Used
- Golang 1.20
- Kafka

## API List
|API|Routes|Method|
|----------------|-------------------------------|-----------------------------|
|Deposit|/api/v1/wallet/deposit  |POST   | 
|Get Details   |/api/v1/wallet/details/:wallet_id     |GET     |

## Getting Started
```
#  Unzip project folder from the zip file

#  Get into the project directory
cd coinbit-wallet

#  Copy .env-example and modify the .env file manually
cp .env-example .env

#  Using Docker Compose
docker compose up -d
```

## Example
 **Deposit Request Payload**

    {
    "wallet_id": "333-999",
    "amount": 6000
    }

**Deposit Response**

    {
    "message": "SUCCESS",
    "status": 200
    }

**Get Details Response**

    {
    "data": {
        "wallet_id": "333-999",
        "balance": 18000,
        "above_threshold": false
    },
    "message": "SUCCESS",
    "status": 200
    }

> Written with [StackEdit](https://stackedit.io/).