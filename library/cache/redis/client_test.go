package redis

import (
	"fmt"
	"testing"
)

func TestNewClusterClient(t *testing.T) {
	a := " "
	fmt.Println(len(a))
	fmt.Println(NewClusterClient())
}