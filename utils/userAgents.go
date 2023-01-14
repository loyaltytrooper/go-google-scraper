package utils

import (
	"math/rand"
	"time"
)

var userAgents = []string{}

func randomUserAgents() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}
