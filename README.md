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
  * [Users](#users)
    * [Register](#register)
* [Development](#development)
  * [Dependencies](#dependencies)
  * [Setup](#setup)
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
$ curl localhost:8080/oauth2/api/v1/tokens \
  -u test_client_id:test_client_password \
  -d "grant_type=password" \
  -d "username=test_username" \
  -d "password=test_password" \
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
$ curl localhost:8080/oauth2/api/v1/tokens \
  -u test_client_id:test_client_password \
  -d "grant_type=client_credentials" \
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

Let's say you have obtained an access token previously. The response included a refresh token which you can use to get a new access token before your current access token expires.

```
$ curl localhost:8080/oauth2/api/v1/tokens \
  -u test_client_id:test_client_password \
  -d "grant_type=refresh_token" \
  -d "refresh_token=6fd8d272-375a-4d8a-8d0f-43367dc8b791"
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

## Users

### Register

To register a new user:

```
$ curl localhost:8080/oauth2/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_username",
    "password": "test_password",
  }'
```

Response:

```json
{
  "id": 1,
  "username": "test_username"
}
```

# Development
## Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in a vendor directory. This approach is called "vendoring" and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

To update dependencies during development:

```
$ make update-deps
```

To install dependencies:

```
$ make install-deps
```

## Setup

If you are developing on OSX, install `etcd`, `Postgres`:

```
$ brew install etcd
$ brew install postgres
```

You might want to create a `Postgres` database:

```
$ createuser --createdb go_oauth2_server
$ createdb -U go_microservice_example go_oauth2_server
```

Load a configuration into `etcd`:

```
$ curl -L http://127.0.0.1:4001/v2/keys/config/go_oauth2_server.json -XPUT -d value='{
  "Database": {
    "Type": "postgres",
    "Host": "127.0.0.1",
    "Port": 5432,
    "User": "go_oauth2_server",
    "Password": "",
    "DatabaseName": "go_oauth2_server"
  },
  "AccessTokenLifetime": 3600,
  "RefreshTokenLifetime": 1209600
}'
```

## Testing

To run tests:

```
make test
```
