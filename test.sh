#!/bin/sh

go run ./httpc/*.go get -v 'http://httpbin.org/get?course=networking&assignment=1'

go run ./httpc/*.go get -v 'http://postman-echo.com/get?foo=bar'

go run ./httpc/*.go post -v 'http://postman-echo.com/post?foo=bar'
