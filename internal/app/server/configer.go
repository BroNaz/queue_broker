package server

type Configer struct {
	BindAddr string `toml:"bint_addr"`
	LogLevel string `toml:"log_level"`
}

func NewConfiger() *Configer {
	return &Configer{
		BindAddr: ":8080",
		LogLevel: "debag",
	}
}
