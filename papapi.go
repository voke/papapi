package papapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiVersion = "70544a5f74e11e13b7b61c4d98ae77e"
const authenticateClassName string = "Gpf_Api_AuthService"
const authenticateMethodName string = "authenticate"

type Role string

const (
	Merchant  Role = "M"
	Affiliate Role = "A"
)

type FilterOperator string

const (
	Like                 FilterOperator = "L"
	NotLike              FilterOperator = "NL"
	Equals               FilterOperator = "E"
	NotEquals            FilterOperator = "NE"
	DateEquals           FilterOperator = "D="
	DateGreater          FilterOperator = "D>"
	DateLower            FilterOperator = "D<"
	DateEqualsGreater    FilterOperator = "D>="
	DateEqualsLower      FilterOperator = "D<="
	DaterangeIs          FilterOperator = "DP"
	TimeEquals           FilterOperator = "T="
	TimeGreater          FilterOperator = "T>"
	TimeLower            FilterOperator = "T<"
	TimeEqualsGreater    FilterOperator = "T>="
	TimeEqualsLower      FilterOperator = "T<="
	RangeToday           FilterOperator = "T"
	RangeYesterday       FilterOperator = "Y"
	RangeLast7Days       FilterOperator = "L7D"
	RangeLast30Days      FilterOperator = "L30D"
	RangeLast90Days      FilterOperator = "L90D"
	RangeThisWeek        FilterOperator = "TW"
	RangeLastWeek        FilterOperator = "LW"
	RangeLast2weeks      FilterOperator = "L2W"
	RangeLastWorkingWeek FilterOperator = "LWW"
	RangeThisMonth       FilterOperator = "TM"
	RangeLastMonth       FilterOperator = "LM"
	RangeThisYear        FilterOperator = "TY"
	RangeLastYear        FilterOperator = "LY"
)

type Session struct {
	AuthToken    string
	APIVersion   string
	Role         Role
	ID           string
	URL          string
	LanguageCode string
}

type Filter struct {
	Name     string
	Operator FilterOperator
	Value    string
}

type Request struct {
	ClassName  string
	MethodName string
	session    *Session
}

type RequestPayload struct {
	ClassName  string                   `json:"C"`
	MethodName string                   `json:"M"`
	IsFromAPI  string                   `json:"isFromApi"`
	SessionID  string                   `json:"S,omitempty"`
	Requests   []map[string]interface{} `json:"requests"`
}

type Response struct {
	Body []byte
}

func (f *Filter) Serialize() []string {
	return []string{f.Name, string(f.Operator), f.Value}
}

type Record map[string]string

// Custom JSON method to avoid escaping greater-signs
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func (rp *RequestPayload) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(rp)
	return buffer.Bytes(), err
}

// Wrapper function to generate necessary payload
func (r *Request) Do(pairs map[string]interface{}) (Response, error) {

	payload := RequestPayload{
		ClassName:  "Gpf_Rpc_Server",
		MethodName: "run",
		IsFromAPI:  "Y",
	}

	if r.session.ID != "" {
		payload.SessionID = r.session.ID
	}

	// NOTE: Only supports one request entry
	pairs["C"] = r.ClassName
	pairs["M"] = r.MethodName
	payload.Requests = append(payload.Requests, pairs)

	jsonBody, _ := payload.JSON()

	return r.Execute(jsonBody)

}

func (r *Request) Execute(jsonBody []byte) (Response, error) {

	data := url.Values{
		"D": {string(jsonBody)},
	}

	resp, err := http.PostForm(r.session.URL, data)

	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Body: body,
	}, nil
}

func (s *Session) LoginWithAuthToken(usr, pw, twofactorToken string) error {
	return nil
}

func (s *Session) Login(usr, pw string) error {

	req := NewFormRequest(authenticateClassName, authenticateMethodName, s)

	req.SetField("username", usr)
	req.SetField("password", pw)
	req.SetField("roleType", string(s.Role))
	req.SetField("apiVersion", apiVersion)
	req.SetField("language", s.LanguageCode)

	res, err := req.Do()

	if err != nil {
		return err
	}

	s.ID = res.Parse()["S"] // Set ID (that will be used in subsequent requests)

	return nil
}

func NewSession(url string, role Role) *Session {
	return &Session{
		URL:          url,
		Role:         role,
		LanguageCode: "en-US",
	}
}
