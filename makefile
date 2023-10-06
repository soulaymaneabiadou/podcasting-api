define setup_env
	$(eval include .env.$(1))
	$(eval export)
endef

up:
	docker compose -f compose.dev.yml up -d --build --force-recreate

migrate:
	$(call setup_env,local)
	go run . db:migrate

seed:
	$(call setup_env,local)
	go run . db:seed

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

stripe-listen:
	stripe listen --forward-to localhost:5000/api/v1/webhooks/stripe --forward-connect-to localhost:5000/api/v1/webhooks/stripe --latest

stripe-trigger:
	stripe trigger --api-version "2023-08-16" customer.subscription.created
