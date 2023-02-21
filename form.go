package papapi

import "encoding/json"

type FormRequest struct {
	fields map[string]string
	Request
}

type FormResponse struct {
	Fields [][]string `json:"fields"`
}

func (fr *FormResponse) Parse() map[string]string {

	vals := make(map[string]string)

	for i, arr := range fr.Fields {
		if i == 0 {
			continue
		}

		vals[arr[0]] = arr[1]
	}

	return vals

}

func (fr *FormRequest) SetField(key, value string) {
	fr.fields[key] = value
}

func (fr *FormRequest) Serialize() map[string]interface{} {

	var dataFields [][]string
	dataFields = append(dataFields, []string{"name", "value"})

	for k, v := range fr.fields {
		dataFields = append(dataFields, []string{k, v})
	}

	data := make(map[string]interface{})
	data["fields"] = dataFields

	return data

}

func (fr *FormRequest) Do() (FormResponse, error) {

	var res []FormResponse

	resp, err := fr.Request.Do(fr.Serialize())

	if err != nil {
		return FormResponse{}, err
	}

	json.Unmarshal(resp.Body, &res)

	return res[0], nil
}

func NewFormRequest(className, methodName string, sess *Session) *FormRequest {
	return &FormRequest{
		fields: make(map[string]string),
		Request: Request{
			ClassName:  className,
			MethodName: methodName,
			session:    sess,
		},
	}
}
