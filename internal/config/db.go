package config

type PgConn struct {
	Host     string `split_words:"true" required:"true"`
	Port     uint16 `split_words:"true" required:"true"`
	Username string `split_words:"true" required:"true"`
	Password string `split_words:"true" required:"true"`
	DBName   string `split_words:"true" required:"true"`
	SSLMode  string `split_words:"true" default:"disable"`
}
