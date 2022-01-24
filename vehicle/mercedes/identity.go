package mercedes

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"golang.org/x/oauth2"
)

type IdentityOptions func(c *Identity) error

// WithToken provides an oauth2.Token to the client for auth.
func WithToken(t *oauth2.Token) IdentityOptions {
	return func(v *Identity) error {
		v.ExpiringSource.Apply(t)
		return nil
	}
}

type Identity struct {
	log *util.Logger
	*ExpiringSource
	sessionSecret []byte
	authC         chan<- bool
	AuthConfig    *oauth2.Config
}

func generateSecret() ([]byte, error) {
	var b [16]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	return b[:], err
}

// TODO SessionSecret from config/persistence
func NewIdentity(log *util.Logger, id, secret string, options ...IdentityOptions) (*Identity, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://id.mercedes-benz.com")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC provider: %s", err)
	}

	v := &Identity{
		log:            log,
		ExpiringSource: new(ExpiringSource),
		AuthConfig: &oauth2.Config{
			ClientID:     id,
			ClientSecret: secret,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOfflineAccess, "mb:vehicle:mbdata:evstatus"},
		},
	}

	v.sessionSecret, err = generateSecret()

	for _, o := range options {
		if err == nil {
			err = o(v)
		}
	}

	return v, err
}

var _ api.ProviderLogin = (*Identity)(nil)

func (v *Identity) SetCallbackParams(uri string, authC chan<- bool) {
	v.AuthConfig.RedirectURL = uri
	v.authC = authC
}

func (v *Identity) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := NewState(v.sessionSecret)

		b, _ := json.Marshal(struct {
			LoginUri string `json:"loginUri"`
		}{
			LoginUri: v.AuthConfig.AuthCodeURL(state.Encrypt(), oauth2.AccessTypeOffline,
				oauth2.SetAuthURLParam("prompt", "login consent"),
			),
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (v *Identity) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.ExpiringSource.Apply(nil)
		v.authC <- false

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(nil)
	}
}

func (v *Identity) CallbackHandler(baseURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.log.TRACE.Println("callback request retrieved")

		data, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Fprintln(w, "invalid response:", data)
			return
		}

		if error, ok := data["error"]; ok {
			fmt.Fprintf(w, "error: %s: %s\n", error, data["error_description"])
			return
		}

		states, ok := data["state"]
		if !ok || len(states) != 1 {
			fmt.Fprintln(w, "invalid state response:", data)
			return
		} else if err := Validate(states[0], v.sessionSecret); err != nil {
			fmt.Fprintf(w, "failed state validation: %s", err)
			return
		}

		codes, ok := data["code"]
		if !ok || len(codes) != 1 {
			fmt.Fprintln(w, "invalid response:", data)
			return
		}

		token, err := v.AuthConfig.Exchange(context.Background(), codes[0])
		if err != nil {
			fmt.Fprintln(w, "token error:", err)
			return
		}

		if token.Valid() {
			v.log.TRACE.Println("sending login update...")
			v.ExpiringSource.Apply(token)
			v.authC <- true
		}

		http.Redirect(w, r, baseURI, http.StatusFound)
	}
}
