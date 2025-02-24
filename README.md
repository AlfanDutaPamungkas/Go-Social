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
- [](github.com/go-chi/chi/v5) v5 or higher

## Instalation
1. Clone the repository:
    ```bash
    git clone https://github.com/AlfanDutaPamungkas/Meals-App-RESTful-API.git
    ```
2. Navigate to the project directory:
    ```bash
    cd Inventory-System-RESTful-API
    ```
3. Install dependencies:
    ```bash
    go mod download
    ```
4. Set up your environment variables:
    Create a `.env` file in the project root and specify the following variables:
    ```env
    JWT_TOKEN_SECRET=your_jwt_secret_key
    CLOUD_NAME=your_cloduinary_cloud_name
    CLOUDINARY_API_KEY=your_cloudinary_api_key
    CLOUDINARY_API_SECRET=your_cloudinary_api_secret
    DB_URL=your_db_url
    ```
5. Start the server:
    ```bash
    go run main.go
    ```
    The API will be running at `http://localhost:3000`.

## API Documentation (OpenAPI 3.0)

The API is fully documented using the OpenAPI 3.0 specification. You can view the  `apispec.json`

Contributing
Feel free to open issues or submit pull requests if you want to contribute to this project.
