# Auth Service 🔑

This service is responsible for creating users & user’s session CRUD(authentication).

## Features 🚀

- Sign-up
- Sign-in
- Logout
- Renew auth token by refresh token

## Flow 🌊

* **Sign-up**
  - User creates a new account with email & password as login credentials.
  - Account confirmation(Email, Phone) is required to use the application.
  - User's info are passed along the request to the auth service.

```mermaid
sequenceDiagram
    autonumber
    API->>Auth Service: Send new account's Info
    activate API
    activate Auth Service
    Auth Service->>Database: Create Account
    activate Database
    Database-->>Auth Service: Account Created
    deactivate Database
    Auth Service-->>API: Account Created
    deactivate Auth Service
    deactivate API
```

* **Sign-in**
  - User signs in with email & password.
  - Generates a new auth token and refresh token for the user.
  - Create a new session for the user in the database.
  - Notifies the user about the new login session by sending an email.

```mermaid
sequenceDiagram
    autonumber
    API->>+Auth Service: Send user's credentials
    Auth Service->>+Database: Get user's account
    Database-->>-Auth Service: User's account
    Auth Service->>+Database: Create new session
    Auth Service->>Auth Service: Validate user's password
    Auth Service->>Auth Service: Generate auth token & refresh token
    Auth Service->>+Database: Create new user's session
    Database-->>-Auth Service: Session created
    Auth Service->>Message Broker: Send email to user
    Note over Auth Service, Message Broker: Send email about the new login session
    Auth Service-->>-API: Auth token & refresh token
```

* **Logout**
  - User signs out from the application.
  - Revokes the user's auth token.
  - Deletes the user's session from the database.

```mermaid
sequenceDiagram
    autonumber
    API->>+Auth Service: Send the session's info to logout from
    Auth Service->>+Database: Get user's session
    Database-->>-Auth Service: User's session
    Auth Service->>+Database: Delete user's session
    Database-->>-Auth Service: Session deleted
    Auth Service-->>-API: User logged out successfully
```

* **Validate Auth Token**
  - Check the token's lifetime has not expired.
  - Check the user's session in the database.
  - Check the user's session is not revoked, expired, or deleted.
  - Return the user's info(ID,Name,Email,Phone...).

```mermaid
sequenceDiagram
    autonumber
    API->>+Auth Service: Send the auth token to validate
    Auth Service->>Auth Service: Validate access token
    Auth Service->>+Database: Get user's session
    Database-->>-Auth Service: User's session
    Auth Service->>Auth Service: Validate user's session
    Auth Service-->>-API: User's info
```

* **Renew Auth Token(Access, Refresh)**
  - Check the refresh token's lifetime has not expired.
  - Check the user's session in the database.
  - Check the user's session is not revoked, expired, or deleted.
  - Generate a new auth token and refresh token for the user.
  - Update the user's session in the database.
  - Return the new auth token and refresh token.

```mermaid
sequenceDiagram
    autonumber
    API->>+Auth Service: Send the refresh token to renew
    Auth Service->>Auth Service: Validate refresh token
    Auth Service->>+Database: Get user's session
    Database-->>-Auth Service: User's session
    Auth Service->>Auth Service: Validate user's session
    Auth Service->>Auth Service: Generate auth token & refresh token
    Auth Service->>+Database: Update user's session
    Database-->>-Auth Service: Session updated
    Auth Service-->>-API: New auth token & refresh token
```

## Database 🗄

![Database](./docs/database.png)
