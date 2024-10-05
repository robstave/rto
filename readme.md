

```
docker compose build
# not docker-compose

docker compose up -d
docker compose down -d
```



# Testing
Generate the mocks with mockery.  
Install as needed

```
go generate ./...
```