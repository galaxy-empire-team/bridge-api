package config

type HTTPServer struct {
	Endpoint string `split_words:"true" required:"true"`
}
