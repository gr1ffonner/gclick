wait-for "${POSTGRESQL}:${5432}" -- "$@"

# Watch your .go files and invoke go build if the files changed.
CompileDaemon --build="go build -o main cmd/main.go"  --command=./main 
