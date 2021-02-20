package ip

import (
	"testing"
)

func TestFetchIP(t *testing.T) {
	t.Log(LocalIP())
	t.Log(IPV4())
	t.Log(IPV6())
}
