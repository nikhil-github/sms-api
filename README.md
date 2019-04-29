# sms-app

This code is 100% in compliance with golint, go_vet and gofmt. Check this for more details: [![Go Report Card](https://goreportcard.com/badge/github.com/nikhil-github/sms-app)](https://goreportcard.com/report/github.com/nikhil-github/sms-app) [![Build](https://travis-ci.org/nikhil-github/sms-app.svg?branch=master)](https://travis-ci.org/nikhil-github/sms-app)


## Introduction:

A simple app built in golang that serves rest endpoint to send sms. Repo include a web form built in react for user to interact with the app.

API provides a single endpoint to send sms

`/api/v1/sms/send` - POST

Payload

```
{
    "phone_number":"number",
    "texts":["text1","text2","text3"]
}
```
Success 

```{status : {"success" , "failure" , "success}}```

Failure

```{"message" : "phone number is invalid"}```

## Project Set up and Structure:

GO version 1.9 is used for building the backend api using dep as dependency manager.

Frontend is built using react + webpack

Parent project folder have been used for the golang API and client folder for the react app.

### GO Unit Tests
- Follows data table approach
- Consistent pattern using Args/Fields/Want format
- Assertion/Mocking using testify

### Config values
- Supplied through .env 

### Pre-Requisites:
- Git (just to clone the repo)
- Docker and Docker-compose

## Installation:
 Clone this repository
`https://github.com/nikhil-github/sms-app.git`

### Run Locally

- server -> `make run`
- client -> `make run-client`

http://localhost:3000 web form

### Run Docker

- server -> `make run-docker`
- client -> `make run-client-docker`

http://localhost:3000 web form


### Make targets

1. `make` - build the project
2. `make fmt` - format the codebase using `go fmt` and `goimports`
3. `make test` - run unit tests for the project


### API client

Simple client is added to the project that consumes the rest endpoint.
To run the client

```go run cmd/sms-app-client/main.go```
### Assumptions:
- API allows maximum of 160 characters per text.
- Secrets/Configs are supplied as env variables.
- Bitly go client library can be used to shorten URL.
- No automated retry incase of rate limited error response from transmit API.
