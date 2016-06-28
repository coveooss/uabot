package natural_language

import (
	"github.com/coveo/go-coveo/search"
	"math"
)

type Index struct {
	Client search.Client
}

func NewIndex(endpoint string, searchToken string) (Index, error){
	client,err := search.NewClient(search.Config{
		Endpoint:endpoint,
		Token:searchToken,
		UserAgent:"",
	})
	return Index{Client:client}, err
}

func (index *Index) FetchLanguages() ([]string, error){
	languageFacetValues, err := index.Client.ListFacetValues("@syslanguage", math.MaxInt16)
	languages := []string{}
	for _, value := range languageFacetValues.Values {
		languages = append(languages, value.Value)
	}
	return languages, err
}

func (index *Index) FetchFieldValues(field string) (*search.FacetValues, error) {
	return index.Client.ListFacetValues(field, 1000)
}

func (index *Index) FindTotalCountFromQuery(query search.Query) (int, error) {
	response, status := index.Client.Query(query)
	return response.TotalCount, status
}

func (index *Index) FetchResponse(queryExpression string, numberOfResults int) (*search.Response, error) {
	return index.Client.Query(search.Query{
		AQ:queryExpression,
		NumberOfResults:numberOfResults,
	})
}

