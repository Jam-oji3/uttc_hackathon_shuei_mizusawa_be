package util

import (
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func ExtractNouns(text string) ([]string, error) {
	// ipadic辞書は内蔵済みなので、New() でOK
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	tokens := t.Tokenize(text)
	var nouns []string

	for _, tok := range tokens {
		features := tok.Features()
		if len(features) > 0 {
			// 品詞の1つ目だけを比較
			pos := strings.Split(features[0], "-")[0]
			if pos == "名詞" {
				nouns = append(nouns, tok.Surface)
			}
		}
	}

	//名刺がない場合明示的に空配列を返す
	if len(nouns) == 0 {
		return []string{}, nil
	}

	return nouns, nil
}
