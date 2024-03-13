package client

import (
	"errors"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/pion/rtp"
	log "github.com/sirupsen/logrus"
)

type RTSPClient struct {
	client        *gortsplib.Client
	url           *base.URL
	packetHandler func(packets []byte) error
}

func NewRTSPClient() *RTSPClient {
	return &RTSPClient{}
}

// @param Host string: rtsp 경로
// @param handler func(packets []byte) error: 패킷 전달 핸들러 함수
func (rc *RTSPClient) Open(Host string, handler func(packets []byte) error) (err error) {
	if handler == nil {
		return errors.New("rtsp open fail: packet handler is nil")
	}
	rc.packetHandler = handler

	rc.client = &gortsplib.Client{}
	rc.url, err = base.ParseURL(Host)
	if err != nil {
		return errors.New("rtsp base url parse fail[" + rc.url.Host + "]: " + err.Error())
	}

	return
}

func (rc *RTSPClient) Close() {
	rc.client.Close()
	rc.packetHandler = nil
}

// @param interval int: 실패 시, 재시도 연결 시도하는 간격(초) [0인 경우 재연결 시도 하지 않음]
func (rc *RTSPClient) Run() (err error) {

	err = rc.client.Start(rc.url.Scheme, rc.url.Host)
	if err != nil {
		return errors.New("rtsp start fail[" + rc.url.Host + "]: " + err.Error())
	}

	// find available medias
	desc, _, err := rc.client.Describe(rc.url)
	if err != nil {
		return errors.New("can not find available medias-> " + err.Error())
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		return errors.New("can not find h264 format")
	}

	// setup RTP/H264 -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		return errors.New("create decorder fail-> " + err.Error())
	}

	// setup a single media
	_, err = rc.client.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		return errors.New("rtsp setup fail-> " + err.Error())
	}

	rc.client.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// decode timestamp
		_, ok := rc.client.PacketPTS(medi, pkt)
		if !ok {
			log.Warn("waiting for timestamp")
		}

		// extract access units from RTP packets
		packets, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Warn("extract fail: ", err)
			}
		}

		for _, packet := range packets {
			packet = append([]uint8{0x00, 0x00, 0x00, 0x01}, packet...)
			if rc.packetHandler != nil {
				rc.packetHandler(packet)
			}
		}
	})

	_, err = rc.client.Play(nil)
	if err != nil {
		log.Warn("rtsp play fail[" + rc.url.Host + "]: " + err.Error())
	}

	err = rc.client.Wait()
	if err != nil {
		log.Info("rtsp [" + rc.url.Host + "]: " + err.Error())
	}

	rc.Run()

	return
}
