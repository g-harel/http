#!/bin/sh

go run ./httpc/*.go get 'http://httpbin.org/get?course=networking&assignment=1'

go run ./httpc/*.go get -v 'http://httpbin.org/get?course=networking&assignment=1'

go run ./httpc/*.go post -h 'Content-Type: application/json' -d '{"Assignment": 1}' http://httpbin.org/post
