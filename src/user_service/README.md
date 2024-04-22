## Getting Started

Before you start using the project, make sure you have completed the following steps:

* **OpenAPI Specification**: Ensure that the API endpoints are described in the OpenAPI specification file located in the `/api` directory.

* **Generate RSA Keys**: Generate a pair of RSA keys and place them in the /build directory. Example keys are provided in the `/build` directory for reference.

* **Generate .env File**: Generate a `.env` file and place it into the `/deployments` directory. File `example.env` is provided for reference in the `/deployments` directory.

* **Configuration Setup**: You can customize the names of the key pair and specify the directory where they should be searched for in the `config.yaml` file located in the `/build` directory. Additionally, you can set the token duration before expiration in the same `config.yaml` file. The path to configuration files is passed as a command-line argument in the `docker-compose.yaml` file.

## Configuration

### RSA Key Pair

* Location: `/build`
* Example: `example_signature.pub` and `example_signature.pem`

### config.yaml

* Location: `/build`
* Customizable Parameters:
  * Key pair names
  * Key pair directory
  * Token expiration duration

### .env File

* Location: `/deployments`
* Example: `example.env`
* Customizable Parameters:
  * Postgres passrod
  * Postgres user
  * Postgres db
  * Postgres port
  * Main service port

## Request examples

### Register

```bash
> curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{
  "username": "example_user",
  "password": "example_password"
}'

{
    "expiration_date": "2024-04-22T00:48:29.682763145Z",
    "token": JWT_TOKEN
}
```

### Login

```bash
> curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{
  "username": "example_user",
  "password": "example_password"
}'

{
    "expiration_date":"2024-04-22T01:00:22.514470901Z",
    "token": JWT_TOKEN
}
```

### Update

```bash
>curl -X PUT http://localhost:8080/update \
-H "Content-Type: application/json" \
-H "Authorization: Bearer JWT_TOKEN" \
-d '{
  "name": "New Name",
  "surname": "New Surname",
  "date_of_birth": "1990-01-01",
  "email": "new@example.com",
  "phone_number": "1234567890"
}'

{
    "message":"User data updated successfully"
}
```
