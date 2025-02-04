package apl

import "github.com/axiomhq/axiom-go/axiom/query"

// Result is the result of an APL query. It adds the APL query request alongside
// the query result it created, making it a superset of `query.Result`
type Result struct {
	*query.Result

	// Request is the APL query request that created the result.
	Request *query.Query `json:"request"`
	// The datasets that were queried in order to create the result.
	Datasets []string `json:"datasetNames"`
}
