# Go REST API with Standard Library

This project implements a simple REST API using only the standard library of Go. It provides endpoints for user management, authentication, file uploading, and serving images.

## Features

- **User Management**: Create, retrieve, update, and delete user records.
- **Authentication**: Basic login functionality with password hashing and JWT authentication.
- **File Upload**: Allows users to upload image files.
- **Image Serving**: Serves uploaded images over HTTP.
- **Middleware**: Includes middleware for request timeout, request ID generation, and JWT authentication.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY.git
    ```

2. Navigate to the project directory:

    ```bash
    cd GO-REST-API-WITH-STANDARD-LIBRARY
    ```

3. Install external dependencies:

    ```bash
    go get github.com/dgrijalva/jwt-go
    go get golang.org/x/crypto/bcrypt
    ```

4. Run the following command to start the server:

    ```bash
    go run main.go
    ```

5. The server will start running on port `8080`.

## Usage

### Endpoints

- **POST /login**: User login endpoint.
- **POST /upload**: File upload endpoint.
- **GET /v1/images/{imageName}**: Image serving endpoint.
- **GET /user**: List all users.
- **GET /user/{userID}**: Get user by ID.
- **POST /user**: Create a new user.
- **PUT /user/{userID}**: Update user by ID.
- **DELETE /user/{userID}**: Delete user by ID.

### Authentication

- Authentication for user-related endpoints is handled via JWT tokens.
- To access user-related endpoints, include a valid JWT token in the request headers.

### File Upload

- To upload a file, send a POST request to `/upload` with a `multipart/form-data` body containing the file.

### Image Serving

- Uploaded images can be accessed via the `/v1/images/{imageName}` endpoint.

## Dependencies

- **github.com/dgrijalva/jwt-go**: Provides functionality for handling JWT tokens in Go.
- **golang.org/x/crypto/bcrypt**: Provides functions for hashing and comparing passwords using the bcrypt hashing algorithm.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

---

Feel free to further customize the README to include any additional information specific to your project. Let me know if you need further assistance!