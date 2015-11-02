# Go Microservice Example

A simple Go microservice example.

This example service implements resource owner password credentials grant of [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3). It can be used as a simple login service.

Goals:

* reusable
* test database
* migrations
* easily packaged as Docker container

- [Third Party Libraries](#third-party-libraries)

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

A test database is created before running the tests and destroyed after the tests have run. Therefor, You need to have Postgres installed in order for tests to run.

To run tests:

```
make test
```
