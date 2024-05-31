wait-for "${CLICKHOUSE_HOST}:${CLICKHOUSE_PORT}" -- "$@"

# Watch your .go files and invoke go build if the files changed.
CompileDaemon --build="go build -o main cmd/main.go"  --command=./main 
