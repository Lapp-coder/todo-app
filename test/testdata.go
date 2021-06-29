package test

const (
	DefaultTime = "19700101 "
)

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}
