# Wallet Service 💰

The service is responsible for storing user’s accounts and manage balances changes.

## Database 🗄

![Diagram](./../docs/fingo_wallet_db.png)

## Features 🚀

### Account
 - [x] Create a user account with a specific currency & name.
 - [x] Get all user's accounts.

### Cards
 - [x] Create a new card for payment.
 - [x] Links a card to account.
 - [x] Get all cards for a specific account.

### Currency
 - [x] Support differecnt currencies(USD, RUB, EGP, GBP, EUR)

## Flow 🌊

* **CreateWallet**
  - Map external user's id to another internal user'id for the current database.

```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make create wallet request
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Database: Create new wallet
    Database-->>-Wallet Service: Wallet created
    Wallet Service-->>-API: Wallet created
```

* **CreateAccount**
  - Create a new account for the user.
  - External user's id is mapped to internal user's id.
  - Internal user's id is used to create a new account.

```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make create account request
    Note over API, Wallet Service: Pass account name & currency
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Database: Create new account
    Database-->>-Wallet Service: Account created
    Wallet Service-->>-API: Account created
```

* **GetAccounts**
  - Get all accounts for the user.
  - Doesn't support pagination.

```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make get accounts request
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Database: Get user's accounts
    Note over Wallet Service, Database: Get all accounts for mapped user's external id
    Database-->>-Wallet Service: User's accounts
    Wallet Service-->>-API: User's accounts
```

* **DeleteAccount**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make delete account request
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Wallet Service: Validate account owner & non empty balance
    Wallet Service->>+Database: Delete account
    Wallet Service-->>-API: Account deleted
```

* **CreateCard**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make create card request
    Note over API, Wallet Service: Pass account id
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Wallet Service: Validate account owner & create card
    Wallet Service->>+Database: Store card
    Database-->>-Wallet Service: Card stored
    Wallet Service-->>-API: Card created
```

* **GetCards**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make get cards request
    Note over API, Wallet Service: Pass account id
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>Wallet Service: Validate account owner
    Wallet Service->>+Database: Get cards
    Database-->>-Wallet Service: Account cards
    Wallet Service-->>-API: Account cards
```

* **DeleteCard**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make delete card request
    Note over API, Wallet Service: Pass card number
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>Wallet Service: Validate card owner
    Wallet Service->>+Database: Delete card
    Database-->>-Wallet Service: Card deleted
    Wallet Service-->>-API: Card deleted
```

* **CreateTransaction**
  - Transaction can be created with 3 different types:
    - **Transfer** - Transfer money from one account to another.
    - **Deposit** - Deposit money to an account.
    - **Withdraw** - Withdraw money from an account.
  - ONLY on transfer transaction  the receiver's account should be passed

```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make create transaction request
    Note over API, Wallet Service: Pass your's &receiver's card  number & amount
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>Wallet Service: Validate card owner
    Wallet Service->>+Database: Get accounts
    Database-->>-Wallet Service: Accounts
    Wallet Service->>Wallet Service: Validate accounts are not the same
    Wallet Service->>+Database: Validate sender's account balance is enough
    Wallet Service->>+Database: Create transaction
    Database-->>-Wallet Service: Transaction created
    Wallet Service->>-API: transaction created
```

* **TransferRollback**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make transfer rollback request
    Note over API, Wallet Service: Pass transaction id
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>+Database: Get transaction
    Database-->>-Wallet Service: Transaction
    Wallet Service->>Wallet Service: Validate transaction type is transfer
    Wallet Service->>Wallet Service: Validate caller is the sender
    Wallet Service->>+Database: Rollback transaction
    Database-->>-Wallet Service: Transaction rolled back
    Wallet Service-->>-API: Transaction rolled back
```

* **GetTransactionHistory**
  -
```mermaid
sequenceDiagram
    autonumber
    API->>+Wallet Service: Make get transaction history request
    Note over API, Wallet Service: Pass account id & pagination & optional filter
    Wallet Service->>+Token Service: Validate token
    Token Service-->>-Wallet Service: User external id
    Wallet Service->>Wallet Service: Validate account owner
    Wallet Service->>+Database: Get transactions
    Database-->>-Wallet Service: Transactions
    Wallet Service-->>-API: Transactions
```
