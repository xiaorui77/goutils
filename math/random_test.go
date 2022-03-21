package math

import "testing"

func TestRandom16Str(t *testing.T) {
	str := Random16Str(10)
	t.Log(str)
}

func TestRandom62Str(t *testing.T) {
	str := Random62Str(10)
	t.Log(str)
}

func TestRandomStr(t *testing.T) {
	str := RandomStr(10, 34)
	t.Log(str)
}
