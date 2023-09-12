package flags

import (
	"flag"
	"os"
)

var ServAddress string
var BaseURL string

func InitFlags() {
	flag.StringVar(&ServAddress, "a", ":8080", "Request URL")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Response URL")
	flag.Parse()

	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		ServAddress = envServAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		BaseURL = envBaseURL
	}
}
