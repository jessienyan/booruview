#!/bin/bash
go install \
	-tags='no_postgres no_mysql no_ydb no_clickhouse no_libsql no_mssql no_vertica' \
	github.com/pressly/goose/v3/cmd/goose@v3.26.0
