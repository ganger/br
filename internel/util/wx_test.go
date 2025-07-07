package util

import (
	"testing"
)

func TestPushWX(t *testing.T) {

	PushWX("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=59b26a0b-9ae2-4e76-b6c0-cb38b1a0c3f9", "test message")

}
