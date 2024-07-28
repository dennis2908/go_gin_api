#!/bin/sh

set -m

nohup sh -c cd rabbitsInsertCustomers && go run InsertCustomers.go &

nohup sh -c cd rabbitsEmailCustomers && go run EmailCustomers.go

fg %1