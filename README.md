# Go-Social (Gopher) API
GopherSocial is a social media platform that allows users to create and share posts, follow other users, and engage in discussions through comments. The GopherSocial API provides endpoints for user authentication, post management, and user interactions, making it easy to integrate with web or mobile applications.

## Features
- Authentication and authorization for secure API access (admin, moderator and user)
- Users can create, update, view, and delete own posts and follow other user
- Moderator can update post user
- Admin can update and delete post user
- Rate limiting
- Swagger documentation
- Graceful shutdown
- Redis caching in get profile user

## Prerequisites
- [Golang](https://golang.org/doc/install) v1.18 or higher
- [pgx](https://github.com/jackc/pgx) or any other postgres connection pool
- [chi](github.com/go-chi/chi/v5) v5 or higher
- [go-redis](https://github.com/redis/go-redis) v9 or higher
- [swag](https://github.com/swaggo/swag) for documentation

## Instalation
1. Clone the repository:
    ```bash
    git clone https://github.com/AlfanDutaPamungkas/Go-Social.git
    ```
2. Navigate to the project directory:
    ```bash
    cd Go-Social
    ```
3. Install dependencies:
    ```bash
    go mod download
    ```
4. Set up your environment variables:
    Create a `.env` file in the project root and specify the following variables:
    ```env
    ENV=
    PORT=
    DB_ADDR=
    EXTERNAL_URL=
    SMTP_USERNAME=
    SMTP_HOST=
    SMTP_PASSWORD=
    AUTH_TOKEN_SECRET=
    REDIS_ENABLED=
    REDIS_PW=
    AUTH_BASIC_USERNAME=
    AUTH_BASIC_PASS=
    ```
5. Start the server:
    ```bash
    go run cmd/api
    ```
    The API will be running at `http://localhost:3000`.

## API Documentation (OpenAPI 3.0)

The API is fully documented using the OpenAPI 3.0 specification. You can view the  `http://localhost:3000/v1/swagger/index.html`

## Contributing
Feel free to open issues or submit pull requests if you want to contribute to this project.
