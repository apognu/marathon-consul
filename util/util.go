package util

func ReverseStringArray(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(ReverseStringArray(input[1:]), input[0])
}
