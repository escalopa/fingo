# gochat 💬

Live chatting distributed system service built with go, websockets, gRPC

## Functionalities & Features 🚀

- Chat with other users as in default messenger apps
- Chats are updated in real time as someone send a message
- Robust authentication with email verification
- Api endpoints to frontend consumption

## Micro-Services Architecture 🏗

Communication between all the microservices is done using `grpc`

The following services are responsible for

- **Chat**: Storing message between users
- **Email**: Sending and email for verifications
- **Auth** CRUD account and token handling
- **API** Providing endpoints to interact with the application

## How to run 🏃‍♂️

## Diagrams & BP 📝

* **Create An Account**
  - User creates a new account with email & password
  - Account confirmation is required to use the application, So a code is sent to your email upon creation.

```mermaid
sequenceDiagram
    Client->>Server: Send new account Info
    activate Server
    Server->>Email Service: Send Confirmation Code
    activate Email Service
    Note over Server,Email Service: Validate email
    deactivate Email Service
    Server->>Database: Create Account
    Server-->>Client: Account Created
    deactivate Server
    Client->>Server: Send Confirmation Code
    activate Server
    activate Email Service
    Server->>Email Service: Verify confirmation code
    Email Service-->>Server: Confirmation Code Verified
    deactivate Email Service
    Server->>Database: Update Account Status
    Server-->>Client: Account activated
    deactivate Server
```

* **Send Messages**
  - User must be authenticated to send messages
  - Messages are should be updated in real time using websockets

```mermaid
sequenceDiagram
    actor Alice
    actor Bob
    Alice->>Server: Send Message to bob
    activate Server
    activate Database
    Server->>Database: Store Message
    deactivate Database
    Server->>Message Queue: Produce a new message in the queue
    Server-->>Alice: Message sent successfully response
    deactivate Server
    activate Message Queue
    Bob->>Message Queue: Consume Message sent to him from the queue
    Message Queue->>Bob: Send delivered message to the user
    deactivate Message Queue
```

## Components Diagram 📊

```mermaid
graph TD
    A[Client] -->| REST | B[API]
    B -->| gRPC | D[Auth Service]
    B -->| gRPC | C[Chat]
    D -->| gRPC | H[Email Service]
    D -->| SQL/NoSQL | E[Database]
    C -->| SQL/NoSQL | F[Database]
    C --> G[Message Queue]
    G --> | Web Sockets | A
```
