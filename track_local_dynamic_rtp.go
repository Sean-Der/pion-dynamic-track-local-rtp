package main

import (
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type trackLocalDynamicRTP struct {
	ssrc         webrtc.SSRC
	writeStream  webrtc.TrackLocalWriter
	bindCallback func(ctx webrtc.TrackLocalContext) webrtc.RTPCodecParameters

	payloadType uint8
	kind        webrtc.RTPCodecType

	id, rid, streamID string
}

func (t *trackLocalDynamicRTP) Bind(ctx webrtc.TrackLocalContext) (webrtc.RTPCodecParameters, error) {
	selectedCodec := t.bindCallback(ctx)

	t.ssrc = ctx.SSRC()
	t.writeStream = ctx.WriteStream()
	t.payloadType = uint8(selectedCodec.PayloadType)

	return selectedCodec, nil
}

func (t *trackLocalDynamicRTP) Unbind(webrtc.TrackLocalContext) error {
	return nil
}

func (t *trackLocalDynamicRTP) WriteRTP(p *rtp.Packet) error {
	p.Header.SSRC = uint32(t.ssrc)
	p.Header.PayloadType = t.payloadType

	_, err := t.writeStream.WriteRTP(&p.Header, p.Payload)
	return err
}

func (t *trackLocalDynamicRTP) ID() string       { return t.id }
func (t *trackLocalDynamicRTP) RID() string      { return t.rid }
func (t *trackLocalDynamicRTP) StreamID() string { return t.streamID }
func (t *trackLocalDynamicRTP) Kind() webrtc.RTPCodecType {
	return webrtc.RTPCodecTypeAudio
}
