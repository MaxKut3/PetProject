.DEFAULT_GOAL = run

run:
	gofmt -w .
	goimports -w .
	go run cmd/PetProject/main.go

mig-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=postgres dbname=postgres sslmode=disable host=localhost port=5434" goose -dir migrations up