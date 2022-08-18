package test

type Example struct{}

func (e *Example) SimpleRefReceiver() {
}

func (e Example) SimpleValueReceiver() {
}

func (*Example) SimpleUnnamedRefReceiver() {
}

func (Example) SimpleUnnamedValueReceiver() {
}

func (e *Example) ParamsInline(a string, b int) {
}

func (e *Example) ParamsMultiline(
	a string,
	b int,
	c interface{},
) {
}

// func (e *Example) ParamsMultiline(
// 	a string,
// 	b int,
// 	c interface{},
// ) (
// 	bool,
// 	error,
// ) {
// 	return false, nil
// }
