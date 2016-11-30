package reconf

import "testing"

func Test_Conf(t *testing.T) {
	err := Init("weixin_article")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("OK")
}
