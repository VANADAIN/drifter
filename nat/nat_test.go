package nat

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {
	ip := GetLocalIP()
	fmt.Println(ip.String())
}
