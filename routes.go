package main

import (
	"fmt"
	"github.com/gokyle/webshell"
	"log"
	"net/http"
	"regexp"
)

var (
	page_title    string
	server_host   string
	server_secure bool
	server_dev    = true
	strip_views   = regexp.MustCompile("^/views/(.+)$")
)

type Page struct {
	Title     string
	Count     string
	Host      string
	ShortCode string
	Posted    bool
	ShowErr   bool
	Scheme    string
	Err       string
	Views     string
}

func getPageCount() (page_count string) {
	count, err := countShortened()
	if err != nil {
		page_count = "are no pages"
	} else {
		if count == 0 {
			page_count = "are no pages"
		} else {
			var verb string
			if count == 1 {
				verb = "is"
				page_count = "page"
			} else {
				verb = "are"
				page_count = "pages"
			}
			page_count = fmt.Sprintf("%s %d %s", verb, count,
				page_count)
		}
	}
	return
}

func NewPage() (page *Page) {
	page = new(Page)
	page.Title = page_title
	page.Host = server_host

	if server_secure {
		page.Scheme = "https"
	} else {
		page.Scheme = "http"
	}
	return
}

func home(w http.ResponseWriter, r *http.Request) {
	page := NewPage()
	servePage(page, w, r)
	return
}

func servePage(page *Page, w http.ResponseWriter, r *http.Request) {
	page.Count = getPageCount()
	out, err := webshell.ServeTemplate("templates/index.html", page)
	if err != nil {
		webshell.Error404(err.Error(), "text/plain", w, r)
	} else {
		w.Write(out)
	}
	LogRequest(page, r)
}

func serveViews(page *Page, w http.ResponseWriter, r *http.Request) {
	page.Count = getPageCount()
	out, err := webshell.ServeTemplate("templates/views.html", page)
	if err != nil {
		webshell.Error404(err.Error(), "text/plain", w, r)
	} else {
		w.Write(out)
	}
	LogRequest(page, r)
}

func serveErr(page *Page, err error, w http.ResponseWriter, r *http.Request) {
	page.ShowErr = true
	page.Err = err.Error()
	servePage(page, w, r)
}

func topRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newShortened(w, r)
		return
	}
	sid := r.URL.Path
	if len(sid) > 1 {
		sid = sid[1:len(sid)]
	}
	url, err := lookupShortCode(sid)
	if err != nil {
		home(w, r)
		return
	} else if url != "" {
		err = updateSidViews(sid)
		if err != nil {
			log.Printf("[+] error updating views: %s\n",
				err.Error())
		}
		http.Redirect(w, r, url, 301)
		return
	}

	home(w, r)
}

func newShortened(w http.ResponseWriter, r *http.Request) {
	page := NewPage()
	err := r.ParseForm()
	if err != nil {
		serveErr(page, err, w, r)
		return
	}

	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	if !authenticate(user, pass) {
		err = fmt.Errorf("Authenticated failed!")
		serveErr(page, err, w, r)
		return
	}
	sid := r.Form.Get("sid")
	url := r.Form.Get("url")
	if len(url) == 0 {
		err := fmt.Errorf("Invalid URL")
		serveErr(page, err, w, r)
		return
	} else if len(sid) > 0 {
		if surl, err := lookupShortCode(sid); err != nil {
			serveErr(page, err, w, r)
			return
		} else if surl != "" && surl != url {
			err = fmt.Errorf("URL already present.")
			sid, db_err := urlToSid(url)
			if db_err == nil {
				page.ShortCode = sid
				page.Posted = true
			} else {
				err = db_err
			}
			serveErr(page, err, w, r)
			return
		} else if err = insertShortened(sid, url); err != nil {
			serveErr(page, err, w, r)
			return
		} else {
			page.ShortCode = sid
			page.Posted = true
		}
	} else {
		sid, err := ShortenUrl(ValidateShortenedUrl)
		if err != nil {
			serveErr(page, err, w, r)
			return
		} else {
			page.Posted = true
			page.ShortCode = sid
		}
	}
	servePage(page, w, r)
	return
}

func getViews(w http.ResponseWriter, r *http.Request) {
	sid := strip_views.ReplaceAllString(r.URL.Path, "$1")
	page := NewPage()
	count, err := getSidViews(sid)

	var views string
	if err != nil {
		log.Println("[!] getViews error: ", err.Error())
		views = "no views"
	} else {
		if count == 0 {
			views = "no views"
		} else if count == 1 {
			views = "one view"
		} else {
			views = fmt.Sprintf("%d views", count)
		}
	}
	log.Printf("[-] %s -> %d\n", sid, count)
	page.ShortCode = sid
	page.Views = views
	serveViews(page, w, r)
}
