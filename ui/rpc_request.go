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

func (w *Window) chceckPlayStatus() *server.PlayAudioInfo {
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

func (w *Window) playOrPause(name string) *server.PlayAudioInfo {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.PlayOrPause(ctx, &server.PlayRequest{
		Name: name,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return res
}

func (w *Window) setVolume(volume float32) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.SetVolume(ctx, &server.VolumeRequest{
		Volume: volume,
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (w *Window) listAll() []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.ListAll(ctx, &server.Empty{})
	if err != nil {
		return []string{}
	}
	return res.GetNames()
}

func (w *Window) search(input string) []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.Search(ctx, &server.SearchRequest{
		Input: input,
	})
	if err != nil {
		return []string{}
	}
	return res.GetNames()
}
