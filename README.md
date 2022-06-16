# Convert JSON from API POINTS CITY DISCORD to CSV

# Create Project
- Run `go mod init github.com/LucioTrucco/discord-points-votes`
- Run `go get .`
- Run `go mod tidy`

## Initial DynamoDBLocal
```shell
docker run -p 8000:8000 amazon/dynamodb-local
```

## HOW TO RUN
- Run `go run . -url https://points.city/api/guilds/961074073868308480/leaderboard -cookie "pointsID=xlTF88yqTNfUKr-NyD-nZdmImAk1S-7P.KCD3%2BTmVMC67bdOuYWa8SgyfhQ4P9c%2B7X8YtSYB43A0; Path=/; Expires=Sun, 17 Jul 2022 01:16:51 GMT; HttpOnly" -table discord-users-dev -phase phase_1`