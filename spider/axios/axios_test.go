package axios

import "testing"

func TestAxios(t *testing.T) {
	for i := 0; i <= 5; i++ {
		GetClient().Request.Get("https://baidu.com")
	}
}
