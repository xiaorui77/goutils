package fileutils

import "testing"

func TestWindowsName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Windows name 01", "  hello  world  ", "hello world"},
		{"Windows name 02", "  hello \tworld  ", "hello world"},
		{"Windows name 03", "  hello >> ? < | world  ", "hello world"},
	}

	for _, test := range tests {
		actual := WindowsName(test.input)
		if actual != test.expected {
			t.Errorf("Test %s: expected %s, actual %s", test.name, test.expected, actual)
		}
	}
}

func TestHTMLDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Windows name 01", "hello&lt;&gt;&amp;&quot;&copy;world", "hello<>&\"©world"},
	}

	for _, test := range tests {
		actual := HTMLDecode(test.input)
		if actual != test.expected {
			t.Errorf("Test %s: expected %s, actual %s", test.name, test.expected, actual)
		}
	}
}

func TestXMLDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Windows name 01", "hello&lt;&gt;&amp;&quot;&apos;world", "hello<>&\"'world"},
	}

	for _, test := range tests {
		actual := XMLDecode(test.input)
		if actual != test.expected {
			t.Errorf("Test %s: expected %s, actual %s", test.name, test.expected, actual)
		}
	}
}
