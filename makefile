dev:
# go install github.com/cosmtrek/air@latest
	PORT=5000 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=password \
	POSTGRES_DB=podcast \
	JWT_ACCESS_SECRET=test \
	JWT_REFRESH_SECRET=testrefresh \
	JWT_ACCESS_EXPIRE=5 \
	JWT_ACCESS_COOKIE_EXPIRE=5 \
	JWT_REFRESH_EXPIRE=7 \
	JWT_REFRESH_COOKIE_EXPIRE=7 \
	SMTP_IDENTITY="Podcast Platform" \
	SMTP_USERNAME=221546739a5c13 \
	SMTP_PASSWORD=4fb36f4ddffd25 \
	SMTP_HOST=smtp.mailtrap.io \
	SMTP_PORT=2525 \
	SMTP_FROM_EMAIL=noreply@podcast.dev \
	SMTP_FROM_NAME="Podcast Platform" \
	SENDGRID_API_KEY="" \
	PUBLIC_URL=http://localhost:5000 \
	air

up:
	docker compose -f compose.dev.yml up -d --build --force-recreate

test:
	PORT=5005 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=password \
	POSTGRES_DB=podcast \
	JWT_ACCESS_SECRET=test \
	JWT_REFRESH_SECRET=testrefresh \
	JWT_ACCESS_EXPIRE=5 \
	JWT_ACCESS_COOKIE_EXPIRE=5 \
	JWT_REFRESH_EXPIRE=7 \
	JWT_REFRESH_COOKIE_EXPIRE=7 \
	SMTP_IDENTITY="Podcast Platform" \
	SMTP_USERNAME=221546739a5c13 \
	SMTP_PASSWORD=4fb36f4ddffd25 \
	SMTP_HOST=smtp.mailtrap.io \
	SMTP_PORT=2525 \
	FROM_EMAIL=noreply@podcast.dev \
	FROM_NAME="Podcast Platform" \
	SENDGRID_API_KEY="" \
	PUBLIC_URL=http://localhost:5000 \
	go test ./...

run:
	PORT=5000 \
	POSTGRES_HOST=db.diwtbugkffxajlpdekss.supabase.co \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=hCaOhqwpxN9XX07m \
	POSTGRES_DB=postgres \
	JWT_ACCESS_SECRET=test \
	JWT_REFRESH_SECRET=testrefresh \
	JWT_ACCESS_EXPIRE=5 \
	JWT_ACCESS_COOKIE_EXPIRE=5 \
	JWT_REFRESH_EXPIRE=7 \
	JWT_REFRESH_COOKIE_EXPIRE=7 \
	SMTP_IDENTITY="Podcast Platform" \
	SMTP_USERNAME=221546739a5c13 \
	SMTP_PASSWORD=4fb36f4ddffd25 \
	SMTP_HOST=smtp.mailtrap.io \
	SMTP_PORT=2525 \
	SMTP_FROM_EMAIL=noreply@podcast.dev \
	SMTP_FROM_NAME="Podcast Platform" \
	SENDGRID_API_KEY="" \
	PUBLIC_URL=http://localhost:5000 \
	go run .

prod:
	go build
	GIN_MODE=release \
	PORT=5000 \
	POSTGRES_HOST=db.diwtbugkffxajlpdekss.supabase.co \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=hCaOhqwpxN9XX07m \
	POSTGRES_DB=postgres \
	JWT_ACCESS_SECRET=test \
	JWT_REFRESH_SECRET=testrefresh \
	JWT_ACCESS_EXPIRE=5 \
	JWT_ACCESS_COOKIE_EXPIRE=5 \
	JWT_REFRESH_EXPIRE=7 \
	JWT_REFRESH_COOKIE_EXPIRE=7 \
	SMTP_IDENTITY="Podcast Platform" \
	SMTP_USERNAME=221546739a5c13 \
	SMTP_PASSWORD=4fb36f4ddffd25 \
	SMTP_HOST=smtp.mailtrap.io \
	SMTP_PORT=2525 \
	SMTP_FROM_EMAIL=noreply@podcast.dev \
	SMTP_FROM_NAME="Podcast Platform" \
	SENDGRID_API_KEY="" \
	PUBLIC_URL=http://localhost:5000 \
	./podcast

wire:
# go install github.com/google/wire/cmd/wire@latest
# wire gen ./middleware ./services ./routes
# wire gen ./...
	wire ./...