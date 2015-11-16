[![Build Status](https://travis-ci.org/RichardKnop/go-oauth2-server.svg?branch=master "Build Status")](https://travis-ci.org/RichardKnop/go-oauth2-server)

# Go OAuth2 Server

This service implements [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3).

# Index

* [Go OAuth2 Server](#go-oauth2-server)
* [Index](#index)
* [API](#api)
  * [OAuth 2.0](#oauth-20)
    * [Client Authentication](#client-authentication)
    * [Grant Types](#grant-types)
      * [Authorization Code](#authorization-code)
      * [Implicit](#implicit)
      * [User Credentials](#user-credentials)
      * [Client Credentials](#client-credentials)
      * [Refresh Token](#refresh-token)
* [Development](#development)
  * [Dependencies](#dependencies)
  * [Setup](#setup)
  * [Test Data](#test-data)
  * [Testing](#Testing)

# API

## OAuth 2.0

### Client Authentication

http://tools.ietf.org/html/rfc6749#section-3.2.1

Clients must authenticate with client credentials (client ID and secret) when issuing requests to `/oauth/api/v1/tokens` endpoint. Basic HTTP authentication should be used.

### Grant Types

#### Authorization Code

http://tools.ietf.org/html/rfc6749#section-4.1

First, you need to send a user to the authorization page with proper parameters, e.g.:

```
http://localhost:8080/web/authorize?client_id=test_client&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=code&state=somestate
```

The user will have to register and log in he/she you hasn't done so yet:

![Log In page screenshot](https://raw.githubusercontent.com/RichardKnop/assets/master/go-oauth2-server/login_screenshot.png)

After logging in, the user will be presented with an authorization page where he/she can accept or decline to authorize the client to act on his/her behalf:

![Authorize page screenshot](https://raw.githubusercontent.com/RichardKnop/assets/master/go-oauth2-server/authorize_screenshot.png)

If the user declines, the user agent will be redirected to the `redirect_uri` with error parameter in the query string, e.g.:

```
https://www.example.com/?error=access_denied&state=somestate
```

Given the user accepts, the user agent will be redirected to the `redirect_uri` and the authorization code will be in the query string, e.g.:

```
https://www.example.com/?code=7afb1c55-76e4-4c76-adb7-9d657cb47a27&state=somestate
```

Once you have an authorization code, you can exchange it for an access token:

```
$ curl -v localhost:8080/oauth/api/v1/tokens \
  -u test_client:test_secret \
  -d "grant_type=authorization_code" \
  -d "code=7afb1c55-76e4-4c76-adb7-9d657cb47a27"
```

Response:

```json
{
    "id": 1,
    "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "foo bar",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Implicit

http://tools.ietf.org/html/rfc6749#section-4.2

Very similar to the authorization code but an access token is returned in URL fragment without a need to make an additional API request to exchange the authorization code for an access token.

So the user would be sent to the authorize page with `response_type` parameter set to `token`, e.g.:

```
http://localhost:8080/web/authorize?client_id=test_client&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=token&state=somestate
```

If the user declines, the user agent will be redirected to the `redirect_uri` with error parameter in the query string, e.g.:

```
https://www.example.com/?error=access_denied&state=somestate
```

Given the user accepts, the user agent will be redirected to the `redirect_uri` and the access token will be directly in the URL:

```
https://www.example.com/?access_token=087902d5-29e7-417b-a339-b57a60d6742a&expires_in=3600&refresh_token=6531b6ae-6db4-4aa6-934f-486807c697f2&state=somestate&token_type=Bearer
```

#### User Credentials

http://tools.ietf.org/html/rfc6749#section-4.3

Given you have a username and password, you can get a new access token.

```
$ curl -v localhost:8080/oauth/api/v1/tokens \
  -u test_client:test_secret \
  -d "grant_type=password" \
  -d "username=test@username" \
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
    "scope": "foo bar",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Client Credentials

http://tools.ietf.org/html/rfc6749#section-4.4

Given you have a client ID and secret, you can get a new access token.

```
$ curl -v localhost:8080/oauth/api/v1/tokens \
  -u test_client:test_secret \
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
    "scope": "foo bar",
    "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Refresh Token

http://tools.ietf.org/html/rfc6749#section-6

Let's say you have obtained an access token previously. The response included a refresh token which you can use to get a new access token before your current access token expires.

```
$ curl -v localhost:8080/oauth/api/v1/tokens \
  -u test_client:test_secret \
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
    "scope": "foo bar",
    "refresh_token": "3a6b45b8-9d29-4cba-8a1b-0093e8a2b933"
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
  "Oauth": {
    "AccessTokenLifetime": 3600,
    "RefreshTokenLifetime": 1209600,
    "AuthCodeLifetime": 3600  
  },
  "Session": {
    "Secret": "test_secret",
    "Path": "/",
    "MaxAge": 604800,
    "HTTPOnly": true
  },
  "TrustedClient": {
    "ClientID": "test_client",
    "Secret": "test_secret"
  }
}'
```

Run migrations:

```
$ go run main.go migrate
```

And finally, run the app:

```
$ go run main.go runserver
```

## Test Data

You might want to insert some test data if you are testing locally using curl:

```sql
insert into scopes(scope, is_default) values('read', true);
insert into scopes(scope, is_default) values('read_write', false);

insert into clients(client_id, secret) values('test_client', '$2a$10$CUoGytf1pR7CC6Y043gt/.vFJUV4IRqvH5R6F0VfITP8s2TqrQ.4e');

insert into users(username, password) values('test@username', '$2a$10$4J4t9xuWhOKhfjN0bOKNReS9sL3BVSN9zxIr2.VaWWQfRBWh1dQIS');
```

## Testing

Some of the tests are functional. You need to have `sqlite` and `etcd` installed and running in order to run the tests.

To run tests:

```
make test
```
