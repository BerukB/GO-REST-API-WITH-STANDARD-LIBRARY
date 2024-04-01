# GO-REST-API-WITH-STANDARD-LIBRARY

This repository contains a simple RESTful API built with Go's standard library. It provides functionality for user authentication, file upload, and image serving.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY.git
    ```

2. Navigate to the project directory:

    ```bash
    cd GO-REST-API-WITH-STANDARD-LIBRARY
    ```

3. Install dependencies:

    ```bash
    # No external dependencies required as the project uses only Go's standard library.
    ```

## Usage

### Running the Server

To start the server, run the following command:

```bash
go run main.go
```

The server will start listening on port `8080` by default.

### Endpoints

- **Login Endpoint**: 
  - URL: `/login`
  - Method: POST
  - Description: Endpoint for user authentication. It returns access and refresh tokens upon successful login.

- **File Upload Endpoint**:
  - URL: `/upload`
  - Method: POST
  - Description: Endpoint for uploading files. It expects a file field named "image" in the request.

- **Image Serving Endpoint**:
  - URL: `/serve/{imageName}`
  - Method: GET
  - Description: Endpoint for serving images. Replace `{imageName}` with the filename of the image you want to retrieve.

- **User Management Endpoints**:
  - `/user`: 
    - Methods: GET (List all users), POST (Create a new user)
  - `/user/{userID}`:
    - Methods: GET (Retrieve user by ID), PUT (Update user by ID), DELETE (Delete user by ID)

- **Home Page**:
  - URL: `/`
  - Description: Simple home page that displays "This is my home page".

### Dependencies

- `net/http`: Provides HTTP client and server implementations.
- `github.com/dgrijalva/jwt-go`: Package for JWT handling.
- `golang.org/x/crypto/bcrypt`: Package for hashing and comparing passwords securely.

## Contributing

Contributions are welcome! Feel free to submit pull requests or open issues for any improvements or bug fixes.

## GitHub Repository

The source code for this project is hosted on GitHub:

[GO-REST-API-WITH-STANDARD-LIBRARY](https://github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY)
