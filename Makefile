DEPS=go list -f '{{range .TestImports}}{{.}} {{end}}' ./...

export GO15VENDOREXPERIMENT=1

update-deps:
	rm -rf Godeps
	rm -rf vendor
	go get github.com/tools/godep
	godep save ./...

install-deps:
	go get github.com/tools/godep
	godep restore
	$(DEPS) | xargs -n1 go get -d

test-set-up:
	createuser --createdb go_microservice_example_test
	createdb -U go_microservice_example_test go_microservice_example_test

run-tests:
	DATABASE_USER=go_microservice_example_test DATABASE_NAME=go_microservice_example_test go list ./... | grep -v vendor | xargs -n1 go test -timeout=3s

test-tear-down:
	dropdb -U go_microservice_example_test go_microservice_example_test
	dropuser go_microservice_example_test

test: test-set-up run-tests test-tear-down
