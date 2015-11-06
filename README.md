[![Build Status](https://travis-ci.org/RichardKnop/go-oauth2-server.svg?branch=master "Build Status")](https://travis-ci.org/RichardKnop/go-oauth2-server)

# Go OAuth2 Server

This service implements [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3).

# Index

* [Go OAuth2 Server](#go-oauth2-server)
* [Index](#index)
* [API](#api)
  * [OAuth 2.0](#oauth-20)
    * [Client Authentication](#client-authorization)
    * [Grant Types](#grant-types)
      * [Authorization Code](#authorization-code)
      * [Implicit](#implicit)
      * [User Credentials](#user-credentials)
      * [Client Credentials](#client-credentials)
      * [Refreshing Token](#refreshing-token)
    * [Scope](#scope)
  * [Users](#users)
    * [Register](#register)
* [Development](#development)
  * [Third Party Libraries](#third-party-libraries)
  * [Dependencies](#dependencies)
  * [Configuration](#configuration)
  * [Database](#database)
  * [Testing](#Testing)

# API

## OAuth 2.0

### Client Authentication

http://tools.ietf.org/html/rfc6749#section-3.2.1

Clients must authenticate with client credentials (client ID and secret) when issuing requests to `/oauth2/api/v1/tokens` endpoint. Basic HTTP authentication should be used.

### Grant Types

#### Authorization Code

http://tools.ietf.org/html/rfc6749#section-4.1

TODO

#### Implicit

http://tools.ietf.org/html/rfc6749#section-4.2

TODO

#### User Credentials

http://tools.ietf.org/html/rfc6749#section-4.3

Given you have a username and password, you can get a new access token.

```
$ curl -u test_client_id:test_client_secret localhost:8080/oauth2/api/v1/tokens
  -d 'grant_type=password' \
  -d "username=test_username" \
  -d "password=test_password' \
  -d "scope=read_write"
```

Response:

```json
{
    "id": 1,
    "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "read_write",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Client Credentials

http://tools.ietf.org/html/rfc6749#section-4.4

Given you have a client ID and secret, you can get a new access token.

```
$ curl -u test_client_id:test_client_password localhost:8080/oauth2/api/v1/tokens \
  -d 'grant_type=client_credentials' \
  -d "scope=read_write"
```

Response:

```json
{
    "id": 1,
    "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "read_write",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

### Refreshing Token

http://tools.ietf.org/html/rfc6749#section-6

Let's say you have created a new access token using client or user credentials grant type. The response included a refresh token which you can use to get a new access token before your current access token expires.

```
$ curl -u test_client_id:test_client_password localhost:8080/oauth2/api/v1/tokens \
  -d 'grant_type=refresh_token" \
  -d "refresh_token=6fd8d272-375a-4d8a-8d0f-43367dc8b791'
```

Response:

```json
{
    "id": 1,
    "access_token": "1f962bd5-7890-435d-b619-584b6aa32e6c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "read_write",
    "refresh_token": "3a6b45b8-9d29-4cba-8a1b-0093e8a2b933"
}
```

### Scope

http://tools.ietf.org/html/rfc6749#section-3.3

Scope is a space delimited case-sensitive string where each part defines a specific access range. It can be used for ACL.

You can define your scopes and insert them into scopes table, is_default flag can be used to specify default scope.

## Users

### Register

To register a new user:

```
curl localhost:8080/oauth2/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testusername",
    "password": "password",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

Response:

```json
{
  "first_name": "John",
  "id": 1,
  "last_name": "Doe",
  "username": "testusername"
}
```

# Development

## Third Party Libraries

Few notable third party libraries used:

* [godep](https://github.com/tools/godep) - build packages reproducibly by fixing their dependencies
* [Viper](https://github.com/spf13/viper) - Go configuration with fangs
* [gorm](https://github.com/jinzhu/gorm) - The fantastic ORM library for Golang, aims to be developer friendly
* [go-json-rest](https://github.com/ant0ine/go-json-rest) - A quick and easy way to setup a RESTful JSON API
* [testify](https://github.com/stretchr/testify) - A sacred extension to the standard go testing package

## Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in a vendor directory. This approach is called "vendoring" and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

To update dependencies:

```
make update-deps
```

To install dependencies:

```
make install-deps
```

## Configuration

The configuration is done via environment variables. Available variables:

* `OAUTH2_DATABASE_TYPE` (defaults to `postgres`)
* `OAUTH2_DATABASE_HOST` (defaults to `127.0.0.1`)
* `OAUTH2_DATABASE_PORT` (defaults to `5432`)
* `OAUTH2_DATABASE_USER` (defaults to `go_oauth2_server`)
* `OAUTH2_DATABASE_PASSWORD` (defaults to empty string)
* `OAUTH2_DATABASE_NAME` (defaults to `go_oauth2_server`)
* `OAUTH2_ACCESS_TOKEN_LIFETIME` (defaults to `3600` or 1 hour)
* `OAUTH2_REFRESH_TOKEN_LIFETIME` (defaults to `1209600` or 14 days)

Variables are prefixed with `OAUTH2_`.

## Database

In order to run this service, create a Postgres user and database:

```
createuser --createdb go_oauth2_server
createdb -U go_microservice_example go_oauth2_server
```

Set environment variables to match your database connection details.

## Testing

To run tests:

```
make test
```
