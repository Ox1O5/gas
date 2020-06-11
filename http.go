package gas

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePaht = "/_gascache/"

//HTTPPool implements PeerPicker for a pool of HTTP peers
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
}

//NewHTTPPool initializes an HTTP Pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePaht,
	}
}

//Log with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Sever %s]%s", p.self, fmt.Sprintf(format, v...))
}

//ServerHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	//<basePath>/<groupName>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
