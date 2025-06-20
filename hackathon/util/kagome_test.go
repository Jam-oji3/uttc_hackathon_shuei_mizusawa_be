package util_test

import (
	"reflect"
	"testing"

	"hackathon/util" // 実際のutilパッケージのパスに書き換えてください
)

func TestExtractNouns(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNouns []string
		wantErr   bool
	}{
		{
			name:      "simple sentence with nouns",
			input:     "私はリンゴを食べる。",
			wantNouns: []string{"私", "リンゴ"},
			wantErr:   false,
		},
		{
			name:      "sentence without nouns",
			input:     "走る速く。",
			wantNouns: []string{},
			wantErr:   false,
		},
		{
			name:      "empty string",
			input:     "",
			wantNouns: []string{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNouns, err := util.ExtractNouns(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ExtractNouns() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotNouns, tt.wantNouns) {
				t.Errorf("ExtractNouns() = %v, want %v", gotNouns, tt.wantNouns)
			}
		})
	}
}
