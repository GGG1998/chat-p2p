package network

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomPort() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	port := fmt.Sprintf(":800%d", random.Intn(9))

	return port
}
