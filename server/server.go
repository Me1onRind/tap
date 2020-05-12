package server

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
	"tap/backend/local"
	"tap/player"
)

const (
	UNIX_SOCK_FILE = "/tmp/tap.sock"
)

type Play struct {
	pushChan chan *PlayAudioInfo
}

var (
	ps *Play
)

func init() {
	ps = &Play{}
}

func (server *Play) Ping(ctx context.Context, empty *Empty) (*Empty, error) {
	return &Empty{}, nil
}

func (server *Play) PlayOrPause(ctx context.Context, request *PlayRequest) (*Empty, error) {
	if authInfo, err := m.PlayOrPause(request.AudioPath); err != nil {
		return nil, err
	} else {
		server.push(authInfo)
		return &Empty{}, nil
	}
}

func (server *Play) Status(ctx context.Context, empty *Empty) (*PlayAudioInfo, error) {
	if authInfo, err := m.PlayWorker.CurrAudioInfo(); err != nil {
		return nil, err
	} else {
		return fommatPlayAudioInfo(authInfo), nil
	}
}

func (server *Play) SetVolume(ctx context.Context, volume *VolumeRequest) (*Empty, error) {
	m.PlayWorker.SetVolume(volume.Volume)
	return &Empty{}, nil
}

func (server *Play) Seek(ctx context.Context, second *Second) (*Empty, error) {
	m.PlayWorker.Seek(second.Value)
	return &Empty{}, nil
}

func (server *Play) Stop(ctx context.Context, empty *Empty) (*Empty, error) {
	m.PlayWorker.Stop()
	return &Empty{}, nil
}

func (server *Play) SetPlayMode(ctx context.Context, playMode *PlayMode) (*Empty, error) {
	m.SetPlayMode(playMode.Mode)
	return &Empty{}, nil
}

func (server *Play) ListAll(ctx context.Context, empty *Empty) (*QueryReplay, error) {
	all, err := m.ListAll()
	if err != nil {
		return nil, err
	}
	return &QueryReplay{Names: all}, nil
}

func (server *Play) Search(ctx context.Context, request *SearchRequest) (*QueryReplay, error) {
	if len(request.Input) == 0 {
		return server.ListAll(ctx, nil)
	}
	all, err := m.Search(request.Input)
	if err != nil {
		return nil, err
	}

	return &QueryReplay{Names: all}, nil
}

func (server *Play) SetLocalProvider(ctx context.Context,
	localPrivoder *LocalProvider) (*Empty, error) {
	if len(localPrivoder.Dirs) == 0 {
		return nil, errors.New("Dirs can't be length 0")
	}
	m.Init(local.NewLocalProvider(), LocalPrivoderType, localPrivoder.Dirs, 1)
	return &Empty{}, nil
}

func (server *Play) Provider(ctx context.Context, empty *Empty) (*ProviderReply, error) {
	return &ProviderReply{
		ProviderType: int32(m.ProviderType),
		Name:         "",
		CurrDir:      m.CurrDir,
		Dirs:         m.Dirs,
	}, nil
}

func (server *Play) SetDir(ctx context.Context, dir *Dir) (*Empty, error) {
	return &Empty{}, nil
}

func (server *Play) PushInfo(empty *Empty, res Play_PushInfoServer) error {
	ch := make(chan *PlayAudioInfo, 100)
	server.pushChan = ch
	for {
		select {
		case info := <-ch:
			res.Send(info)
		case <-res.Context().Done():
			server.pushChan = nil
			return nil
		}
	}
}

func (server *Play) push(info *player.AudioInfo) {
	if server.pushChan != nil {
		select {
		case server.pushChan <- fommatPlayAudioInfo(info):
		}
	}
}

func fommatPlayAudioInfo(authInfo *player.AudioInfo) *PlayAudioInfo {
	return &PlayAudioInfo{
		Status:   uint32(authInfo.Status),
		Duration: authInfo.Duration,
		Curr:     authInfo.CurrSecond,
		Pathinfo: authInfo.Pathinfo,
		Volume:   authInfo.Volume,
		Name:     pathToName(authInfo.Pathinfo),
		Mode:     m.PlayMode,
	}
}

func pathToName(pathinfo string) string {
	if pathinfo == "" {
		return ""
	}
	return filepath.Base(pathinfo)
}

func RunServer() {
	os.Remove(UNIX_SOCK_FILE)

	l, err := net.Listen("unix", UNIX_SOCK_FILE)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	RegisterPlayServer(s, ps)
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("start rpc server")
	s.Serve(l)
}
