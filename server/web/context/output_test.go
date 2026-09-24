package context

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStringsToJSON(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"ascii":"hello"}`,
			expected: `{"ascii":"hello"}`,
		},
		{
			input:    `{"bmp":"hello 世界"}`,
			expected: `{"bmp":"hello \u4e16\u754c"}`,
		},
		{
			input:    `{"emoji":"😀"}`,
			expected: `{"emoji":"\ud83d\ude00"}`,
		},
		{
			input:    `{"mixed":"A😀B"}`,
			expected: `{"mixed":"A\ud83d\ude00B"}`,
		},
		{
			input:    `{"supplementary":"𠀀"}`,
			expected: `{"supplementary":"\ud840\udc00"}`,
		},
	}

	for _, tc := range testCases {
		res := stringsToJSON(tc.input)
		if res != tc.expected {
			t.Errorf("stringsToJSON(%q) = %q, want %q", tc.input, res, tc.expected)
		}
		var decoded map[string]string
		if err := json.Unmarshal([]byte(res), &decoded); err != nil {
			t.Fatalf("json.Unmarshal failed on %q: %v", res, err)
		}
		var original map[string]string
		if err := json.Unmarshal([]byte(tc.input), &original); err != nil {
			t.Fatalf("json.Unmarshal failed on input %q: %v", tc.input, err)
		}
		if !reflect.DeepEqual(decoded, original) {
			t.Errorf("decoded = %+v, want %+v", decoded, original)
		}
	}
}

func TestOutput_JSON_Encoding_SupplementaryUnicode(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()

	ctx := NewContext()
	ctx.Reset(&Response{ResponseWriter: rec}, req)

	data := map[string]string{
		"emoji":         "😀",
		"chinese":       "世界",
		"supplementary": "𠀀",
	}
	err := ctx.Output.JSON(data, false, true)
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed on %q: %v", rec.Body.String(), err)
	}

	if decoded["emoji"] != "😀" {
		t.Errorf("decoded emoji = %q, want %q (raw wire: %s)", decoded["emoji"], "😀", rec.Body.String())
	}
	if decoded["chinese"] != "世界" {
		t.Errorf("decoded chinese = %q, want %q", decoded["chinese"], "世界")
	}
	if decoded["supplementary"] != "𠀀" {
		t.Errorf("decoded supplementary = %q, want %q", decoded["supplementary"], "𠀀")
	}
}
