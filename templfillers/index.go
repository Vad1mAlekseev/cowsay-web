package templfillers

type Index struct {
	Title   string
	Figures []string
}

func NewIndex(figures *[]string) Index {
	return Index{"Cowsay Web!", *figures}
}

func (Index) Fill(templ string) (string, error) {
	return "", nil
}
