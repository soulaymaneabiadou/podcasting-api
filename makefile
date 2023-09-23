dev:
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
	air

up:
	docker compose -f compose.dev.yml up -d 
