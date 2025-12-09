PHONY: lint e2e

lint:
	go run lint.go

e2e:
	docker-compose -f test/compose.test.yml up -d
	go test -v -count=1 ./test
	docker-compose -f test/compose.test.yml down -v