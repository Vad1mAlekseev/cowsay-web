package templfillers

type Figure struct {
	Title  string
	figure string
}

func (Figure) Fill(templ string) (string, error) {
	return "", nil
}
