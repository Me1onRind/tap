package ui

import (
	//"github.com/gizak/termui/v3"
	"context"
	"log"
	"tap/server"
	"time"
)

const (
	_TIME_OUT = 500
)

func (w *Window) ChceckPlayStatus() *server.PlayAudioInfo {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	a, err := w.playerClient.Status(ctx, &server.Empty{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return a
}

func (w *Window) ServerIsHealthLive() bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.Ping(ctx, &server.Empty{})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (w *Window) SetLocalProvider(dirs []string) bool {
	request := server.LocalProvider{
		Dirs: dirs,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.SetLocalProvider(ctx, &request)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (w *Window) PlayOrPause(index int) *server.PlayAudioInfo {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.PlayOrPause(ctx, &server.PlayRequest{
		Index: uint32(index),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return res
}
