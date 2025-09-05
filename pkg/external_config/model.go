package external_config

type Config struct {
	PostgreSQL PostgreSQL `json:"postgre_sql"`
}

type PostgreSQL struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Port    string `json:"port"`
	DbName  string `json:"db_name"`
	Pass    string `json:"pass"`
	SslMode string `json:"ssl_mode"`
}
