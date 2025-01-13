### 1. PREQUISITES
- Go installed
- Docker installed (or use the command below)
- Docker Compose installed (or use the command below)

### 2. INSTALL DOCKER AND DOCKER COMPOSE
To install Docker and Docker Compose, use the following commands:

1. **Install Docker**:
   ```sh
   curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh && rm get-docker.sh
   ```

2. **Check Docker Version**:
   ```sh
   docker --version
   ```

3. **Install Docker Compose**:
   ```sh
   sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
   ```

4. **Check Docker Compose Version**:
   ```sh
   docker-compose --version
   ```

### 3. CLONE THE FRONTEND REPOSITORY
For development purposes, the frontend repository is included as a submodule.
Before running Docker Compose, clone the frontend repository into the `frontend` folder:

```sh
git clone https://github.com/hotaryuzaki/portier-nextjs-frontend.git frontend
```

### 4. RUN THE APPLICATION
To run the setup, follow these steps:

1. **Start Docker Compose**:
   Use the `docker-compose-up` command from the `Makefile` to start all services.
   ```sh
   make docker-compose-up
   ```

2. **Build the Backend**:
   Use the `build-backend` command from the `Makefile` to build the Go backend.
   ```sh
   make build-backend
   ```

3. **Build the Frontend**:
   Use the `build-frontend` command from the `Makefile` to build the Next.js frontend.
   ```sh
   make build-frontend
   ```

4. **Run Database Migrations**:
   Use the `migrate-up` command from the `Makefile` to run the database migrations.
   ```sh
   make migrate-up
   ```

### 5. CONNECT TO THE DATABASE
To connect to the PostgreSQL database, use the following command:
```sh
make psql
```

This command will open a `psql` session in the PostgreSQL container using the credentials specified in the `.env` file.

### 6. CHECK DOCKER CONTAINERS
To check the status of Docker containers, use the following command:
```sh
docker ps
```

### 7. DOCKER COMMANDS
Additional Docker commands available in the `Makefile`:

- **Build Docker images**:
  ```sh
  make docker-build
  ```

- **Stop Docker Compose**:
  ```sh
  make docker-compose-down
  ```

- **Restart Docker Compose**:
  ```sh
  make docker-compose-restart
  ```

- **View Docker Compose logs**:
  ```sh
  make docker-compose-logs
  ```

- **Clean Docker containers, networks, and volumes**:
  ```sh
  make docker-clean
  ```

These commands will start the Docker containers, build the backend and frontend applications, and run the database migrations. Make sure your `.env` file is correctly configured with the necessary environment variables.
