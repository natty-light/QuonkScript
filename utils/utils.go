package utils

func Pop[T any](arr []T) []T {
	return arr[1:]
}
