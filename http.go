package gas

import (
	"fmt"
	"github.com/Ox1O5/gas/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_gascache/"
	defaultReplicas = 50

)

//HTTPPool implements PeerPicker for a pool of HTTP peers
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
	mu sync.Mutex
	peers *consistenthash.Map
	httpGetter map[string]*httpGetter
}

//NewHTTPPool initializes an HTTP Pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
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
	// /<basePath>/<groupName>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
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

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get (group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
		)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading respones body: %v", err)
	}

	return bytes, nil
}

var _PeerGetter = (*httpGetter)(nil)

//Set updates the pool's peers
func (p *HTTPPool)Set(peers...string)  {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetter = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers{
		p.httpGetter[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}


//PickPeer picks a peer according to key
func (p *HTTPPool)PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key) ; peer != "" || peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetter[peer], true
	}
	return nil, false
}

var _PickPeer = (*PeerPicker)(nil)