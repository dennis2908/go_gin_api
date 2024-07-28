set -m

cd rabbitsInsertCustomers && \ go run InsertCustomers.go &

cd rabbitsEmailCustomers && \ go run EmailCustomers.go

fg %1