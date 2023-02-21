package papapi

import "encoding/json"

type GridRequest struct {
	filters       []Filter
	columns       []string
	sortColumn    string
	sortAscending bool
	offset        int
	limit         int
	Request
}

type GridResponse struct {
	Response Response
}

type GridResponsePayload struct {
	Rows [][]string
}

// Parse rows and build records
func (gr *GridResponse) Records() []map[string]string {

	var payload []GridResponsePayload
	json.Unmarshal(gr.Response.Body, &payload)

	// NOTE: We only support single request/response entry
	headers, entries := payload[0].Rows[0], payload[0].Rows[1:]

	var records []map[string]string

	for _, entry := range entries {
		record := make(map[string]string)
		for i := 0; i < len(entry); i += 2 {
			record[headers[i]] = entry[i]
		}
		records = append(records, record)
	}

	return records
}

func (gr *GridRequest) AddFilter(name string, op FilterOperator, val string) {
	f := Filter{
		Name:     name,
		Operator: op,
		Value:    val,
	}
	gr.filters = append(gr.filters, f)
}

func (gr *GridRequest) AddColumn(name string) {
	gr.columns = append(gr.columns, name)
}

func (gr *GridRequest) AddColumns(names ...string) {
	gr.columns = append(gr.columns, names...)
}

func (gr *GridRequest) SetOffset(val int) {
	gr.offset = val
}

func (gr *GridRequest) SetLimit(val int) {
	gr.limit = val
}

func (gr *GridRequest) Serialize() map[string]interface{} {
	m := make(map[string]interface{})
	// Filters are passed as an Array e.g. (name, op, val)
	var serFilters [][]string
	for _, f := range gr.filters {
		serFilters = append(serFilters, f.Serialize())
	}
	m["filters"] = serFilters
	m["columns"] = gr.columns
	m["sort_col"] = gr.sortColumn
	m["sort_asc"] = gr.sortAscending
	m["offset"] = gr.offset
	m["limit"] = gr.limit
	return m
}

func (gr *GridRequest) Do() (GridResponse, error) {
	res, err := gr.Request.Do(gr.Serialize())
	if err != nil {
		return GridResponse{}, err
	}
	return GridResponse{Response: res}, nil
}

func (gr *GridRequest) SetSorting(column string, ascending bool) {
	gr.sortColumn = column
	gr.sortAscending = ascending
}

func NewGridRequest(className, methodName string, sess *Session) *GridRequest {
	return &GridRequest{
		Request: Request{
			ClassName:  className,
			MethodName: methodName,
			session:    sess,
		},
	}
}
