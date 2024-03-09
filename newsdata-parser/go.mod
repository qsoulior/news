module github.com/qsoulior/news/newsdata-parser

go 1.22.0

require (
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/qsoulior/news/aggregator v0.0.0-00010101000000-000000000000
	github.com/qsoulior/news/parser v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/redis/go-redis/v9 v9.5.1 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/qsoulior/news/aggregator => ../aggregator

replace github.com/qsoulior/news/parser => ../parser
