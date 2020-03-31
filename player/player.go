package player

// #include "player.h"
// #include "stdlib.h"
// void playFinishCallback(void*);
import "C"
import "unsafe"
import "errors"
import "log"
import "sync"
import "time"

const (
	_ uint32 = iota
	INIT
	PLAY
	PAUSE
	STOP
)

type FinishCallback func(p *PlayerWorker)

type PlayerWorker struct {
	cPlayer  unsafe.Pointer
	cDecoder unsafe.Pointer

	CurrAudiopath string
	PlayerStatus  uint32

	volume float32

	Callback []FinishCallback
	mutex    sync.Mutex
}

type AudioInfo struct {
	Status     uint32
	Pathinfo   string
	Duration   int64
	CurrSecond int64
	Volume     float32
}

//export playFinishCallback
func playFinishCallback(pw unsafe.Pointer) {
	begin := time.Now()
	p := (*PlayerWorker)(unsafe.Pointer(pw))
	p.mutex.Lock()
	p.PlayerStatus = PAUSE
	go func() {
		p.mrReset()
		p.mutex.Unlock()
		for _, f := range p.Callback {
			f(p)
		}
		log.Printf("callback duration:%s\n", time.Since(begin))
	}()
}

func NewPlayerWorker() *PlayerWorker {
	p := &PlayerWorker{
		CurrAudiopath: "",
		PlayerStatus:  INIT,
		cPlayer:       C.malloc(C.sizeof_mr_player),
		cDecoder:      C.malloc(C.sizeof_ma_decoder),
		Callback:      make([]FinishCallback, 0),

		volume: 0.5,
	}

	return p
}

func (p *PlayerWorker) AddCallback(f FinishCallback) {
	p.Callback = append(p.Callback, f)
}

func (p *PlayerWorker) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.mrStop()
	C.free((p.cPlayer))
}

func (p *PlayerWorker) Play(audiopath string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var err error
	switch p.PlayerStatus {
	case INIT:
		err = p.mrInit(audiopath)
	case PLAY:
		if audiopath == p.CurrAudiopath {
			err = errors.New(audiopath + " is playing now")
		} else {
			p.mrStop()
			err = p.mrInit(audiopath)
		}
	case PAUSE:
		if audiopath != p.CurrAudiopath {
			p.mrStop()
			err = p.mrInit(audiopath)
		}
	case STOP:
		err = p.mrInit(audiopath)
	}
	if err == nil {
		p.CurrAudiopath = audiopath
		p.PlayerStatus = PLAY
		return p.mrStart()
	}
	return err
}

func (p *PlayerWorker) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	switch p.PlayerStatus {
	case INIT:
		return
	case PLAY:
		p.mrStop()
	case PAUSE:
		p.mrStop()
	case STOP:
		p.mrStop()
	}
	p.PlayerStatus = STOP
}

func (p *PlayerWorker) Pause() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	switch p.PlayerStatus {
	case INIT:
		return
	case PLAY:
		p.mrPause()
		p.PlayerStatus = PAUSE
	case PAUSE:
		p.mrPause()
		p.PlayerStatus = PAUSE
	case STOP:
		return
	}
}

func (p *PlayerWorker) CurrAudioInfo() (*AudioInfo, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	a := &AudioInfo{
		Status: uint32(p.PlayerStatus),
		Volume: p.volume,
	}

	if p.PlayerStatus == INIT || p.PlayerStatus == STOP {
		return a, nil
	}
	if p.CurrAudiopath == "" {
		return nil, errors.New("CurrAudiopath is empty, but status wrong")
	}
	a.Pathinfo = p.CurrAudiopath
	p.mrCurrAudioinfo(a)
	return a, nil
}

func (p *PlayerWorker) SetVolume(volume float32) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if volume < 0 || p.volume > 1 {
		log.Printf("can't set volume to %f\n", volume)
		return
	}
	p.volume = volume
	if p.PlayerStatus == PLAY || p.PlayerStatus == PAUSE {
		p.mrSetVolume()
	}
}

func (p *PlayerWorker) Seek(second int64) {
	p.mrSeekFrame(second)
}

func (p *PlayerWorker) mrInit(audiopath string) error {
	// init decoder
	err := C.mr_decoder_init_file((*C.ma_decoder)(p.cDecoder), C.CString(audiopath))
	if err != nil {
		return errors.New(C.GoString(err))
	}

	// init device
	err = C.mr_player_init((*C.mr_player)(p.cPlayer), (*C.ma_decoder)(p.cDecoder),
		C.callback(C.playFinishCallback), unsafe.Pointer(p), (C.float)(p.volume))
	if err != nil {
		return errors.New(C.GoString(err))
	}

	log.Printf("init play worker: %s\n", audiopath)
	return nil
}

func (p *PlayerWorker) mrStart() error {
	err := C.mr_player_start((*C.mr_player)(p.cPlayer))
	if err != nil {
		e := C.GoString(err)
		log.Println(e)
		return errors.New(e)
	}
	log.Printf("start play: %s\n", p.CurrAudiopath)
	return nil
}

func (p *PlayerWorker) mrStop() {
	C.mr_player_destory((*C.mr_player)(p.cPlayer))
	log.Printf("stop play: %s\n", p.CurrAudiopath)
}

func (p *PlayerWorker) mrPause() {
	C.mr_player_stop((*C.mr_player)(p.cPlayer))
	log.Printf("pause play: %s\n", p.CurrAudiopath)
}

func (p *PlayerWorker) mrReset() {
	C.mr_player_reset((*C.mr_player)(p.cPlayer))
	log.Printf("reset play: %s\n", p.CurrAudiopath)
}

func (p *PlayerWorker) mrCurrAudioinfo(info *AudioInfo) {
	C.mr_curr_audio_info((*C.mr_player)(p.cPlayer),
		(*C.int64_t)(&info.Duration), (*C.int64_t)(&info.CurrSecond))
}

func (p *PlayerWorker) mrSetVolume() {
	C.mr_player_set_volume((*C.mr_player)(p.cPlayer), (C.float)(p.volume))
}

func (p *PlayerWorker) mrSeekFrame(second int64) {
	C.mr_player_seek_frame((*C.mr_player)(p.cPlayer), (C.int64_t)(second))
}
