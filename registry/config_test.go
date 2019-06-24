package registry

import "testing"

func TestReadConf(t *testing.T) {
	data, err := ReadConf("redis/config/pwd")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}