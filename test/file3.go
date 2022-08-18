package test

func (t *Example) ReturnsInline() (int, string) {
	return 0, ""
}

func (t *Example) ReturnsInlineNamed() (a, b int, c string) {
	return 0, 0, ""
}

func (t *Example) ReturnsMultilineNamed() (
	a, b int,
	c string,
	d error,
) {
	return 0, 0, "", nil
}

func (t *Example) ReturnsMultilineNamedWithComment() (
	a, b int,
	//	c string,
	d error,
) {
	return 0, 0, nil
}

/*

func (t *Example) InCommentBlock() {
}

*/
