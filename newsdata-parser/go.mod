module github.com/qsoulior/news/newsdata-parser

go 1.22.0

require (
	github.com/qsoulior/news/aggregator v0.0.0-00010101000000-000000000000
	github.com/qsoulior/news/parser v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/redis/go-redis/v9 v9.5.1 // indirect
)

replace github.com/qsoulior/news/aggregator => ../aggregator

replace github.com/qsoulior/news/parser => ../parser
