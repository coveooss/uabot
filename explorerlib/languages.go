package explorerlib

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var (
	supportedLanguage = []language.Tag{
		language.English,
		language.French,
		language.German,
		language.Spanish,
		language.Dutch,
		language.Portuguese,
		language.Armenian,
		language.Catalan,
		language.Danish,
		language.Finnish,
		language.Italian,
		language.Hungarian,
		language.Norwegian,
		language.Romanian,
		language.Russian,
		language.Swedish,
		language.Turkish,
		language.Afrikaans,
		language.Arabic,
		language.Bulgarian,
		language.Chinese,
		language.Croatian,
		language.Czech,
		language.Estonian,
		language.Filipino,
		language.Greek,
		language.Hebrew,
		language.Hindi,
		language.Icelandic,
		language.Indonesian,
		language.Japanese,
		language.Korean,
		language.Latvian,
		language.Lithuanian,
		language.Persian,
		language.Polish,
		language.Serbian,
		language.Slovak,
		language.Slovenian,
		language.Swahili,
		language.Thai,
		language.Ukrainian,
		language.Vietnamese,
	}
	englishDisplay = display.English.Tags()

	exceptionLanguages = map[string]string{
		"Esperanto": "eo",
		"Norwegian": "no",
	}
)

func LanguageToTag(language string) string {
	for lang, code := range exceptionLanguages {
		if lang == language {
			return code
		}
	}

	for _, lang := range supportedLanguage {
		if language == englishDisplay.Name(lang) {
			return lang.String()
		}
	}
	return ""
}
