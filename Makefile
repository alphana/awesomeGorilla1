install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	swagger generate spec -i ./swagger.yml -o ./swagger.json
