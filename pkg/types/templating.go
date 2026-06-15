package types

const (
	PageIndex Page = iota
	PageSettings
)

type (
	Page int
)

func (p Page) String() string {
	return [...]string{
		"index.html",
		"settings.html",
	}[p]
}
