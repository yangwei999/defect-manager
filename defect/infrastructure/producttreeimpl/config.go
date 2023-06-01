package producttreeimpl

type Config struct {
	Token  string `json:"token"        required:"true"`
	PkgRPM PkgRPM `json:"pkg_rpm"      required:"true"`
}

type PkgRPM struct {
	Org        string `json:"org"         required:"true"`
	Repo       string `json:"repo"        required:"repo"`
	PathPrefix string `json:"path_prefix" required:"true"`
	Branch     string `json:"branch"      required:"true"`
}
