package utils

import "testing"

func TestGetLocalIP(t *testing.T) {
	localIP := GetLocalIP()
	if localIP == DefaultIp {
		t.Errorf("获得内网IP失败 %s", localIP)
	} else {
		t.Log("内网IP", localIP)
	}
}
