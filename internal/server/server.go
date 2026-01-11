package server

type Server struct {
	Port string
}

func NewServer() *Server {
	server := &Server{
		Port: "8080",
	}
	return server
}

func StartServer() {
	//start server and pass params into redis
}
