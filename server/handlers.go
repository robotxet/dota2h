package server

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/satori/go.uuid"
)

type TfResponse struct {
    Hero    string
    Rating  float64
    ImgData []byte
    History string
}

type TfResponseData []TfResponse

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

func (s *Server) saveImage(content []byte, iType string) string {
    date := time.Now().Local()
    imgPath := s.config.ImagePath + "/" + date.Format("06-01-02")
    if _, err := os.Stat(imgPath); os.IsNotExist(err) {
        os.Mkdir(imgPath, 0777)
    }
    imgName := uuid.NewV4();
    file, err := os.Create(imgPath + "/" + imgName.String() + "." + iType)
    defer file.Close()

    if err != nil {
        log.Println("Failed to save image: " + err.Error())
        return ""
    }

    var decoded []byte
    strContent := string(content)
    var regex = regexp.MustCompile(`base64,(.*)`)
    imgstring := regex.FindStringSubmatch(strContent)

    decoded, err = base64.StdEncoding.DecodeString(imgstring[1])
    if err != nil {
        log.Println("Failed to save image: " + err.Error())
        return ""
    }

    _, err = file.Write(decoded)
    if err != nil {
        log.Println("Failed to save image: " + err.Error())
        return ""
    }
    return date.Format("06-01-02") + "/" + imgName.String() + "." + iType
}

func (s *Server) imageLoadHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/load_image" {
        s.errorHandler(w, r, 404)
        return
    }
    contentType := r.Header.Get("Content-Type")
    imgType := contentType[6:len(contentType) - 8]
    log.Println(imgType)
    if !ImageTypes[imgType] {
        log.Println("Empty Content-Type")
        return
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println("Failed to save image")
        return
    }
    filename := s.saveImage(body, imgType)
    if filename == "" {
        s.errorHandler(w, r, http.StatusNotFound)
        return
    } else {
        w.Write([]byte(filename))
    }
}

func top3(value string) map[string]float64 {
    result := make(map[string]float64)
    splitted := strings.Split(value, "\n")
    top3 := splitted[:3]
    for _, item := range top3 {
        itemKv := strings.Split(item, ":")
        result[itemKv[0]], _ = strconv.ParseFloat(itemKv[1], 64)
    }
    log.Println(result)
    return result
}

func (s *Server) tfHandler( w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/process_tf" {
        s.errorHandler(w, r, 404)
        return
    }
    body, err := ioutil.ReadAll(r.Body)
    if (err != nil) {
        log.Println("Error reading request body")
        return
    }
    filename := string(body[:])
    cmd := exec.Command(s.config.ScriptPath, s.config.DataPath, s.config.ImagePath + "/" + filename)
    getPredictions(s.config.ImagePath + "/" + filename, s.config.DataPath)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err == nil {
        topresults := top3(string(out.Bytes()))
        var avatar  []byte
        var history string
        var response = make(TfResponseData, len(topresults))
        index := 0
        for key, value := range topresults {
            if HeroMap[key] != "" {
                lorePath := s.config.LorePath + "/" + HeroMap[key] + "/"
                avatar, err = ioutil.ReadFile(lorePath + "avatar.png")
                if err != nil {
                    log.Println("Failed to get avatar")
                }
                data, err := ioutil.ReadFile(lorePath + "history_english.txt")
                if err != nil {
                    log.Println("Failed to get history")
                }
                history = strings.TrimSpace(string(data))
            }
            item := TfResponse {
                Hero: key,
                Rating: value,
                ImgData: avatar,
                History: history,
            }
            response[index] = item
            index += 1
        }
        js, err := json.Marshal(response)
        if err != nil {
            log.Println("failed to create json")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(js)
    } else {
        log.Println(stderr.String())
        log.Println("Failed to run tf script: " + err.Error())
    }
}
