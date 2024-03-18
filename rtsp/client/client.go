package client

import (
	"errors"
	"time"

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
	retryCh       chan error
}

func NewRTSPClient() *RTSPClient {
	return &RTSPClient{
		retryCh: make(chan error, 2),
	}
}

func (rc *RTSPClient) Open(Host string, handler func(packets []byte) error) (err error) {
	if handler == nil {
		return errors.New("rtsp open fail: packet handler is nil")
	}
	rc.packetHandler = handler

	rc.url, err = base.ParseURL(Host)
	if err != nil {
		return errors.New("rtsp base url parse fail[" + rc.url.Host + "]: " + err.Error())
	}

	go rc.retry()

	return
}

func (rc *RTSPClient) Close() {
	rc.client.Close()
	rc.packetHandler = nil
}

func (rc *RTSPClient) Run() {
	rc.client = &gortsplib.Client{}
	err := rc.client.Start(rc.url.Scheme, rc.url.Host)
	if err != nil {
		rc.retryCh <- err
		return
	}

	// find available medias
	desc, _, err := rc.client.Describe(rc.url)
	if err != nil {
		rc.retryCh <- err
		return
	}

	// find the H264 media and format
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		rc.retryCh <- errors.New("can not find h264 format")
		return
	}

	// setup RTP/H264 -> H264 decoder
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		rc.retryCh <- err
		return
	}

	// setup a single media
	_, err = rc.client.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		rc.retryCh <- err
		return
	}

	rc.client.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// decode timestamp
		rc.client.PacketPTS(medi, pkt)

		// extract access units from RTP packets
		packets, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				rc.retryCh <- err
				return
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
		rc.retryCh <- err
		return
	}

	// log.Info("rtsp connect success [" + rc.url.Host + "]")
	err = rc.client.Wait()
	if err != nil {
		rc.retryCh <- err
		return
	}
}

func (rc *RTSPClient) retry() (err error) {
	for err := range rc.retryCh {
		log.Info("retry to connection rtsp ["+rc.url.Host+"]: ", err)
		rc.client.Close()
		rc.client = nil
		rc.Run()
		time.Sleep(time.Second * 1)
	}

	return
}
