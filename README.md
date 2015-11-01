# Go Microservice Example

A simple Go microservice example.

This example service implements resource owner password credentials grant of [OAuth 2.0 specification](http://tools.ietf.org/html/rfc6749#section-4.3). It can be used as a simple login service.

Goals:

* reusable
* test database
* migrations
* easily packaged as Docker container

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

A test database is created before running the tests and destroyed after the tests have run.

To run tests:

```
make test
```
