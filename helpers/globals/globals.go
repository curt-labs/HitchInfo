package globals

import (
	"flag"
	"github.com/gorilla/sessions"
)

const (
	API_DOMAIN = "http://goapi.curtmfg.com"
	API_KEY    = "8aee0620-412e-47fc-900a-947820ea1c1d"

	SESSION_KEY          = "curt_reports"
	SESSION_ERROR_KEY    = "Errors"
	SESSION_CUSTOMER_KEY = "Customer"

	SESSION_USER_PUBLIC_KEY  = "UserPublicKey"
	SESSION_USER_PRIVATE_KEY = "UserPrivateKey"
)

var (
	GLOBALS = map[string]string{
		"LOGIN_HEADING":      "Welcome",
		"APPLICATION_DOMAIN": "http://goGoAdmin.curtmfg.com",
		"SITE_NAME":          "eCommerce Management Portal",
	}

	// Gorilla Session Store
	Store = sessions.NewCookieStore([]byte("doughboy"))
)

func SetGlobals() {

	flag.Parse()
}

func GetGlobal(k string) string {
	if str, ok := GLOBALS[k]; ok {
		return str
	}
	return ""
}
