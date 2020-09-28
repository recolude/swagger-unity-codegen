package security

// OAuth is a guard on route that requires oauth to be used
type OAuth struct {
	Name string
	In   string
}
