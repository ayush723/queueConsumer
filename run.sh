MONGO_URI="mongodb://root:rootpassword@localhost:27017/test_db?authSource=admin" \
MONGO_DATABASE=test_db \
RABBITMQ_URI="amqp://guest:guest@localhost:5672/" \
RABBITMQ_QUEUE=rss_urls \
go run main.go
