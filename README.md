# Crypto Exchange

A crypto exchange platform written in Golang that simulates an order book for trading cryptocurrencies. The platform supports placing limit and market orders, canceling orders, and retrieving the current state of the order book.

## Features

- **Order Book Management**: Maintains a list of buy and sell orders with price and size.
- **Order Types**:
  - **Limit Orders**: Place orders at a specific price.
  - **Market Orders**: Match orders with the best available price.
- **Order Matching**: Automatically matches buy and sell orders based on price and size.
- **Order Cancellation**: Cancel existing orders by their ID.
- **Ethereum Integration**: Simulates ETH transfers using the Ethereum blockchain.
- **REST API**: Exposes endpoints for interacting with the exchange.

## Project Structure

- **`main.go`**: Entry point of the application. Sets up the REST API and initializes the exchange.
- **`orderbook/orderbook.go`**: Core logic for managing the order book, including placing and matching orders.
- **`orderbook/orderbook_test.go`**: Unit tests for the order book functionality.
- **`util.go`**: Utility functions for interacting with the Ethereum blockchain.
- **`Makefile`**: Build and test automation.
- **`apiTest.http`**: HTTP request examples for testing the REST API.

## API Endpoints

### 1. Place an Order

- **POST** `/order`
- **Request Body**:

  ```json
  {
    "userID": 8,
    "type": "LIMIT",
    "bid": true,
    "size": 20,
    "price": 9000.0,
    "market": "ETH"
  }
  ```

- **Response**:
  - For limit orders: `{ "message": "Limit Order placed successfully" }`
  - For market orders: `{ "matches": [...] }`

### 2. Get Order Book

- **GET** `/book/:market`
- **Response**:

  ```json
  {
    "TotalBidVolume": 100.0,
    "TotalAskVolume": 200.0,
    "Asks": [...],
    "Bids": [...]
  }
  ```

### 3. Cancel an Order

- **DELETE** `/order/:id`
- **Response**:

  ```json
  { "message": "Limit Order Deleted successfully" }
  ```

## Setup Instructions

### Prerequisites

- Go 1.18 or later
- Ethereum client (e.g., Geth or Ganache)
- `make` (optional, for build automation)

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/your-repo/crypto-exchange.git
   cd crypto-exchange
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Start an Ethereum client (e.g., Ganache) locally.

4. Run the application:

   ```bash
   go run main.go
   ```

5. Use the provided `apiTest.http` file to test the API using tools like Postman or VS Code REST Client.

### Using Makefile

- Build and run the application:

```bash
make run
```

- Run tests:

```bash
make test
```

- Clean build artifacts:

```bash
make clean
```

## Testing

Unit tests for the order book are located in `orderbook/orderbook_test.go`

## Ethereum Integration

The platform uses the Ethereum blockchain to simulate ETH transfers. Ensure that your Ethereum client is running and accessible at `http://localhost:8545`. Update the private keys and addresses in `main.go` as needed.

## Example Workflow

1. Place a limit order to buy ETH at a specific price.
2. Place a market order to sell ETH, which matches the existing buy order.
3. Retrieve the order book to view the updated state.
4. Cancel an order by its ID.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a detailed description of your changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- [Echo Framework](https://echo.labstack.com/) for building REST APIs.
- [Go-Ethereum](https://github.com/ethereum/go-ethereum) for Ethereum blockchain integration.
