# User Service 👤

This service is responsible for storing and managing user’s information & updating all the user’s info

## Features 🚀

### User

- Get user information
- Update user’s email, password
- Restore forgotten password
- User’s Info
  - First name
  - Last name
  - Username
  - Email

## Flow 🌊

* **Get User Info**
  - User's info are retrieved from the database by the user's ID.
  - User's info are returned to the user after filtering the sensitive data.

```mermaid
sequenceDiagram
    autonumber
    API->>+User Service: Send user's ID
    User Service->>+Database: Get user's account
    Database-->>-User Service: User's account
    User Service-->>-API: User's info
```

* **Update User Info**
  - User's info are updated right away in the database after validating the user's password.
  - User's info are returned to the user after filtering the sensitive data.

```mermaid
sequenceDiagram
  autonumber
  API->>+User Service: Send user's ID & info
  Note over API, User Service: Info is not email or password
  User Service->>+Database: Update user's info
  Database-->>-User Service: User's info updated
  User Service-->>-API: Success response
```

* **Update User Email**
  - User's email is NOT updated right away.
  - User must confirm his new email first so that changes take place.

```mermaid
sequenceDiagram
  autonumber
  API->>+User Service: Send user's ID & new email
  User Service->>+Database: Get user's account
  Database-->>-User Service: User's account
  User Service->>User Service: Generate confirmation code
  User Service->>+Cache: Store confirmation code
  Cache-->>-User Service: Confirmation code stored
  User Service->>Message Broker: Send confirmation code to user's new email/phone
  User Service->>API: Success response
```

* **Verify User Email**

- User's email is updated after verifying the confirmation code.
- User's email is update in the database.

```mermaid
sequenceDiagram
  autonumber
  API->>+User Service: Send user's ID & confirmation code
  User Service->>+Database: Get user's account
  Database-->>-User Service: User's account
  User Service->>+Cache: Get confirmation code
  Cache-->>-User Service: Confirmation code
  User Service->>User Service: Validate confirmation code
  User Service->>+Database: Update user's account
  Database-->>-User Service: User's account updated
  User Service-->>-API: Success response
```

* **Update User Password**
  - User's password is updated right away in the database after validating the old password.
  - Password is hashed before storing it in the database.

```mermaid
sequenceDiagram
  autonumber
  API->>+User Service: Send user's ID & passwords
  Note over API, User Service: Old & new passwords are sent
  User Service->>+Database: Get user's account
  Database-->>-User Service: User's account
  User Service->>User Service: Validate old password
  User Service->>User Service: Hash new password
  User Service->>+Database: Update user's passoword
  Database-->>-User Service: User's password updated
  User Service-->>-API: Success response
```
* **Forgot Password**
  - User provides his email with which he signed up.
  - Email user's account with a reset password token.
  - User can reset his password with the reset password token and his new password.

```mermaid
sequenceDiagram
    autonumber
    API->>+User Service: Send user's ID, email
    User Service->>+Database: Get user's account
    Database-->>-User Service: User's account
    User Service->>User Service: Generate reset password token
    User Service->>+Cache: Store reset password token
    Cache-->>-User Service: Reset password token stored
    User Service->>Message Broker: Send email to user
    Note over User Service, Message Broker: With reset password token
    User Service-->>-API: Success response
    API->>+User Service: Send user's new password
    Note over API, User Service: With reset password token
    User Service->>+Cache: Get reset password token
    Cache-->>-User Service: Reset password token
    User Service->>User Service: Validate reset password token
    User Service->>User Service: Hash user's password
    User Service->>+Database: Update user's account
    Database-->>-User Service: User's account updated
    User Service-->>-API: Success response
```
