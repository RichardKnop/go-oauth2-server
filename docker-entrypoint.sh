#!/bin/bash

# 1. Run database migrations
/go/bin/go-oauth2-server migrate

# 2. Load fixtures
/go/bin/go-oauth2-server loaddata \
  oauth/fixtures/scopes.yml

# Finally, run the server
/go/bin/go-oauth2-server runserver
