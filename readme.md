## Running Service locally
### Start the service
`docker compose up --build`

This would spin up 3 docker containers:
1. meeting-service-api: to access -> `http://localhost:8080`
2. postgres:15
3. swagger-ui
To access swagger UI hit -> `http://localhost:8081`

### Stop the service
`docker compose down`

#### Also, remove DB
`docker compose down -v`


## Deploying Service

#### Deployed currently on
`https://meeting-service.fly.dev/health`

#### Fly - Prerequisites
```
fly apps create meeting-service
fly postgres create --name meeting-service-db --region bom
fly postgres attach meeting-service-db --app meeting-service
```

#### Deploy to Fly
`fly deploy`

Have tested deploying it to fly


## About code 
#### Flow
Handler/controller → Service → Repository → Database

#### Core logic
service -> scheduler.go 
This has entire recommendation logic

## Other
### Open API Spec
Can be tested locally via `http://localhost:8081`

### Tech stack:
Go, Postgres