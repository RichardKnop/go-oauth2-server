[1]: ../../../assets/blob/master/go-oauth2-server/login_screenshot.png
[2]: ../../../assets/blob/master/go-oauth2-server/authorization_code_screenshot.png
[3]: ../../../assets/blob/master/go-oauth2-server/implicit_screenshot.png
[4]: http://patreon_public_assets.s3.amazonaws.com/sized/becomeAPatronBanner.png
[5]: http://richardknop.com/images/btcaddress.png

## Go OAuth2 Server

This service implements [OAuth 2.0 specification](https://tools.ietf.org/html/rfc6749). Excerpts from the specification are included in this README file to describe different grant types. Please read the full spec for more detailed information.

[![Travis Status for RichardKnop/go-oauth2-server](https://travis-ci.org/RichardKnop/go-oauth2-server.svg?branch=master&label=linux+build)](https://travis-ci.org/RichardKnop/go-oauth2-server)
[![godoc for RichardKnop/go-oauth2-server](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/RichardKnop/go-oauth2-server)
[![codecov for RichardKnop/go-oauth2-server](https://codecov.io/gh/RichardKnop/go-oauth2-server/branch/master/graph/badge.svg)](https://codecov.io/gh/RichardKnop/go-oauth2-server)

[![Sourcegraph for RichardKnop/go-oauth2-server](https://sourcegraph.com/github.com/RichardKnop/go-oauth2-server/-/badge.svg)](https://sourcegraph.com/github.com/RichardKnop/go-oauth2-server?badge)
[![Donate Bitcoin](https://img.shields.io/badge/donate-bitcoin-orange.svg)](https://richardknop.github.io/donate/)

---

* [OAuth 2.0](#oauth-20)
  * [Client Authentication](#client-authentication)
  * [Grant Types](#grant-types)
    * [Authorization Code](#authorization-code)
    * [Implicit](#implicit)
    * [Resource Owner Password Credentials](#resource-owner-password-credentials)
    * [Client Credentials](#client-credentials)
  * [Refreshing An Access Token](#refreshing-an-access-token)
  * [Token Introspection](#token-introspection)
* [Plugins](#plugins)
* [Session Storage](#session-storage)
* [Dependencies](#dependencies)
* [Setup](#setup)
  * [etcd](#etcd)
  * [consul](#consul)
  * [postgres](#postgres)
* [Compile & Run Data](#compile--run)
* [Testing](#testing)
* [Docker](#docker)
* [Docker Compose](#docker-compose)
* [Supporting the project](#supporting-the-project)

## OAuth 2.0

### Client Authentication

http://tools.ietf.org/html/rfc6749#section-3.2.1

Clients must authenticate with client credentials (client ID and secret) when issuing requests to `/v1/oauth/tokens` endpoint. Basic HTTP authentication should be used.

### Grant Types

#### Authorization Code

http://tools.ietf.org/html/rfc6749#section-4.1

The authorization code grant type is used to obtain both access tokens and refresh tokens and is optimized for confidential clients. Since this is a redirection-based flow, the client must be capable of interacting with the resource owner's user-agent (typically a web browser) and capable of receiving incoming requests (via redirection) from the authorization server.

```
+----------+
| Resource |
|   Owner  |
|          |
+----------+
     ^
     |
    (B)
+----|-----+          Client Identifier      +---------------+
|         -+----(A)-- & Redirection URI ---->|               |
|  User-   |                                 | Authorization |
|  Agent  -+----(B)-- User authenticates --->|     Server    |
|          |                                 |               |
|         -+----(C)-- Authorization Code ---<|               |
+-|----|---+                                 +---------------+
  |    |                                         ^      v
 (A)  (C)                                        |      |
  |    |                                         |      |
  ^    v                                         |      |
+---------+                                      |      |
|         |>---(D)-- Authorization Code ---------'      |
|  Client |          & Redirection URI                  |
|         |                                             |
|         |<---(E)----- Access Token -------------------'
+---------+       (w/ Optional Refresh Token)
```

The client initiates the flow by directing the resource owner's user-agent to the authorization endpoint. The client includes its client identifier, requested scope, local state, and a redirection URI to which the authorization server will send the user-agent back once access is granted (or denied).

```
http://localhost:8080/web/authorize?client_id=test_client_1&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=code&state=somestate&scope=read_write
```

The authorization server authenticates the resource owner (via the user-agent).

![Log In page screenshot][1]

The authorization server then establishes whether the resource owner grants or denies the client's access request.

![Authorize page screenshot][2]

If the request fails due to a missing, invalid, or mismatching redirection URI, or if the client identifier is missing or invalid, the authorization server SHOULD inform the resource owner of the error and MUST NOT automatically redirect the user-agent to the invalid redirection URI.

If the resource owner denies the access request or if the request fails for reasons other than a missing or invalid redirection URI, the authorization server informs the client by adding the error parameter to the query component of the redirection URI.

```
https://www.example.com/?error=access_denied&state=somestate
```

Assuming the resource owner grants access, the authorization server redirects the user-agent back to the client using the redirection URI provided earlier (in the request or during client registration). The redirection URI includes an authorization code and any local state provided by the client earlier.

```
https://www.example.com/?code=7afb1c55-76e4-4c76-adb7-9d657cb47a27&state=somestate
```

The client requests an access token from the authorization server's token endpoint by including the authorization code received in the previous step. When making the request, the client authenticates with the authorization server. The client includes the redirection URI used to obtain the authorization code for verification.

```sh
curl --compressed -v localhost:8080/v1/oauth/tokens \
	-u test_client_1:test_secret \
	-d "grant_type=authorization_code" \
	-d "code=7afb1c55-76e4-4c76-adb7-9d657cb47a27" \
	-d "redirect_uri=https://www.example.com"
```
The authorization server authenticates the client, validates the authorization code, and ensures that the redirection URI received matches the URI used to redirect the client before. If valid, the authorization server responds back with an access token and, optionally, a refresh token.

```json
{
  "user_id": "1",
  "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "read_write",
  "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Implicit

http://tools.ietf.org/html/rfc6749#section-4.2

The implicit grant type is used to obtain access tokens (it does not support the issuance of refresh tokens) and is optimized for public clients known to operate a particular redirection URI. These clients are typically implemented in a browser using a scripting language such as JavaScript.

Since this is a redirection-based flow, the client must be capable of interacting with the resource owner's user-agent (typically a web browser) and capable of receiving incoming requests (via redirection) from the authorization server.

Unlike the authorization code grant type, in which the client makes separate requests for authorization and for an access token, the client receives the access token as the result of the authorization request.

The implicit grant type does not include client authentication, and relies on the presence of the resource owner and the registration of the redirection URI.  Because the access token is encoded into the redirection URI, it may be exposed to the resource owner and other applications residing on the same device.

```
+----------+
| Resource |
|  Owner   |
|          |
+----------+
     ^
     |
    (B)
+----|-----+          Client Identifier     +---------------+
|         -+----(A)-- & Redirection URI --->|               |
|  User-   |                                | Authorization |
|  Agent  -|----(B)-- User authenticates -->|     Server    |
|          |                                |               |
|          |<---(C)--- Redirection URI ----<|               |
|          |          with Access Token     +---------------+
|          |            in Fragment
|          |                                +---------------+
|          |----(D)--- Redirection URI ---->|   Web-Hosted  |
|          |          without Fragment      |     Client    |
|          |                                |    Resource   |
|     (F)  |<---(E)------- Script ---------<|               |
|          |                                +---------------+
+-|--------+
  |    |
 (A)  (G) Access Token
  |    |
  ^    v
+---------+
|         |
|  Client |
|         |
+---------+
```

The client initiates the flow by directing the resource owner's user-agent to the authorization endpoint. The client includes its client identifier, requested scope, local state, and a redirection URI to which the authorization server will send the user-agent back once access is granted (or denied).

```
http://localhost:8080/web/authorize?client_id=test_client_1&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=token&state=somestate&scope=read_write
```

The authorization server authenticates the resource owner (via the user-agent).

![Log In page screenshot][1]

The authorization server then establishes whether the resource owner grants or denies the client's access request.

![Authorize page screenshot][3]

If the request fails due to a missing, invalid, or mismatching redirection URI, or if the client identifier is missing or invalid, the authorization server SHOULD inform the resource owner of the error and MUST NOT automatically redirect the user-agent to the invalid redirection URI.

If the resource owner denies the access request or if the request fails for reasons other than a missing or invalid redirection URI, the authorization server informs the client by adding the following parameters to the fragment component of the redirection URI.

```
https://www.example.com/#error=access_denied&state=somestate
```

Assuming the resource owner grants access, the authorization server redirects the user-agent back to the client using the redirection URI provided earlier.  The redirection URI includes he access token in the URI fragment.

```
https://www.example.com/#access_token=087902d5-29e7-417b-a339-b57a60d6742a&expires_in=3600&scope=read_write&state=somestate&token_type=Bearer
```

The user-agent follows the redirection instructions by making a request to the web-hosted client resource (which does not include the fragment per [RFC2616]).  The user-agent retains the fragment information locally.

The web-hosted client resource returns a web page (typically an HTML document with an embedded script) capable of accessing the full redirection URI including the fragment retained by the user-agent, and extracting the access token (and other parameters) contained in the fragment.

The user-agent executes the script provided by the web-hosted client resource locally, which extracts the access token.

The user-agent passes the access token to the client.

#### Resource Owner Password Credentials

http://tools.ietf.org/html/rfc6749#section-4.3

The resource owner password credentials grant type is suitable in cases where the resource owner has a trust relationship with the client, such as the device operating system or a highly privileged application. The authorization server should take special care when enabling this grant type and only allow it when other flows are not viable.

This grant type is suitable for clients capable of obtaining the resource owner's credentials (username and password, typically using an interactive form). It is also used to migrate existing clients using direct authentication schemes such as HTTP Basic or Digest authentication to OAuth by converting the stored credentials to an access token.

```
+----------+
| Resource |
|  Owner   |
|          |
+----------+
     v
     |    Resource Owner
     (A) Password Credentials
     |
     v
+---------+                                  +---------------+
|         |>--(B)---- Resource Owner ------->|               |
|         |         Password Credentials     | Authorization |
| Client  |                                  |     Server    |
|         |<--(C)---- Access Token ---------<|               |
|         |    (w/ Optional Refresh Token)   |               |
+---------+                                  +---------------+

```

The resource owner provides the client with its username and password.

The client requests an access token from the authorization server's token endpoint by including the credentials received from the resource owner. When making the request, the client authenticates with the authorization server.

```sh
curl --compressed -v localhost:8080/v1/oauth/tokens \
	-u test_client_1:test_secret \
	-d "grant_type=password" \
	-d "username=test@user" \
	-d "password=test_password" \
	-d "scope=read_write"
```

The authorization server authenticates the client and validates the resource owner credentials, and if valid, issues an access token.

```json
{
  "user_id": "1",
  "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "read_write",
  "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

#### Client Credentials

http://tools.ietf.org/html/rfc6749#section-4.4

The client can request an access token using only its client credentials (or other supported means of authentication) when the client is requesting access to the protected resources under its control, or those of another resource owner that have been previously arranged with the authorization server (the method of which is beyond the scope of this specification).

The client credentials grant type MUST only be used by confidential clients.

```
+---------+                                  +---------------+
|         |                                  |               |
|         |>--(A)- Client Authentication --->| Authorization |
| Client  |                                  |     Server    |
|         |<--(B)---- Access Token ---------<|               |
|         |                                  |               |
+---------+                                  +---------------+
```

The client authenticates with the authorization server and requests an access token from the token endpoint.

```sh
curl --compressed -v localhost:8080/v1/oauth/tokens \
	-u test_client_1:test_secret \
	-d "grant_type=client_credentials" \
	-d "scope=read_write"
```

The authorization server authenticates the client, and if valid, issues an access token.

```json
{
  "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "read_write",
  "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```

### Refreshing An Access Token

http://tools.ietf.org/html/rfc6749#section-6

If the authorization server issued a refresh token to the client, the client can make a refresh request to the token endpoint in order to refresh the access token.

```sh
curl --compressed -v localhost:8080/v1/oauth/tokens \
	-u test_client_1:test_secret \
	-d "grant_type=refresh_token" \
	-d "refresh_token=6fd8d272-375a-4d8a-8d0f-43367dc8b791"
```

The authorization server MUST:

* require client authentication for confidential clients or for any client that was issued client credentials (or with other authentication requirements),

* authenticate the client if client authentication is included and ensure that the refresh token was issued to the authenticated client, and

* validate the refresh token.

If valid and authorized, the authorization server issues an access token.

```json
{
  "user_id": "1",
  "access_token": "1f962bd5-7890-435d-b619-584b6aa32e6c",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "read_write",
  "refresh_token": "3a6b45b8-9d29-4cba-8a1b-0093e8a2b933"
}
```

The authorization server MAY issue a new refresh token, in which case the client MUST discard the old refresh token and replace it with the new refresh token.  The authorization server MAY revoke the old refresh token after issuing a new refresh token to the client.  If a new refresh token is issued, the refresh token scope MUST be identical to that of the refresh token included by the client in the request.

### Token Introspection

https://tools.ietf.org/html/rfc7662

If the authorization server issued a access token or refresh token to the client, the client can make a request to the introspect endpoint in order to learn meta-information about a token.

```sh
curl --compressed -v localhost:8080/v1/oauth/introspect \
	-u test_client_1:test_secret \
	-d "token=00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d "token_type_hint=access_token"
```

The authorization server responds meta-information about a token.

```json
{
  "active": true,
  "scope": "read_write",
  "client_id": "test_client_1",
  "username": "test@username",
  "token_type": "Bearer",
  "exp": 1454868090
}
```

## Plugins

This server is easily extended or modified through the use of plugins. Four services, [health](https://github.com/RichardKnop/go-oauth2-server/tree/master/health), [oauth](https://github.com/RichardKnop/go-oauth2-server/tree/master/oauth), [session](https://github.com/RichardKnop/go-oauth2-server/tree/master/session) and [web](https://github.com/RichardKnop/go-oauth2-server/tree/master/web) are available for modification.

In order to implement a plugin:
1. Create your own interface that implements all of methods of the service you are replacing.
2. Modify `cmd/run_server.go` to use your service by calling the `session.Use[service-you-are-replaceing]Service(yourCustomService.NewService())` before the services are initialized via `services.Init(cnf, db)`.

For example, to implement an available [redis session storage plugin](https://github.com/adam-hanna/redis-sessions):

~~~go
// $ go get https://github.com/adam-hanna/redis-sessions
//
// cmd/run_server.go
import (
    ...
    "github.com/adam-hanna/redis-sessions/redis"
    ...
)

// RunServer runs the app
func RunServer(configBackend string) error {
    ...

    // configure redis for session store
    sessionSecrets := make([][]byte, 1)
    sessionSecrets[0] = []byte(cnf.Session.Secret)
    redisConfig := redis.ConfigType{
        Size:           10,
        Network:        "tcp",
        Address:        ":6379",
        Password:       "",
        SessionSecrets: sessionSecrets,
    }

    // start the services
    services.UseSessionService(redis.NewService(cnf, redisConfig))
    if err := services.InitServices(cnf, db); err != nil {
        return err
    }
    defer services.CloseServices()

    ...
}
~~~

## Session Storage

By default, this server implements in-memory, cookie sessions via [gorilla sessions](https://github.com/gorilla/sessions).

However, because the session service can be replaced via a plugin, any of the available [gorilla sessions store implementations](https://github.com/gorilla/sessions#store-implementations) can be wrapped by `session.ServiceInterface`.

## Dependencies

Since Go 1.11, a new recommended dependency management system is via [modules](https://github.com/golang/go/wiki/Modules).

This is one of slight weaknesses of Go as dependency management is not a solved problem. Previously Go was officially recommending to use the [dep tool](https://github.com/golang/dep) but that has been abandoned now in favor of modules.

## Setup

For distributed config storage you can use either etcd or consul (etcd being the default)

If you are developing on OSX, install `etcd` or `consul`, `Postgres` and `nats-streaming-server`:

### etcd

```sh
brew install etcd
```

Load a development configuration into `etcd`:

```sh
ETCDCTL_API=3 etcdctl put /config/go_oauth2_server.json '{
  "Database": {
    "Type": "postgres",
    "Host": "localhost",
    "Port": 5432,
    "User": "go_oauth2_server",
    "Password": "",
    "DatabaseName": "go_oauth2_server",
    "MaxIdleConns": 5,
    "MaxOpenConns": 5
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
  "IsDevelopment": true
}'
```

If you are using etcd API version 3, use `etcdctl put` instead of `etcdctl set`.

Check the config was loaded properly:

```sh
ETCDCTL_API=3 etcdctl get /config/go_oauth2_server.json
```

### consul

```sh
brew install consul
```

Load a development configuration into `consul`:

```sh
consul kv put /config/go_oauth2_server.json '{
  "Database": {
    "Type": "postgres",
    "Host": "localhost",
    "Port": 5432,
    "User": "go_oauth2_server",
    "Password": "",
    "DatabaseName": "go_oauth2_server",
    "MaxIdleConns": 5,
    "MaxOpenConns": 5
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
  "IsDevelopment": true
}'
```

Check the config was loaded properly:

```sh
consul kv get /config/go_oauth2_server.json
```

### Postgres

```sh
brew install postgres
```

You might want to create a `Postgres` database:

```sh
createuser --createdb go_oauth2_server
createdb -U go_oauth2_server go_oauth2_server
```

## Compile & Run

Compile the app:

```sh
go install .
```

The binary accepts an optional flag of `--configBackend` which can be set to `etcd | consul`, defaults to `etcd`

Run migrations:

```sh
go-oauth2-server migrate
```

And finally, run the app:

```sh
go-oauth2-server runserver
```

When deploying, you can set etcd related environment variables:

* `ETCD_ENDPOINTS`
* `ETCD_CERT_FILE`
* `ETCD_KEY_FILE`
* `ETCD_CA_FILE`
* `ETCD_CONFIG_PATH`

You can also set consul related variables

* `CONSUL_ENDPOINT`
* `CONSUL_CERT_FILE`
* `CONSUL_KEY_FILE`
* `CONSUL_CA_FILE`
* `CONSUL_CONFIG_PATH`

and the equivalent above commands would be

```sh
go-oauth2-server --configBackend consul migrate
```
```sh
go-oauth2-server --configBackend consul runserver
```

## Testing

I have used a mix of unit and functional tests so you need to have `sqlite` installed in order for the tests to run successfully as the suite creates an in-memory database.

To run tests:

```sh
make test
```

## Docker

Build a Docker image and run the app in a container:

```sh
docker build -t go-oauth2-server:latest .
docker run -e ETCD_ENDPOINTS=localhost:2379 -p 8080:8080 --name go-oauth2-server go-oauth2-server:latest
```

You can load fixtures with `docker exec` command:

```sh
docker exec <container_id> /go/bin/go-oauth2-server loaddata \
  oauth/fixtures/scopes.yml \
  oauth/fixtures/roles.yml \
  oauth/fixtures/test_clients.yml
```

## Docker Compose

You can use [docker-compose](https://docs.docker.com/compose/) to start the app, postgres, etcd in separate linked containers:

```sh
docker-compose up
```

During `docker-compose up` process all configuration and fixtures will be loaded. After successful up you can check, that app is running using for example the health check request:

```sh
curl --compressed -v localhost:8080/v1/health
```

## Supporting the project

Donate BTC to my wallet if you find this project useful: `12iFVjQ5n3Qdmiai4Mp9EG93NSvDipyRKV`

![Donate BTC][5]
