package dp

type arch string

type Arch interface {
	String() string
}

func NewArch(s string) Arch {
	return arch(s)
}

func (a arch) String() string {
	return string(a)
}
