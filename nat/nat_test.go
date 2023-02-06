package nat

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {
	ip := getLocalIP()
	fmt.Println(ip.String())
}
