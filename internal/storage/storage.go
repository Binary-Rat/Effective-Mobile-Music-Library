package storage

type Interface interface {
	Song() (string, error)
	Text() (string, error)
	Delete() error
	AddSong() error
	ChangeSong() error
}
