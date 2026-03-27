package echo

import "testing"

func TestServiceEcho(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "preserves content",
			input: "Please echo 35",
			want:  "Please echo 35",
		},
	}

	service := NewService()

	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := service.Echo(testCase.input); got != testCase.want {
				t.Fatalf("Echo(%q) = %q, want %q", testCase.input, got, testCase.want)
			}
		})
	}
}
