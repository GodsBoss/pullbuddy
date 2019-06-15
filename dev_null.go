package pullbuddy

type devNull string

func (w devNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}
