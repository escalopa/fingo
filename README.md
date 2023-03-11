# fingo 💸

[![Go Report Card](https://goreportcard.com/badge/github.com/escalopa/fingo)](https://goreportcard.com/report/github.com/escalopa/fingo)
[![codecov](https://codecov.io/gh/escalopa/fingo/branch/master/graph/badge.svg?token=QZQZQZQZQZ)](https://codecov.io/gh/escalopa/fingo)
[![Build Status](https://travis-ci.com/escalopa/fingo.svg?branch=master)](https://travis-ci.com/escalopa/fingo)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub issues](https://img.shields.io/github/issues/escalopa/fingo.svg)](https://github.com/escalopa/fingo/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/escalopa/fingo.svg)](https://github.com/escalopa/fingo/pulls)


## Features 🚀

TODO

## Microservices Architecture 🏗

Communication between all the microservices is done using `grpc` or `message brokers`.

In fingo we have the following services, Where each one is responsible for a specific set of tasks.

1. [**API**](./api)
2. [**Auth**](./auth)
3. [**User**](./user)
4. [**Chat**](./chat)
5. [**Wallet**](./wallet)
6. [**Payment**](./payment)
7. [**Email**](./email)
8. [**Phone**](./phone)

Our microservices are also powered by the best logging, tracing applications like

1. OpenTelemetry
2. Datadog
3. Prometheus

## Project Architecture 🏘

![Diagram](./fingo.png)

## How to run 🏃‍♂️

### Prerequisites

1. [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
2. [Docker](https://docs.docker.com/get-docker/)
3. [Docker Compose](https://docs.docker.com/compose/install/)
4. [Go](https://golang.org/doc/install)

### Running the project

First clone the project

```bash
git clone github.com/escalopa/fingo && cd fingo
```

Copy

- `.env.example` file to `.env` and fill in the empty fields.
- `.db.env.example` file to `.db.env` and fill in the empty fields.

```bash
cp .env.example .env
cp .db.env.example .db.env
```

Run the project

```bash
docker-compose up
```
