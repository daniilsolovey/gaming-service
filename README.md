# Handle http requests from client, send requests to game platform

Build app:
```
make build
```

Run tests:
```
make test
```

Build docker container:
```
docker build -t gaming-service .
```

Run docker container:

Command for run database:

```
docker run --name gaming-service-postgres-db -e POSTGRES_PASSWORD=your_password -p 5432:5432 postgres
```

Don't forget to create a database:
```
CREATE DATABASE databasename;
```

Run container with app:

```
docker-start -i gaming-service
```
