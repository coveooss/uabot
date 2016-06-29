package explorerlib

type Config struct {
	FetchNumberOfResults           int        `json:"fetchQueryNumber"`
	DocumentsExplorationPercentage float64 `json:"explorationRatio"`
	FieldsToExploreEqually         []string `json:"fields"`
	SearchEndpoint                 string        `json:"searchEndpoint"`
	SearchToken	string        `json:"searchToken"`
	NumberOfQueryByLanguage        int `json:"numberOfQueryPerLanguage"`
	AnalyticsEndpoint              string `json:"analyticsEndpoint"`
	Org                            string `json:"org"`
	OutputFilePath                 string `json:"outputFilePath"`
	AverageNumberOfWordsPerQuery   int        `json:"avgNumberWordsPerQuery"`
}
