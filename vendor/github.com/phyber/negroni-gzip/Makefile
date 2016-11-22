PKGNAME=github.com/phyber/negroni-gzip/gzip
GOCMD=go
GOTOOL=$(GOCMD) tool
GOTEST=$(GOCMD) test
COVERFILE=cover.out
TAGSFILE=tags
TESTCOVER=$(GOTEST) -coverprofile $(COVERFILE)
GOCOVER=$(GOTOOL) cover -func=$(COVERFILE)

clean:
	rm -f $(COVERFILE) gzip/$(COVERFILE) gzip/$(TAGSFILE)

test:
	$(GOTEST) $(PKGNAME)

testcover:
	$(TESTCOVER) $(PKGNAME)
	$(GOCOVER)
