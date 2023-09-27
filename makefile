define setup_env
	$(eval include .env.$(1))
	$(eval export)
endef

up:
	docker compose -f compose.dev.yml up -d --build --force-recreate

dev:
# go install github.com/cosmtrek/air@latest
	$(call setup_env,local)
	rm -rf tmp/
	air

test:
	$(call setup_env,test)
	go test ./...

run:
	$(call setup_env,prod)
	go run .

build:
	rm ./podcast
	go build -o ./podcast

prod: build
	$(call setup_env,prod)
	./podcast

clean:
	rm ./podcast
	rm -rf tmp/

gen:
# go install github.com/google/wire/cmd/wire@latest
# wire gen ./middleware ./services ./routes
# wire gen ./...
	wire ./...

install-deps:
	go install github.com/cosmtrek/air@latest
	go install github.com/google/wire/cmd/wire@latest