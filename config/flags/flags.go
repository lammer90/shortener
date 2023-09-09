package flags

import "flag"

var A = flag.String("a", ":8080", "Request URL")
var B = flag.String("b", "http://localhost:8080", "Response URL")

func InitFlags() {
	flag.Parse()
}
