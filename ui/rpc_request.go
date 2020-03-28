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

func (w *Window) PlayStatus() *server.PlayAudioInfo {
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

func (w *Window) PlayOrPause(name string) *server.PlayAudioInfo {
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

func (w *Window) SetVolume(volume float32) bool {
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

func (w *Window) ListAll() []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.ListAll(ctx, &server.Empty{})
	if err != nil {
		return []string{}
	}
	return res.GetNames()
}

func (w *Window) Search(input string) []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.Search(ctx, &server.SearchRequest{
		Input: input,
	})
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return res.GetNames()
}

func (w *Window) ChangeLoopModel(mode uint32) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.SetPlayMode(ctx, &server.PlayMode{
		Mode: mode,
	})
	if err != nil {
		log.Println(err)
	}
}

func (w *Window) subscribe() {
	res, _ := w.playerClient.PushInfo(context.Background(), &server.Empty{})
	for {
		info, err := res.Recv()
		if err != nil {
			log.Println(err)
			res.CloseSend()
			return
		}

		w.ps.Notify(info)
		w.al.NotifyPlayNameChange(info.Name)
	}
}

func (w *Window) forward() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	<-w.rewindOrForwardChan
	_, err := w.playerClient.Forward(ctx, &server.Second{
		Value: 2,
	})
	if err != nil {
		log.Println(err)
	}
	w.rewindOrForwardChan <- struct{}{}
}

func (w *Window) rewind() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	<-w.rewindOrForwardChan
	_, err := w.playerClient.Rewind(ctx, &server.Second{
		Value: 2,
	})
	if err != nil {
		log.Println(err)
	}
	w.rewindOrForwardChan <- struct{}{}
}
