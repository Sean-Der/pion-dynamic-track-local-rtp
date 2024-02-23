package main

import (
	"fmt"

	"github.com/pion/webrtc/v4"
)

func main() {
	offerer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assertNilErr(err)

	track := &trackLocalDynamicRTP{
		kind: webrtc.RTPCodecTypeAudio,
		bindCallback: func(ctx webrtc.TrackLocalContext) webrtc.RTPCodecParameters {
			codecs := ctx.CodecParameters()
			if len(codecs) == 0 {
				panic("No codecs found")
			}

			fmt.Printf("Selecting codec %s \n", codecs[0].MimeType)
			return codecs[0]
		},
	}

	_, err = offerer.AddTrack(track)
	assertNilErr(err)

	// Create a PeerConnection that only accepts PCMU
	m := &webrtc.MediaEngine{}
	assertNilErr(m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypePCMU,
			ClockRate:    8000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 0,
	}, webrtc.RTPCodecTypeAudio))

	answerer, err := webrtc.NewAPI(webrtc.WithMediaEngine(m)).NewPeerConnection(webrtc.Configuration{})
	assertNilErr(err)

	offer, err := offerer.CreateOffer(nil)
	assertNilErr(err)

	assertNilErr(offerer.SetLocalDescription(offer))
	assertNilErr(answerer.SetRemoteDescription(offer))

	answer, err := answerer.CreateAnswer(nil)
	assertNilErr(err)

	assertNilErr(answerer.SetLocalDescription(answer))
	assertNilErr(offerer.SetRemoteDescription(answer))

	fmt.Println("Session negotiated without error!")
}

func assertNilErr(err error) {
	if err != nil {
		panic(err)
	}
}

// func
