package vars

type Types string

const (
	LoginPassword Types = "Login password data"
	Text          Types = "Text data"
	Card          Types = "Bank card data"
)

func (t Types) ToString() string {
	return string(t)
}
