#!/bin/bash

# 1. Run database migrations
/go/bin/go-oauth2-server migrate

# Finally, run the server
/go/bin/go-oauth2-server runserver
