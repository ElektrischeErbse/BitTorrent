package main

import "testing"

func TestBuildBencode(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		// 整数测试
		{"Integer Positive", 42, "i42e"},
		{"Integer Negative", -42, "i-42e"},
		{"Integer Zero", 0, "i0e"},

		// 字符串测试
		{"String Normal", "spam", "4:spam"},
		{"String Empty", "", "0:"},

		// 列表测试
		{"List Mixed", []any{"spam", 42}, "l4:spami42ee"},
		{"List Empty", []any{}, "le"},
		{"List Nested", []any{[]any{1, 2}, 3}, "lli1ei2eei3ee"},

		// 字典测试
		{"Dict Sorted", map[string]any{"bar": "spam", "foo": 42}, "d3:bar4:spam3:fooi42ee"},
		{"Dict Empty", map[string]any{}, "de"},
		{"Dict Nested List", map[string]any{"spam": []any{"a", "b"}}, "d4:spaml1:a1:bee"},
		{"Dict Nested Dict", map[string]any{"info": map[string]any{"length": 1024}}, "d4:infod6:lengthi1024eee"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildBencode(tt.input)
			if err != nil {
				t.Errorf("%s test failed: %v", tt.name, err)
			}
			if result != tt.expected {
				t.Errorf("%s test failed: expected: %s, actually: %s", tt.name, tt.expected, result)
			}
		})
	}
}
