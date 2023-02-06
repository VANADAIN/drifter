package nat

import "context"

func checkNat() bool {
	ctx := context.Background()
	router, _ := PickRouterClient(ctx)

	// no nat
	if router == nil {
		return false
	} else {
		return true
	}
}

func createPortForwarding(ctx context.Context) {
	err := GetIPAndForwardPort(ctx)
	// panic if there is NAT, but we can connect behind it
	if err != nil {
		panic(err)
	}
}
