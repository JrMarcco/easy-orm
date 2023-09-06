cover_test:
	go test -v -cover

integration_test:
	go test -race -tags=integration
