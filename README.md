
# Hospital Management API

## Overview

The **Hospital Management API** is designed to manage authentication and related operations for a hospital management system. This API handles tasks such as user login, registration, Personel management, polyclinic management, and password reset functionalities. 

## Prerequisites

Before you start, ensure you have the following installed:

- [Go 1.20+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose (optional)](https://docs.docker.com/compose/install/)

## Getting Started

### 1. Clone the Repository

Clone the repository to your local machine:

bash
git clone https://github.com/your-username/your-repo.git
cd your-repo


### 2. Dockerize the Application

To run the application using Docker:

#### Build the Docker Image

bash
docker build -t auth-api .


#### Run the Docker Container

bash
docker run -d -p 3000:3000 auth-api

## API Documentation

### Base URL

The API is hosted at:


http://localhost:3000/


### Authentication

All requests to protected endpoints require a Bearer Token in the `Authorization` header.

### Endpoints

#### Auth Endpoints

- **Login a User**
  - **POST** `/login`
  - Description: Login a user and return a JWT token.
  - Parameters: `email` or `tel_no`, `password`
  - Responses: 
    - \`200 OK\`: Returns a JWT token
    - \`400 Bad Request\`
    - \`401 Unauthorized\`

- **Register a New User**
  - **POST** `/register`
  - Description: Register a new user and their associated hospital.
  - Parameters: `hospital`, `user`
  - Responses:
    - \`200 OK\`: Returns the created user and hospital
    - \`400 Bad Request\`
    - \`500 Internal Server Error\`

#### Personel Endpoints

- **Add New Personel**
  - **POST** `/personel`
  - Description: Add a new Personel to the hospital.
  - Parameters: `name`, `surname`, `title`, `tc_no`, `tel_no`, `job_type`, `working_days`
  - Responses:
    - \`200 OK\`: Returns the created Personel
    - \`400 Bad Request\`
    - \`401 Unauthorized\`
    - \`500 Internal Server Error\`

- **Get Personel by ID**
  - **GET** `/personel/{id}`
  - Description: Retrieve Personel details by ID.
  - Parameters: `id` (in path)
  - Responses:
    - \`200 OK\`: Returns the Personel details
    - \`404 Not Found\`

- **Update Personel**
  - **PUT** `/personel/{id}`
  - Description: Update Personel details by ID.
  - Parameters: `id` (in path), `name`, `surname`, `title`, `tc_no`, `tel_no`, `job_type`, `working_days`
  - Responses:
    - \`200 OK\`: Returns the updated Personel
    - \`400 Bad Request\`
    - \`404 Not Found\`
    - \`500 Internal Server Error\`

- **Delete Personel**
  - **DELETE** `/personel/{id}`
  - Description: Delete Personel by ID.
  - Parameters: `id` (in path)
  - Responses:
    - \`200 OK\`
    - \`404 Not Found\`
    - \`500 Internal Server Error\`

#### Polyclinic Endpoints

- **Get Polyclinics**
  - **GET** `/polyclinics`
  - Description: Retrieve a list of polyclinics associated with the hospital.
  - Responses:
    - \`200 OK\`: Returns a list of polyclinics
    - \`500 Internal Server Error\`

- **Add Polyclinic**
  - **POST** `/polyclinics`
  - Description: Add a new polyclinic to the hospital.
  - Parameters: `polyclinic_name`, `city`, `district`
  - Responses:
    - \`200 OK\`: Returns the created polyclinic
    - \`400 Bad Request\`
    - \`500 Internal Server Error\`

- **Delete Polyclinic**
  - **DELETE** `/polyclinics/{id}`
  - Description: Delete a polyclinic by ID.
  - Parameters: `id` (in path)
  - Responses:
    - \`200 OK\`
    - \`404 Not Found\`
    - \`500 Internal Server Error\`

#### Password Reset Endpoints

- **Initiate Password Reset**
  - **POST** `/reset-password/initiate`
  - Description: Send a reset code to the user's phone number.
  - Parameters: `tel_no`
  - Responses:
    - \`200 OK\`
    - \`400 Bad Request\`
    - \`404 Not Found\`
    - \`409 Conflict\`

- **Finalize Password Reset**
  - **POST** `/reset-password/finalize`
  - Description: Validate the reset code and get a token to change the password.
  - Parameters: `reset code`, `tel_no`
  - Responses:
    - \`200 OK\`: Returns the token for password change
    - \`400 Bad Request\`
    - \`401 Unauthorized\`
    - \`404 Not Found\`

- **Reset Password**
  - **POST** `/reset-password`
  - Description: Change the user's password using the validation token.
  - Parameters: `password`, `repeat_password`, `validate_code`
  - Responses:
    - \`200 OK\`
    - \`400 Bad Request\`
    - \`404 Not Found\`
    - \`500 Internal Server Error\`

#### User Management Endpoints

- **Get All Users**
  - **GET** `/users`
  - Description: Retrieve a list of all users associated with the hospital.
  - Responses:
    - \`200 OK\`: Returns a list of users
    - \`500 Internal Server Error\`

- **Create a New User**
  - **POST** `/users`
  - Description: Add a new user to the hospital.
  - Parameters: `name`, `surname`, `email`, `password`, `role`, `tcNo`, `telNo`
  - Responses:
    - \`200 OK\`: Returns the created user
    - \`400 Bad Request\`
    - \`500 Internal Server Error\`

- **Get User by ID**
  - **GET** `/users/{id}`
  - Description: Retrieve a single user by ID.
  - Parameters: `id` (in path)
  - Responses:
    - \`200 OK\`: Returns the user details
    - \`404 Not Found\`

- **Update User**
  - **PUT** `/users/{id}`
  - Description: Update user details by ID.
  - Parameters: `id` (in path), `name`, `surname`, `email`, `password`, `role`, `tcNo`, `telNo`
  - Responses:
    - \`200 OK\`: Returns the updated user
    - \`400 Bad Request\`
    - \`404 Not Found\`
    - \`500 Internal Server Error\`

## Testing

To run the tests:

bash
go test ./...


## Deployment

You can deploy this application using Docker or any platform that supports containerized applications. Build and push your Docker image to a registry, then deploy it on your server or cloud platform.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (\`git checkout -b feature/your-feature\`).
3. Make your changes and commit them (\`git commit -m 'Add some feature'\`).
4. Push to the branch (\`git push origin feature/your-feature\`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
