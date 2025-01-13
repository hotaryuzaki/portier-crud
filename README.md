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

5. **Add Docker Compose to PATH** (if needed):
   If Docker Compose is not found in your PATH, you can add it manually:
   ```sh
   export PATH=$PATH:/usr/local/bin
   ```

   To make this change permanent, add the above line to your `.bashrc` file:
   ```sh
   echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
   source ~/.bashrc
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

2. **Run Database Migrations**:
   Use the `migrate-up` command from the `Makefile` to run the database migrations.
   ```sh
   make migrate-up
   ```

After completing these steps, all applications should be running. You can verify this by checking the status of the Docker containers.

### 5. CHECK DOCKER CONTAINERS
To check the status of Docker containers, use the following command:
```sh
docker ps
```

### 6. CONNECT TO THE DATABASE
To connect to the PostgreSQL database, use the following command:
```sh
make psql
```

This command will open a `psql` session in the PostgreSQL container using the credentials specified in the `.env` file.

### 7. CREATE INITIAL DATA
**Note: Before creating a user, you must first create a tenant.**

To create a tenant, use the following `curl` command:
```sh
curl -X POST http://localhost:4000/tenants \
-H "Content-Type: application/json" \
-d '{"name": "PT ZIG ZAG", "address": "Jln banyak belok", "status": "Active"}'
```

**Note: Before creating keys and copies, you must first create a user.**

To create a user, use the following `curl` command:
```sh
curl -X POST http://localhost:4000/users \
-H "Content-Type: application/json" \
-d '{"username": "ahmad", "email": "ahmadamri.id@gmail.com", "password": "securepassword123", "name": "ahmad amri sanusi", "gender": "1", "id_number": "123456789", "user_image": "http://example.com/image.jpg", "tenant_id": 1}'
```

To create a key, use the following `curl` command:
```sh
curl -X POST http://localhost:4000/keys \
-H "Content-Type: application/json" \
-d '{"name": "TEST Key"}'
```

To create a copy, use the following `curl` command:
```sh
curl -X POST http://localhost:4000/copies \
-H "Content-Type: application/json" \
-d '{"name": "TEST Copy", "key_id": 1}'
```

### 8. BACKEND AND FRONTEND URLS
- **Backend URL**: `http://localhost:4000`
- **Frontend URL**: `http://localhost:3000`

### 9. AVAILABLE ROUTES

#### Backend Routes
- **User Routes**:
  - `GET /users`
  - `GET /users/:id`
  - `POST /users`
  - `PUT /users/:id`
  - `DELETE /users/:id`

- **Key Routes**:
  - `GET /keys`
  - `GET /keys/:id`
  - `POST /keys`
  - `PUT /keys/:id`
  - `DELETE /keys/:id`

- **Copy Routes**:
  - `GET /copies`
  - `GET /copies/:id`
  - `POST /copies`
  - `PUT /copies/:id`
  - `DELETE /copies/:id`

- **Tenant Routes**:
  - `GET /tenants`
  - `GET /tenants/:id`
  - `POST /tenants`
  - `PUT /tenants/:id`
  - `DELETE /tenants/:id`

#### Frontend Routes
The frontend routes are defined in the `frontend/src/pages` directory of the frontend repository.

### 10. DOCKER COMMANDS
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

These commands will start the Docker containers and run the database migrations. Make sure your `.env` file is correctly configured with the necessary environment variables.

### 11. DEVELOPMENT (TEST) PURPOSES
For development (test) purposes, the following are allowed:
- Creating a user without the `tenant_id` object.
- Creating a copy without the `key_id` object.
- `created_by` is stored by default value due to no authentication.
