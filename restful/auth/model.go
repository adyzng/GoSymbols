package auth

/**
*
*  Following struct are from graph api response json object.
*
*  refer: https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-v2-protocols-oauth-code
*
**/

// GraphUser from microsoft Graph API
//
type GraphUser struct {
	ID             string   `json:"id,omitempty"`
	Enabled        bool     `json:"accountEnabled,omitempty"`
	UserType       string   `json:"userType,omitempty"`
	DisplayName    string   `json:"displayName,omitempty"`
	GivenName      string   `json:"givenName,omitempty"`
	AboutMe        string   `json:"aboutMe,omitempty"`
	Mail           string   `json:"mail,omitempty"`
	JobTitle       string   `json:"jobTitle,omitempty"`
	MobilePhone    string   `json:"mobilePhone,omitempty"`
	CompanyName    string   `json:"companyName,omitempty"`
	Department     string   `json:"department,omitempty"`
	BusinessPhones []string `json:"businessPhones,omitempty"`
}

// GraphToken response data from ad server
//
type GraphToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	Type         string `json:"token_type,omitempty"`
	ExpireAt     int64  `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	State        string `json:"-"`
	UserName     string `json:"-"`
}

// GraphError wrap of Graph API error
//
type GraphError struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
	Codes       []int  `json:"error_codes,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	TraceID     string `json:"trace_id,omitempty"`
}
