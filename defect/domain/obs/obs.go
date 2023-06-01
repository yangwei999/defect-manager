package obs

type OBS interface {
	Upload(string) error
}
