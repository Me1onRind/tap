package ui

import (
	//"github.com/gizak/termui/v3"
	"context"
	//"log"
	//"tap/backend"
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
		w.op.Println(err)
		return false
	}
	return true
}

func (w *Window) PlayStatus() *server.PlayAudioInfo {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	a, err := w.playerClient.Status(ctx, &server.Empty{})
	if err != nil {
		w.op.Println(err)
		return nil
	}
	return a
}

func (w *Window) ServerIsHealthLive() bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.Ping(ctx, &server.Empty{})
	if err != nil {
		w.op.Println(err)
		return false
	}
	return true
}

func (w *Window) PlayOrPause(audioPath string) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.PlayOrPause(ctx, &server.PlayRequest{
		AudioPath: audioPath,
	})
	if err != nil {
		w.op.Println(err)
	}
}

func (w *Window) SetVolume(volume float32) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.SetVolume(ctx, &server.VolumeRequest{
		Volume: volume,
	})
	if err != nil {
		w.op.Println(err)
		return false
	}
	return true
}

func (w *Window) ListAll() []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := w.playerClient.ListAll(ctx, &server.Empty{})
	if err != nil {
		w.op.Println(err)
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
		w.op.Println(err)
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
		w.op.Println(err)
	}
}

func (w *Window) SeekAudioFile(second int64) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.Seek(ctx, &server.Second{
		Value: second,
	})
	if err != nil {
		w.op.Printf("request seek, value:%d, err:%s\n", second, err.Error())
	}
}

func (w *Window) Provider() *server.ProviderReply {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	p, err := w.playerClient.Provider(ctx, &server.Empty{})
	if err != nil {
		w.op.Println(err)
		return nil
	}
	return p
}

func (w *Window) SetDir(dir string) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := w.playerClient.SetDir(ctx, &server.Dir{Value: dir})
	if err != nil {
		w.op.Println(err)
	}
}

func (w *Window) subscribe() {
	res, _ := w.playerClient.PushInfo(context.Background(), &server.Empty{})
	for {
		info, err := res.Recv()
		if err != nil {
			w.op.Println(err)
			res.CloseSend()
			return
		}

		if info != nil {
			w.playStatus.Notify(info)
			w.audioList.NotifyAudioPathChange(info.Pathinfo)
		}
	}
}
