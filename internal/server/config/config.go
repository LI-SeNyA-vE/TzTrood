package config

type Server struct {
	Redis             Redis  `json:"Redis"`
	HttpPort          string `json:"HttpPort"`
	ServiceHumanAgent string `json:"ServiceHumanAgent"`
	NLPAddress        string `json:"NLPAddress"`
}

type Redis struct {
	Addr     string `json:"Addr"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

func NewMock() *Server {
	return &Server{
		Redis: Redis{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		HttpPort:          ":8080",
		ServiceHumanAgent: "http://localhost:9000/human",
		NLPAddress:        "http://localhost:8000/analyze",
	}
}
