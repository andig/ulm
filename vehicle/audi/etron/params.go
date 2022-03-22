package etron

import "net/url"

const AuthClientID = ""

// Authorization parameters
var AuthParams = url.Values(map[string][]string{
	"response_type": {"code id_token token"},
	"client_id":     {"f4d0934f-32bf-4ce4-b3c4-699a7049ad26@apps_vw-dilab_com"},
	"redirect_uri":  {"myaudi:///"},
	"scope":         {"openid profile mbb"}, // vin badge birthdate nickname email address phone name picture
	"prompt":        {"login"},
	"ui_locales":    {"de-DE"},
})
