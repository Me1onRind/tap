package guider

type Guider interface {
	NextAudioPath() string
	PreAudioPath() string
}
