package use

func IsZero[T comparable](input T) bool {
	return input == Zero[T]()
}

func Zero[T any]() T {
	return *new(T)
}

func GetOrDefaultFunc[T comparable](input T, getDefault func() T) T {
	return ifOrElseFunc(!IsZero(input), input, getDefault)
}

func GetOrDefault[T comparable](input T, defaultVal T) T {
	return If(!IsZero(input), input, defaultVal)
}

func ifOrElseFunc[T any](cond bool, thenVal T, elseFunc func() T) T {
	if cond {
		return thenVal
	}
	return elseFunc()
}

func If[T any](cond bool, thenVal, elseVal T) T {
	return ifOrElseFunc(cond, thenVal, func() T {
		return elseVal
	})
}
