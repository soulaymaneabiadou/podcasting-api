dev:
	PORT=5000 \
	POSTGRES_HOST=db.diwtbugkffxajlpdekss.supabase.co \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=hCaOhqwpxN9XX07m \
	POSTGRES_DB=postgres \
	air

up:
	docker compose -f compose.dev.yml up -d 

prod:
	PORT=5000 \
	POSTGRES_HOST=db.diwtbugkffxajlpdekss.supabase.co \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=hCaOhqwpxN9XX07m \
	POSTGRES_DB=postgres \
	go run .