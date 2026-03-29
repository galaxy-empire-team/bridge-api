package config

type GRPCServer struct {
	Endpoint string `split_words:"true" required:"true"`
}
