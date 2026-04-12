package mongo

// Config database configuration loader.
type Config struct {
	DSN      string `env:"DSN,default=mongodb://localhost:27017/database?retryWrites=false" json:"dsn"`
	Username string `env:"USERNAME,default=username"                                        json:"username"`
	Password string `env:"PASSWORD,default=password"                                        json:"password"`
}
