package configs

type Settings struct {
	MongoConnString string `env:"MONGO_CONN_STRING"`
}
