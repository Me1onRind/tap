package player

// #include "player.h"
// #include "stdlib.h"
// void playFinishCallback(void*);
import "C"
import "unsafe"
import "errors"
import "log"
import "sync"

const (
	_ uint32 = iota
	INIT
	PLAY
	PAUSE
	STOP
)

type FinishCallback func()

type PlayerWorker struct {
	cPlayer  unsafe.Pointer
	cDecoder unsafe.Pointer

	currAudiopath string
	playerStatus  uint32

	Callback FinishCallback
	mutex    sync.Mutex
}

type AudioInfo struct {
	Status     uint32
	Pathinfo   string
	Duration   uint32
	CurrSecond uint32
	SampleRate uint32
}

//export playFinishCallback
func playFinishCallback(pw unsafe.Pointer) {
	p := (*PlayerWorker)(unsafe.Pointer(pw))
	p.Callback()
	p.playerStatus = PAUSE
	go func() {
		p.mrReset()
	}()
	log.Println("finish play end callback")
}

func NewPlayerWorker() *PlayerWorker {
	p := &PlayerWorker{
		currAudiopath: "",
		playerStatus:  INIT,
		cPlayer:       C.malloc(C.sizeof_mr_player),
		cDecoder:      C.malloc(C.sizeof_ma_decoder),
	}

	p.Callback = func() {
	}
	return p
}

func (p *PlayerWorker) Close() {
	p.mrStop()
	C.free((p.cPlayer))
}

func (p *PlayerWorker) Play(audiopath string) error {
	log.Printf("will play %s, now status is %d\n", audiopath, p.playerStatus)
	var err error
	switch p.playerStatus {
	case INIT:
		err = p.mrInit(audiopath)
	case PLAY:
		if audiopath == p.currAudiopath {
			err = errors.New(audiopath + " is playing now")
		} else {
			p.mrStop()
			err = p.mrInit(audiopath)
		}
	case PAUSE:
		if audiopath != p.currAudiopath {
			p.mrStop()
			err = p.mrInit(audiopath)
		}
	case STOP:
		err = p.mrInit(audiopath)
	}
	if err == nil {
		p.currAudiopath = audiopath
		p.playerStatus = PLAY
		return p.mrStart()
	}
	return err
}

func (p *PlayerWorker) Stop() {
	switch p.playerStatus {
	case INIT:
		return
	case PLAY:
		p.mrStop()
	case PAUSE:
		p.mrStop()
	case STOP:
		p.mrStop()
	}
	p.playerStatus = STOP
}

func (p *PlayerWorker) Pause() {
	switch p.playerStatus {
	case INIT:
		return
	case PLAY:
		p.mrPause()
		p.playerStatus = PAUSE
	case PAUSE:
		p.mrPause()
		p.playerStatus = PAUSE
	case STOP:
		return
	}
}

func (p *PlayerWorker) CurrAudioInfo() (*AudioInfo, error) {
	a := &AudioInfo{
		Status: uint32(p.playerStatus),
	}
	if p.playerStatus == INIT || p.playerStatus == STOP {
		return a, nil
	}
	if p.currAudiopath == "" {
		return nil, errors.New("currAudiopath is empty, but status wrong")
	}
	a.Pathinfo = p.currAudiopath
	p.mrCurrAudioinfo(a)
	return a, nil
}

func (p *PlayerWorker) mrInit(audiopath string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// init decoder
	err := C.mr_decoder_init_file((*C.ma_decoder)(p.cDecoder), C.CString(audiopath))
	if err != nil {
		return errors.New(C.GoString(err))
	}

	// init device
	err = C.mr_player_init((*C.mr_player)(p.cPlayer), (*C.ma_decoder)(p.cDecoder),
		C.callback(C.playFinishCallback), unsafe.Pointer(p))
	if err != nil {
		return errors.New(C.GoString(err))
	}

	log.Printf("init play worker: %s\n", audiopath)
	return nil
}

func (p *PlayerWorker) mrStart() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	err := C.mr_player_start((*C.mr_player)(p.cPlayer))
	if err != nil {
		log.Println(err)
		return errors.New(C.GoString(err))
	}
	log.Printf("start play: %s\n", p.currAudiopath)
	return nil
}

func (p *PlayerWorker) mrStop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	C.mr_player_destory((*C.mr_player)(p.cPlayer))
	log.Printf("stop play: %s\n", p.currAudiopath)
}

func (p *PlayerWorker) mrPause() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	C.mr_player_stop((*C.mr_player)(p.cPlayer))
	log.Printf("pause play: %s\n", p.currAudiopath)
}

func (p *PlayerWorker) mrReset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	C.mr_player_reset((*C.mr_player)(p.cPlayer))
	log.Printf("reset play: %s\n", p.currAudiopath)
}

func (p *PlayerWorker) mrCurrAudioinfo(info *AudioInfo) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	C.mr_curr_audio_info((*C.mr_player)(p.cPlayer),
		(*C.uint32_t)(&info.Duration), (*C.uint32_t)(&info.CurrSecond),
		(*C.uint32_t)(&info.SampleRate))
}
