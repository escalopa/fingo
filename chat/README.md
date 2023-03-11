# Chat Service 💬

This service is responsible for handling all the chat related operations, such as sending messages, receiving messages, etc.

## Features 🚀


## Flow

* **Send Messages**
  - User must be authenticated to send messages
  - Messages are should be updated in real time using websockets

```mermaid
sequenceDiagramautonumber
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
