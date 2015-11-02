* [Introduction](#introduction)
  * [Third Party Libraries](#third-party-libraries)
  * [Dependencies](#dependencies)
  * [Database](#database)
  * [Testing](#Testing)
* [API](#api)
  * [Grant Types](#grant-types)
    * [User Credentials](#user-credentials)

# Introduction

A simple Go microservice example.

This example service implements resource owner password credentials grant of [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3). It can be used as a simple login service.

Goals:

* reusable
* test database
* migrations
* easily packaged as Docker container

## Third Party Libraries

Few notable third party libraries used:

* [godep](https://github.com/tools/godep) - build packages reproducibly by fixing their dependencies
* [Viper](https://github.com/spf13/viper) - Go configuration with fangs
* [gorm](https://github.com/jinzhu/gorm) - The fantastic ORM library for Golang, aims to be developer friendly

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

## Database

Create a Postgres user and database:

```
createuser --createdb go_microservice_example
createdb -U go_microservice_example go_microservice_example
```

## Testing

A test database is created before running the tests and destroyed after the tests have run. Therefor, You need to have Postgres installed in order for tests to run (i.e. `createdb` and `createuser` commands must be available).

To run tests:

```
make test
```

# API

## Grant Types

### User Credentials

Given you have a username and password, you can get a new access token:

```
$ curl -u testclient:testpassword localhost:8080/api/v1/tokens/ -d 'grant_type=password&username=testuser@example.com&password=testpassword'
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
