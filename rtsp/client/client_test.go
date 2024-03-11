package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHandler(packets []byte) (err error) {
	for _, p := range packets {
		fmt.Println(p)
	}

	return
}

func TestRTSPClient(t *testing.T) {
	rtsp := "rtsp://admin:laonpeople!@10.30.8.159:554/profile2/media.smp"

	client := NewRTSPClient()
	err := client.Open(rtsp, testHandler)
	defer client.Close()
	assert.Nil(t, err)
}
