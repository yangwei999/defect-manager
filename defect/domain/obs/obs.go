package obs

type OBS interface {
	Upload([]byte) error
}
