package server

import (
	"errors"
	"strconv"
	"sync"

	"github.com/pion/rtp"
	log "github.com/sirupsen/logrus"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
)

type RTSPServer struct {
	port      int
	s         *gortsplib.Server
	mutex     sync.Mutex
	stream    *gortsplib.ServerStream
	publisher *gortsplib.ServerSession
	isConn    chan bool
}

func NewRTSPSever() *RTSPServer {
	return &RTSPServer{}
}

func (rs *RTSPServer) Open(port int, isConn chan bool) (err error) {
	if port <= 0 || port > 655365 {
		return errors.New("rtsp server open fail: port is invalid")
	}

	if isConn == nil {
		return errors.New("rtsp server open fail: connection check channel is nil")
	}

	rs.port = port
	rs.isConn = isConn

	return
}

func (rs *RTSPServer) Close() {
	rs.s.Close()
}

func (rs *RTSPServer) Run() (err error) {
	p := strconv.Itoa(rs.port)
	rs.s = &gortsplib.Server{
		Handler:     rs,
		RTSPAddress: ":" + p,
	}

	if err = rs.s.StartAndWait(); err != nil {
		return
	}

	return
}

// called when a connection is opened.
func (sh *RTSPServer) OnConnOpen(ctx *gortsplib.ServerHandlerOnConnOpenCtx) {
}

// called when a connection is closed.
func (sh *RTSPServer) OnConnClose(ctx *gortsplib.ServerHandlerOnConnCloseCtx) {
}

// called when a session is opened.
func (sh *RTSPServer) OnSessionOpen(ctx *gortsplib.ServerHandlerOnSessionOpenCtx) {
}

// called when a session is closed.
func (sh *RTSPServer) OnSessionClose(ctx *gortsplib.ServerHandlerOnSessionCloseCtx) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	// if the session is the publisher,
	// close the stream and disconnect any reader.
	if sh.stream != nil && ctx.Session == sh.publisher {
		sh.stream.Close()
		sh.stream = nil
	}
}

// called when receiving a DESCRIBE request.
func (sh *RTSPServer) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	// no one is publishing yet
	if sh.stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	// send medias that are being published to the client
	return &base.Response{
		StatusCode: base.StatusOK,
	}, sh.stream, nil
}

// called when receiving an ANNOUNCE request.
func (sh *RTSPServer) OnAnnounce(ctx *gortsplib.ServerHandlerOnAnnounceCtx) (*base.Response, error) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	// disconnect existing publisher
	if sh.stream != nil {
		sh.stream.Close()
		sh.publisher.Close()
	}

	// create the stream and save the publisher
	sh.stream = gortsplib.NewServerStream(sh.s, ctx.Description)
	sh.publisher = ctx.Session

	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}

// called when receiving a SETUP request.
func (sh *RTSPServer) OnSetup(ctx *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	// no one is publishing yet
	if sh.stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, sh.stream, nil
}

// called when receiving a PLAY request.
func (sh *RTSPServer) OnPlay(ctx *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	log.Info("play request: ", ctx.Conn.NetConn().LocalAddr().String())

	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}

// called when receiving a RECORD request.
func (sh *RTSPServer) OnRecord(ctx *gortsplib.ServerHandlerOnRecordCtx) (*base.Response, error) {
	// called when receiving a RTP packet
	ctx.Session.OnPacketRTPAny(func(medi *description.Media, forma format.Format, pkt *rtp.Packet) {
		// route the RTP packet to all readers
		sh.stream.WritePacketRTP(medi, pkt)
	})

	if sh.isConn != nil {
		sh.isConn <- true
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}
