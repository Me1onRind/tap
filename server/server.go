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
)

const (
	UNIX_SOCK_FILE = "/tmp/tap.sock"
)

type PlayServer struct {
}

func (server *PlayServer) Ping(ctx context.Context, empty *Empty) (*Empty, error) {
	return &Empty{}, nil
}

func (server *PlayServer) PlayOrPause(ctx context.Context, request *PlayRequest) (*PlayAudioInfo, error) {
	if authinfo, err := PlayOrPause(request.Name); err != nil {
		return nil, err
	} else {
		return &PlayAudioInfo{
			Status:     uint32(authinfo.Status),
			SampleRate: authinfo.SampleRate,
			Duration:   authinfo.Duration,
			Curr:       authinfo.CurrSecond,
			Pathinfo:   authinfo.Pathinfo,
			Volume:     authinfo.Volume,
			Name:       filepath.Base(authinfo.Pathinfo),
		}, nil
	}
}

func (server *PlayServer) Status(ctx context.Context, empty *Empty) (*PlayAudioInfo, error) {
	if authinfo, err := Status(); err != nil {
		return nil, err
	} else {
		return &PlayAudioInfo{
			Status:     uint32(authinfo.Status),
			SampleRate: authinfo.SampleRate,
			Duration:   authinfo.Duration,
			Curr:       authinfo.CurrSecond,
			Pathinfo:   authinfo.Pathinfo,
			Volume:     authinfo.Volume,
			Name:       filepath.Base(authinfo.Pathinfo),
		}, nil
	}
}

func (server *PlayServer) SetVolume(ctx context.Context, volume *VolumeRequest) (*Empty, error) {
	SetVolume(volume.Volume)
	return &Empty{}, nil
}

func (server *PlayServer) Stop(ctx context.Context, empty *Empty) (*Empty, error) {
	Stop()
	return &Empty{}, nil
}

func (server *PlayServer) ListAll(ctx context.Context, empty *Empty) (*QueryReplay, error) {
	all, err := ListAll()
	if err != nil {
		return nil, err
	}
	return &QueryReplay{Names: all}, nil
}

func (server *PlayServer) Search(ctx context.Context, request *SearchRequest) (*QueryReplay, error) {
	if len(request.Input) == 0 {
		return server.ListAll(ctx, nil)
	}
	all, err := Search(request.Input)
	if err != nil {
		return nil, err
	}

	return &QueryReplay{Names: all}, nil
}

func (server *PlayServer) SetLocalProvider(ctx context.Context,
	localPrivoder *LocalProvider) (*Empty, error) {
	if len(localPrivoder.Dirs) == 0 {
		return nil, errors.New("Dirs can't be length 0")
	}
	provider = local.NewLocalProvider(localPrivoder.Dirs)
	return &Empty{}, nil
}

func (server *PlayServer) Provider(ctx context.Context, empty *Empty) (*ProviderReply, error) {
	return &ProviderReply{
		ProviderType: int32(providerType),
		Name:         providerName,
	}, nil
}

func RunServer() {
	os.Remove(UNIX_SOCK_FILE)

	l, err := net.Listen("unix", UNIX_SOCK_FILE)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	s := grpc.NewServer()
	RegisterPlayServerServer(s, &PlayServer{})
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("start rpc server")
	s.Serve(l)
}
