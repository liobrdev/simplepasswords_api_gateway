# SimplePasswords - API Gateway

To be used in front of SimplePasswords vaults microservice.

## System Requirements
- Latest version of Go

## Download

Clone this repository: `git clone https://github.com/liobrdev/simplepasswords_api_gateway.git`

## Required Environment Variables

Each of the environment variables in the following table **must** be an **absolute path** to an existent `UTF-8`-encoded text file. The contents of the first line of this text file will be parsed to a Go data type (empty file contents will be parsed with a default value). This data will then be saved to the AppConfig field whose name matches that of the corresponding environment variable.

| **Required Environment Variable** | **Description of File Contents** | **Data Type** | **Default Value** |
|-----------------------------------|----------------------------------|---------------|-------------------|
| GO_FIBER_ENVIRONMENT | Should be either `development`, `production`, or `testing`. | `string` | `"development"` |
| GO_FIBER_BEHIND_PROXY | Configures Fiber server setting `EnableTrustedProxyCheck bool`. Should be either `true` or `false`. | `bool` | `false` |
| GO_FIBER_PROXY_IP_ADDRESSES | Configures Fiber server setting `TrustedProxies []string`. Should be comma-separated IP address string(s). | `[]string` | `[]string{""}` |
| GO_FIBER_API_GATEWAY_DB_USER | Main PostgreSQL database user. | `string` | `""` |
| GO_FIBER_API_GATEWAY_DB_PASSWORD | Main PostgreSQL database password. | `string` | `""` |
| GO_FIBER_API_GATEWAY_DB_HOST | Main PostgreSQL database host. | `string` | `""` |
| GO_FIBER_API_GATEWAY_DB_PORT | Main PostgreSQL database port. | `string` | `""` |
| GO_FIBER_API_GATEWAY_DB_NAME | Main PostgreSQL database name. | `string` | `""` |
| GO_FIBER_LOGGER_DB_USER | Logger PostgreSQL database user. | `string` | `""` |
| GO_FIBER_LOGGER_DB_PASSWORD | Logger PostgreSQL database password. | `string` | `""` |
| GO_FIBER_LOGGER_DB_HOST | Logger PostgreSQL database host. | `string` | `""` |
| GO_FIBER_LOGGER_DB_PORT | Logger PostgreSQL database port. | `string` | `""` |
| GO_FIBER_LOGGER_DB_NAME | Logger PostgreSQL database name. | `string` | `""` |
| GO_FIBER_REDIS_PASSWORD | Redis cache password. | `string` | `""` |
| GO_FIBER_SECRET_KEY | Secret key for various app-level encryption methods. | `string` | `""` |
| GO_FIBER_SERVER_HOST | Fiber app will be run from this host. | `string` | `"localhost"` |
| GO_FIBER_SERVER_PORT | Fiber app will be run from host using this port. | `string` | `"5050"` |
| GO_FIBER_VAULTS_URL | URL of SimplePasswords vaults microservice. | `string` | `"http://localhost:8080"` |

### Methods For Setting Environment Variables

Required environment variables may be sourced by:

1. setting variables via the command line at compile time,

and/or,

2. setting variables in the shell environment before compile time,

and/or,

3. including a `.env` file in the root application folder before compile time.

Each environment variable **must** be set using at least one of these three methods. Variables set in the shell environment will override duplicate variables included in a `.env` file, and variables set via the command line will override duplicate variables set in the shell environment and/or included in a `.env` file. That is, environment variable sources have the following precedence: command line, *then* shell, *then* `.env` file.

#### An example using all three methods:

With the following `.env` file present in the root application folder:

```bash
# .env

GO_FIBER_ENVIRONMENT=/path/to/secret_files/environment_3
GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_3
GO_FIBER_PROXY_IP_ADDRESSES=/path/to/secret_files/proxy_ip_addresses
GO_FIBER_API_GATEWAY_DB_USER=/path/to/secret_files/api_gateway_db_user
GO_FIBER_API_GATEWAY_DB_PASSWORD=/path/to/secret_files/api_gateway_db_password
GO_FIBER_API_GATEWAY_DB_HOST=/path/to/secret_files/api_gateway_db_host
GO_FIBER_API_GATEWAY_DB_PORT=/path/to/secret_files/api_gateway_db_port
GO_FIBER_API_GATEWAY_DB_NAME=/path/to/secret_files/api_gateway_db_name
GO_FIBER_LOGGER_DB_USER=/path/to/secret_files/logger_db_user
GO_FIBER_LOGGER_DB_PASSWORD=/path/to/secret_files/logger_db_password
GO_FIBER_LOGGER_DB_HOST=/path/to/secret_files/logger_db_host
GO_FIBER_LOGGER_DB_PORT=/path/to/secret_files/logger_db_port
GO_FIBER_LOGGER_DB_NAME=/path/to/secret_files/logger_db_name
GO_FIBER_REDIS_PASSWORD=/path/to/secret_files/redis_password
GO_FIBER_SECRET_KEY=/path/to/secret_files/secret_key
GO_FIBER_SERVER_HOST=/path/to/secret_files/server_host
GO_FIBER_SERVER_PORT=/path/to/secret_files/server_port
GO_FIBER_VAULTS_URL=/path/to/secret_files/vaults_url
```

