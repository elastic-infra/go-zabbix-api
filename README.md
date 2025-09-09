# Go zabbix api

[![GoDoc](https://godoc.org/github.com/claranet/go-zabbix-api?status.svg)](https://godoc.org/github.com/claranet/go-zabbix-api) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) [![Build Status](https://travis-ci.org/claranet/go-zabbix-api.svg?branch=master)](https://travis-ci.org/claranet/go-zabbix-api)

This Go package provides access to Zabbix API.

Tested on Zabbix 6.0, 6.2, 6.4.

This package aims to support multiple zabbix resources from its API like trigger, application, host group, host, item, template..

## Install

Install it: `go get github.com/elastic-infra/go-zabbix-api`

## Getting started

```go
package main

import (
	"fmt"

	"github.com/elastic-infra/go-zabbix-api"
)

func main() {
	user := "MyZabbixUsername"
	pass := "MyZabbixPassword"
	api := zabbix.NewAPI("http://localhost/api_jsonrpc.php")
	api.Login(user, pass)

	res, err := api.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to zabbix api v%s\n", res)
}
```

## Tests

### Run tests

```bash
export TEST_ZABBIX_URL=http://localhost:8080/zabbix/api_jsonrpc.php
export TEST_ZABBIX_USER=Admin
export TEST_ZABBIX_PASSWORD=zabbix
export TEST_ZABBIX_VERBOSE=1
go test -v
```

`TEST_ZABBIX_URL` may contain HTTP basic auth username and password: `http://username:password@host/api_jsonrpc.php`. Also, in some setups URL should be like `http://host/zabbix/api_jsonrpc.php`.

## References

Documentation is available on [godoc.org](https://godoc.org/github.com/claranet/go-zabbix-api).
Also, Rafael Fernandes dos Santos wrote a [great article](http://www.sourcecode.net.br/2014/02/zabbix-api-with-golang.html) about using and extending this package.

License: Simplified BSD License (see [LICENSE](LICENSE)).
