export APP_PORT=":8080"
export DB_LOCATION="$PWD/posts.db"
cd src; go run cmd/api/main.go