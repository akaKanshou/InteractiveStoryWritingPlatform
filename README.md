# ForgeVerse
ForgeVerse is an interactive story writing platform. Here users can interactive stories having which runs on user choices made in the story. Users can create a story, add multiple chapters to it, and connect those chapters together with edges so readers can jump from one chapter to another and follow different paths through the story.This project is built for Devkriti 2026.

Features Implemented
Frontend
> It contains a leftside horizontal navbar to give it a better look and easy navigation.
> The home contains a searchbar which can be used to search stories.
> The home page contains a number of stories having different lists having a scroll view.
> login page by user can create an account and login using google authentication.
> User registration page to setting usernames for new users.
> My Stories page to view all created stories and manage these.
> Create story page to create new story having sections Title, description and a public private system.
> Graph view to visualize how chapters are connected to each other.
> View story page, which shows story details and chapters.
>View chapter page to simulate a chapter.

Backend
----Backend Features----
> Users can log-in using Google-OAuth for seamless and secure authentication
> Users can create and write stories, adding chapters as they go
> Users can create more than just linear stories by adding arbitrary edges between chapters
> Users can add edges to other people's stories to create the best environment of collaborative work and fan-fiction
> The back-end is written completely in Go, making it responsive and sturdy, while being as modern and idiomatic as possible


Technologies/Libraries/Packages Used
>Go – core backend language
>Gin – HTTP web framework for routing and middleware
>PostgreSQL (via pgx) – database
>Gorilla Sessions – cookie based session management
>golang.org/x/oauth2 – Google OAuth authentication
>godotenv – loading environment variables from .env file
>gin-contrib/cors – handling CORS
>HTML, CSS, and JavaScript – frontend pages and scripts

Local Setup
>Clone the repository
>Install Go dependencies
>Set up PostgreSQL database
>Create a PostgreSQL database.
>Import the schema/data using the provided dump: db_dump.sql
>Create a .env file in the project root with the following variables:
   PORT=8080
   PG_USER=your_postgres_user
   PG_PASSWORD=your_postgres_password
   PG_HOST=localhost
   PG_PORT=5432
   PG_DB=your_database_name
   SESSION_SECRET=your_session_secret
   GOOGLE_OAUTH_CLIENTID=your_google_client_id
   GOOGLE_OAUTH_CLIENTSECRET=your_google_client_secret

>The server runs on HTTPS, so add a public_key and private_key file in the project root (you can generate a self-signed certificate for local development).
>Run the project using this command in terminal in root directory:  go run cmd/api/main.go
>Visit https://localhost:8080 in your browser.
> Execute the following commands (in the postgres binary directory, which you is possibly at "C:\Program Files\PostgreSQL\18\bin")
   .\dropdb.exe -U postgres Forgeverse
   .\createdb.exe -U postgres Forgeverse
   .\psql.exe -U postgres -h localhost -d Forgeverse -f <path_to_sql_dump>
> Remove the leading underscores from the env and TLS key files.

Team Members
>Siddarth Anil Nair
>Chetan Disania
>Deepak Paswan

