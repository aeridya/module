package cookie

import (
	"fmt"

	"github.com/aeridya/core"
	"github.com/aeridya/module"

	"github.com/gorilla/securecookie"

	"net/http"
	"time"
)

type cookie struct {
	module.Module
	handler *securecookie.SecureCookie

	cookieHash  []byte
	cookieBlock []byte
}

var (
	cook *cookie
)

func Setup(options ...module.Option) {
	cook = &cookie{}
	cook.ParseOpts(options)
	cook.handler = securecookie.New(cook.cookieHash, cook.cookieBlock)
}

func SetHash(h string) module.Option {
	return func() {
		cook.cookieHash = []byte(h)
	}
}

func SetBlock(b string) module.Option {
	return func() {
		cook.cookieBlock = []byte(b)
	}
}

func Generate() module.Option {
	return func() {
		cook.cookieHash = securecookie.GenerateRandomKey(32)
		cook.cookieBlock = securecookie.GenerateRandomKey(32)
	}
}

func Add(resp *core.Response, name string, hour int, values map[string]string) error {
	enc, err := cook.handler.Encode(name, values)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:     name,
		Value:    enc,
		Path:     "/",
		MaxAge:   60 * hour,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	}
	http.SetCookie(resp.W, &c)
	return nil
}

func Delete(resp *core.Response, name string) {
	c := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(resp.W, &c)
}

func Get(resp *core.Response, name string) (map[string]string, error) {
	ck, err := resp.R.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cook.handler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func GetValues(resp *core.Response, name string, vars ...string) ([]string, error) {
	ck, err := resp.R.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cook.handler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(vars))
	for i := range vars {
		if v, ok := values[vars[i]]; !ok {
			return nil, fmt.Errorf("No value for key: %s", vars[i])
		} else {
			out[i] = v
		}
	}
	return out, nil
}

func GetRaw(w http.ResponseWriter, r *http.Request, name string) (map[string]string, error) {
	ck, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cook.handler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func GetValuesRaw(w http.ResponseWriter, r *http.Request, name string, vars ...string) ([]string, error) {
	ck, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cook.handler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(vars))
	for i := range vars {
		if v, ok := values[vars[i]]; !ok {
			return nil, fmt.Errorf("No value for key: %s", vars[i])
		} else {
			out[i] = v
		}
	}
	return out, nil
}

func AddRaw(w http.ResponseWriter, r *http.Request, name string, hour int, values map[string]string) error {
	enc, err := cook.handler.Encode(name, values)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:     name,
		Value:    enc,
		Path:     "/",
		MaxAge:   60 * hour,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	}
	http.SetCookie(w, &c)
	return nil
}

func DeleteRaw(w http.ResponseWriter, r *http.Request, name string) {
	c := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