Then running the following commands:

```bash
export GO_FIBER_ENVIRONMENT=/path/to/secret_files/environment_2
export GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_2
GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_1 go build
```

The resulting executable will be built with `GO_FIBER_ENVIRONMENT` set to `/path/to/secret_files/environment_2`, `GO_FIBER_BEHIND_PROXY` set to `/path/to/secret_files/behind_proxy_1`, and all other required variables set to their corresponding values in the `.env` file.

### Rationale

The scheme described above is particularly convenient for use with Docker Compose `secrets` configuration:

```yaml
# simplepasswords/docker-compose.yml

version: '3.9'

services:
    api_gateway:
        build:
            context: https://github.com/liobrdev/simplepasswords_api_gateway
        command: go run main.go
        secrets:
            - redis_password
            - api_gateway_environment
            - api_gateway_behind_proxy
            - api_gateway_proxy_ip_addresses
            - api_gateway_db_host
            - api_gateway_db_name
            - api_gateway_db_password
            - api_gateway_db_port
            - api_gateway_db_user
            - logger_db_host
            - logger_db_name
            - logger_db_password
            - logger_db_port
            - logger_db_user
            - api_gateway_secret_key
            - api_gateway_server_host
            - api_gateway_server_port
            - vaults_url
        environment:
            GO_FIBER_ENVIRONMENT: /run/secrets/api_gateway_environment
            GO_FIBER_BEHIND_PROXY: /run/secrets/api_gateway_behind_proxy
            GO_FIBER_PROXY_IP_ADDRESSES: /run/secrets/api_gateway_proxy_ip_addresses
            GO_FIBER_API_GATEWAY_DB_HOST: /run/secrets/api_gateway_db_host
            GO_FIBER_API_GATEWAY_DB_NAME: /run/secrets/api_gateway_db_name
            GO_FIBER_API_GATEWAY_DB_PASSWORD: /run/secrets/api_gateway_db_password
            GO_FIBER_API_GATEWAY_DB_PORT: /run/secrets/api_gateway_db_port
            GO_FIBER_API_GATEWAY_DB_USER: /run/secrets/api_gateway_db_user
            GO_FIBER_LOGGER_DB_HOST: /run/secrets/logger_db_host
            GO_FIBER_LOGGER_DB_NAME: /run/secrets/logger_db_name
            GO_FIBER_LOGGER_DB_PASSWORD: /run/secrets/logger_db_password
            GO_FIBER_LOGGER_DB_PORT: /run/secrets/logger_db_port
            GO_FIBER_LOGGER_DB_USER: /run/secrets/logger_db_user
            GO_FIBER_REDIS_PASSWORD: /run/secrets/redis_password
            GO_FIBER_SECRET_KEY: /run/secrets/api_gateway_secret_key
            GO_FIBER_SERVER_HOST: /run/secrets/api_gateway_server_host
            GO_FIBER_SERVER_PORT: /run/secrets/api_gateway_server_port
            GO_FIBER_VAULTS_URL: /run/secrets/vaults_url
        ports:
            - 5050:5050
        depends_on:
            - api_gateway_db
            - logger_db
            - redis
            - vaults
    api_gateway_db:
        # ...
    logger_db:
        # ...
    redis:
        # ...
    vaults:
        # ...
    # ...
secrets:
    redis_password:
        file: ./secret_files/redis_password.txt
    api_gateway_behind_proxy:
        file: ./secret_files/api_gateway_behind_proxy.txt
    api_gateway_db_host:
        file: ./secret_files/api_gateway_db_host.txt
    api_gateway_db_name:
        file: ./secret_files/api_gateway_db_name.txt
    api_gateway_db_password:
        file: ./secret_files/api_gateway_db_password.txt
    api_gateway_db_port:
        file: ./secret_files/api_gateway_db_port.txt
    api_gateway_db_user:
        file: ./secret_files/api_gateway_db_user.txt
    logger_db_host:
        file: ./secret_files/logger_db_host.txt
    logger_db_name:
        file: ./secret_files/logger_db_name.txt
    logger_db_password:
        file: ./secret_files/logger_db_password.txt
    logger_db_port:
        file: ./secret_files/logger_db_port.txt
    logger_db_user:
        file: ./secret_files/logger_db_user.txt
    api_gateway_environment:
        file: ./secret_files/api_gateway_environment.txt
    api_gateway_secret_key:
        file: ./secret_files/api_gateway_secret_key.txt
    api_gateway_server_host:
        file: ./secret_files/api_gateway_server_host.txt
    api_gateway_server_port:
        file: ./secret_files/api_gateway_server_port.txt
    api_gateway_proxy_ip_addresses:
        file: ./secret_files/api_gateway_proxy_ip_addresses.txt
    vaults_url:
        file: ./secret_files/vaults_url.txt
    # ...
```

## Install Dependencies, Build, & Run

After configuring required environment variables as explained above, run the following commands from the root application folder:

```bash
go mod download
go build
./simplepasswords_api_gateway
```

The server should now be running at whichever host and port were loaded from `GO_FIBER_SERVER_HOST` and `GO_FIBER_SERVER_PORT` environment variables respectively.
