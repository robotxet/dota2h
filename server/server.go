package server

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "sync"
    "time"

    "github.com/fsnotify/fsnotify"
)

//Config is a struct that stores server configuration
type Config struct {
    HTTPPort      int            `json:"httpPort"`
}

//Server is a main server struct
type Server struct {
    config Config

    tMutex    sync.RWMutex
    templates map[string]*template.Template
}

//ParseConfig returns server Config from file on the given path
func ParseConfig(path string) Config {
    file, err := os.Open(path)
    if err != nil {
        log.Panicf("Can't load config: %s", err.Error())
    }
    decoder := json.NewDecoder(file)
    var config Config
    err = decoder.Decode(&config)
    if err != nil {
        log.Panicf("Can't parse config file: %s. Error: %s", path, err.Error())
    }
    return config
}

func projectPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Printf("Can't get project path: %s", err.Error())
        return ""
    }
    return dir
}

func templatePath(pattern string) string {
    return path.Join(projectPath(), "template", pattern)
}

func (s *Server) parseTemplates() {
    t := template.New("layout")

    _, err := t.ParseGlob(templatePath("*.html"))

    if err != nil {
        log.Fatalf("Error loading templates: %s", err.Error())
    }

    s.tMutex.Lock()
    s.templates["layout"] = t
    s.tMutex.Unlock()
}

func (s *Server) watchTemplates() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Printf("Can't watch for templates! Static won't be reloaded. Error: %s", err.Error())
        return
    }

    defer watcher.Close()
    err = watcher.Add(templatePath(""))
    if err != nil {
        log.Printf("Can't watch for templates! Static won't be reloaded. Error: %s", err.Error())
        return
    }
    for {
        <-watcher.Events

    wait:
        select {
        case <-watcher.Events:
            goto wait
        case <-time.After(time.Second):
        }

        log.Printf("Parse template activated...")
        s.parseTemplates()
    }
}

func (s *Server) Run() {
    s.parseTemplates()
    go s.watchTemplates()


    http.HandleFunc("/", s.indexHandler)
    // http.HandleFunc("/process_tf", s.tfHandler)
    // htpp.HandleFunc("/load_image", s.imageLoadHandler))

    http.HandleFunc("/static/", s.staticHandler)

    log.Printf("starting server: %d ...", s.config.HTTPPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.config.HTTPPort), nil))
    log.Printf("stopping server")
}

//New creates new Server
func New(config Config) *Server {
    // templates
    tMutex := sync.RWMutex{}
    templates := make(map[string]*template.Template)

    return &Server{config: config,
        tMutex:         tMutex,
        templates:      templates,
    }
}