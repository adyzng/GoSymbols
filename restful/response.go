package restful

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/adyzng/GoSymbols/symbol"

	log "gopkg.in/clog.v1"
)

var (
	buffSize4k = 4096
	/* buffPool4K = sync.Pool{
		New: func() interface{} {
			// 4k buffer pool by default
			return bytes.NewBuffer(make([]byte, buffSize4k))
		},
	} */
)

var (
	ErrSucceed      = ErrCodeMsg{0, "ok"}
	ErrInvalidParam = ErrCodeMsg{100, "invalid parameter"}
	ErrServerInner  = ErrCodeMsg{101, "server inner error"}

	ErrInvalidBranch = ErrCodeMsg{200, "branch unavailable"}
	ErrExistOnLocal  = ErrCodeMsg{201, "branch exist in symbol store"}
	ErrUnknownBranch = ErrCodeMsg{202, "unknown branch"}
	ErrUnauthorized  = ErrCodeMsg{203, "unauthorized operation"}
)

// ErrCodeMsg is predefined error code and error message
//
type ErrCodeMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BranchList return branch list of current symbol store
//
type BranchList struct {
	Total   int              `json:"total"`
	Branchs []*symbol.Branch `json:"branchs"`
}
type BuildList struct {
	Branch string          `json:"branchName"`
	Total  int             `json:"total"`
	Builds []*symbol.Build `json:"builds"`
}
type SymbolList struct {
	Branch  string           `json:"branchName"`
	Build   string           `json:"buildID"`
	Total   int              `json:"total"`
	Symbols []*symbol.Symbol `json:"symbols"`
}

// RestResponse is the basic struct used to wrap data back to client in json format.
//
type RestResponse struct {
	ErrCodeMsg
	Data interface{} `json:"data"`
}

// ToJSON encoding to json string
//
func (r *RestResponse) ToJSON() string {
	buff := bytes.NewBuffer(make([]byte, buffSize4k))
	err := json.NewEncoder(buff).Encode(r)
	if err != nil {
		log.Error(2, "[Restful] json encoding failed with %v.", err)
		return ""
	}
	return buff.String()
}

// WriteJSON write json reponse to client
//
func (r *RestResponse) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Error(2, "[Restful] JSON.Encode(%+v) failed with %v.", r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
