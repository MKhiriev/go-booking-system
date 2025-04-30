# 📅Go Booking System

Created for small companies which are working within one building.

## ▶ Run Project
### Prerequisites
- `go 1.24.0`
- `Terminal`: `git` command support
- SQL Database:
  - `PostgreSQL` or other
- `Docker`: optional

### Create Database (For Docker users)
1. Clone repository
```bash
git clone https://github.com/MKhiriev/go-booking-system.git
```

2. Clone `PostgreSQL` image
```bash
docker pull postgres
```
Version is optional
```bash
docker pull postgres:17.3
```

3. Create and run PostgreSQL container. 

You can change given command arguments:
- Name of the container: `some-postgres`
- PostgreSQL Password: `yourpassword` 
```bash
docker run -d --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=yourpassword postgres
```

4. Verify the Running `postgres` Container
```bash
docker ps
```

5. Connect to PostgreSQL
Replace `some-postgres` with name of the container specified from Step 3: `3. Create and run PostgreSQL container.`.

```bash
docker exec -it some-postgres psql -U postgres
```

6. Create database

```bash
CREATE DATABASE "go-booking";
```

### Connect database to `go-booking-system`
1. Create `.env` config file
```bash
touch .env
```
Replace `yourpassword` with password specified from: `3. Create and run PostgreSQL container.`.
```bash
echo DB_PASSWORD=yourpassword > .env
```

2. Create `config.yaml` config file
```bash
mkdir -p internal/configs
touch internal/configs/config.yaml
```
Replace `Username` and `DBName` with your PostgreSQL username and database name
```bash
cat > internal/configs/config.yaml <<EOL
server:
  port: "8080"
    
db:
  Host: "localhost"
  Port: 5432
  Username: "postgres"
  DBName: "go-booking"
EOL
```


### Create tables in database
1. Ensure that database (PostgreSQL) is running
```bash
docker ps
```
2. Connect to database (PostgreSQL)

Replace `some-postgres` with name of your container.
```bash
docker exec -it some-postgres psql -U postgres
```
3. Execute sql queries **line by line** or **by executing whole files** from `pkg/schema` directory
- Initialize tables: `1_init_up.sql`
- Fill tables with test data: `1_init_data.sql`
- Delete all tables and data: `1_init_down.sql`