package main

import (
	"flag"
	"github.com/sourcegraph/buckler/shield"
	"go/build"
	"image/color"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
)

var bind = flag.String("bind", ":3000", "HTTP bind address")
var assets = flag.String("assets", filepath.Join(defaultBase("github.com/sqs/counter"), "assets"), "badge assets directory")

var badgeColor color.RGBA

func init() {
	shield.Init(*assets)

	var err error
	badgeColor, err = shield.GetColor("blue")
	if err != nil {
		log.Fatal("GetColor:", err)
	}
}

func main() {
	flag.Parse()
	http.Handle("/favicon.ico", http.HandlerFunc(http.NotFound))
	http.Handle("/", &counters{counts: make(map[string]int)})
	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type counters struct {
	counts map[string]int
	mu     sync.Mutex
}

func (c *counters) incr(name string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[name]++
	return c.counts[name]
}

func (c *counters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origPath := r.URL.Path
	path := filepath.Clean(origPath)
	if path != origPath {
		log.Printf("REDIRECT %s -> %s", origPath, path)
		http.Redirect(w, r, path, http.StatusMovedPermanently)
		return
	}

	count := c.incr(path)
	log.Printf("COUNT    %d %s", count, path)
	data := shield.Data{Vendor: "views", Status: strconv.Itoa(count), Color: badgeColor}
	w.Header().Add("content-type", "image/png")
	w.Header().Add("Cache-Control", "max-age=0, private, no-cache, must-revalidate")
	shield.PNG(w, data)
}

func defaultBase(path string) string {
	p, err := build.Default.Import(path, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}
