[![Build Status](https://travis-ci.org/RichardKnop/go-oauth2-server.svg?branch=master "Build Status")](https://travis-ci.org/RichardKnop/go-oauth2-server)

# Go OAuth2 Server

This service implements [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3).

# Index

* [Go OAuth2 Server](#go-oauth2-server)
* [Index](#index)
* [API](#api)
  * [OAuth 2.0](#oauth-20)
    * [Grant Types](#grant-types)
      * [Authorization Code](#authorization-code)
      * [Implicit](#implicit)
      * [User Credentials](#user-credentials)
      * [Client Credentials](#client-credentials)
    * [Refreshing Token](#refreshing-token)
    * [Scope](#scope)
    * [Authentication](#authorization)
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

Either by using HTTP Basic Authentication:

```
$ curl -u testusername:testpassword localhost:8080/api/v1/tokens
  -d 'grant_type=password'
```

Or using POST body:

```
$ curl localhost:8000/api/v1/tokens \
  -d 'grant_type=password" \
  -d `"username=testusername" \
  -d "password=testpassword'
```

You should get a response like:

```json
{
    "id": 1,
    "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "foo bar qux",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Client Credentials

http://tools.ietf.org/html/rfc6749#section-4.4

Given you have a client ID and secret, you can get a new access token.

Either by using HTTP Basic Authentication:

```
$ curl -u testclient:testpassword localhost:8080/api/v1/tokens \
  -d 'grant_type=client_credentials'
```

Or using POST body:

```
$ curl localhost:8000/api/v1/tokens \
  -d 'grant_type=client_credentials" \
  -d "client_id=testclient" \
  -d "client_secret=testpassword'
```

You should get a response like:

```json
{
    "id": 1,
    "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "foo bar qux",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

### Refreshing Token

http://tools.ietf.org/html/rfc6749#section-6

Let's say you have created a new access token using client or user credentials grant type. The response included a refresh token which you can use to get a new access token before your current access token expires.

```
$ curl localhost:8080/api/v1/tokens \
  -d 'grant_type=refresh_token" \
  -d "refresh_token=6fd8d272-375a-4d8a-8d0f-43367dc8b791'
```

And you get a new access token:

```json
{
    "id": 1,
    "access_token": "1f962bd5-7890-435d-b619-584b6aa32e6c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "foo bar qux",
    "refresh_token": "3a6b45b8-9d29-4cba-8a1b-0093e8a2b933"
}
```

### Scope

http://tools.ietf.org/html/rfc6749#section-3.3

Scope is quite arbitrary. Basically it is a space delimited case-sensitive string where each part defines a specific access range.

You can define your scopes and insert them into scopes table, is_default flag can be used to specify default scope.

### Authentication

TODO

## Users

### Register

To register a new user, POST to `/api/v1/users`:

```
curl localhost:8080/api/v1/users \
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

* `DATABASE_TYPE` (defaults to `postgres`)
* `DATABASE_HOST` (defaults to `127.0.0.1`)
* `DATABASE_PORT` (defaults to `5432`)
* `DATABASE_USER` (defaults to `go_oauth2_server`)
* `DATABASE_PASSWORD` (defaults to empty string)
* `DATABASE_NAME` (defaults to `go_oauth2_server`)
* `ACCESS_TOKEN_LIFETIME` (defaults to `3600` or 1 hour)
* `REFRESH_TOKEN_LIFETIME` (defaults to `1209600` or 14 days)

Variables are not prefixed as this service is intended to run inside a Docker container so there should be no conflict with some other service's configuration.

## Database

In order to run this service, create a Postgres user and database:

```
createuser --createdb go_microservice_example
createdb -U go_microservice_example go_microservice_example
```

Set environment variables to match your database connection details.

## Testing

To run tests:

```
make test
```
