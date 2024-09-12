package sqs_test

func must[T any](val T, e error) T {
	if e != nil {
		panic(e)
	}

	return val
}
