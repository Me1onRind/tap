package rpc_client

import (
	"context"
	"log"
	"os"
	"tap/server"
	"tap/server/guider"
	"time"
)

var (
	playerClient server.PlayClient
	op           Output
)

const (
	_TIME_OUT = 500
)

func init() {
	op = log.New(os.Stdout, "", log.LstdFlags)
}

type Output interface {
	Println(values ...interface{})
	Printf(format string, values ...interface{})
}

func SetRpcClient(p server.PlayClient) {
	playerClient = p
}

func SetOutput(o Output) {
	op = o
}

func SetLocalProvider(dirs []string) bool {
	request := server.LocalProvider{
		Dirs: dirs,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.SetLocalProvider(ctx, &request)
	if err != nil {
		op.Println(err)
		return false
	}
	return true
}

func PlayStatus() *server.PlayAudioInfo {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	a, err := playerClient.Status(ctx, &server.Empty{})
	if err != nil {
		op.Println(err)
		return nil
	}
	return a
}

func ServerIsHealthLive() bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.Ping(ctx, &server.Empty{})
	if err != nil {
		op.Println(err)
		return false
	}
	return true
}

func PlayOrPause(audioPath string) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.PlayOrPause(ctx, &server.PlayRequest{
		AudioPath: audioPath,
	})
	if err != nil {
		op.Println(err)
	}
}

func SetVolume(volume float32) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.SetVolume(ctx, &server.VolumeRequest{
		Volume: volume,
	})
	if err != nil {
		op.Println(err)
		return false
	}
	return true
}

func ListAll() []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := playerClient.ListAll(ctx, &server.Empty{})
	if err != nil {
		op.Println(err)
		return []string{}
	}
	return res.GetNames()
}

func Search(input string) []string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	res, err := playerClient.Search(ctx, &server.SearchRequest{
		Input: input,
	})
	if err != nil {
		op.Println(err)
		return []string{}
	}
	return res.GetNames()
}

func ChangeLoopModel(mode guider.Mode) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.SetPlayMode(ctx, &server.PlayMode{
		Mode: uint32(mode),
	})
	if err != nil {
		op.Println(err)
	}
}

func SeekAudioFile(second int64) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.Seek(ctx, &server.Second{
		Value: second,
	})
	if err != nil {
		op.Printf("request seek, value:%d, err:%s\n", second, err.Error())
	}
}

func Provider() *server.ProviderReply {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	p, err := playerClient.Provider(ctx, &server.Empty{})
	if err != nil {
		op.Println(err)
		return nil
	}
	return p
}

func SetDir(dir string) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*_TIME_OUT)
	defer cancel()
	_, err := playerClient.SetDir(ctx, &server.Dir{Value: dir})
	if err != nil {
		op.Println(err)
	}
}

func Subscribe(f func(info *server.PlayAudioInfo)) {
	res, _ := playerClient.PushInfo(context.Background(), &server.Empty{})
	for {
		info, err := res.Recv()
		if err != nil {
			op.Println(err)
			res.CloseSend()
			return
		}

		if info != nil {
			f(info)
		}
	}
}
