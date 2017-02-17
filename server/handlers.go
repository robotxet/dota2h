package server

import (
    "io"
    "io/ioutil"
    "log"
    "os"
    "net/http"
    "strconv"
)

func (s *Server) renderTemplate(wr io.Writer, key string, name string, data interface{}) {
    s.tMutex.RLock()
    t := s.templates[key]
    s.tMutex.RUnlock()
    err := t.ExecuteTemplate(wr, name, data)
    if err != nil {
        log.Printf("Error rendering template: %s", err.Error())
    }
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    s.renderTemplate(w, "layout", "error"+strconv.Itoa(status)+".html", nil)
}

func (s *Server) staticHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[1:]
    end := path[len(path)-1:]
    if "/" == end {
        s.errorHandler(w, r, http.StatusForbidden)
        return
    }
    if _, err := os.Stat(path); os.IsNotExist(err) {
        s.errorHandler(w, r, http.StatusNotFound)
        return
    }
    log.Printf(path)
    http.ServeFile(w, r, path)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        s.errorHandler(w, r, 404)
        return
    }
    s.renderTemplate(w, "layout", "index.html", nil)
}

func (s *Server) imageLoadHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/load_image" {
        s.errorHandler(w, r, 404)
        return
    }
    body, err := ioutil.ReadAll(r.Body);
    if err != nil {
        log.Println("error")
        return
    }
    log.Println(r)
    log.Println(r.Header.Get("Content-Type"))
    log.Println(body[0])
}