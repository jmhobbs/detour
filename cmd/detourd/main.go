package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/jmhobbs/detour/pkg/hosts"
)

type APIResponse struct {
	Hosts hosts.HostMapping `json:"hosts"`
}

var (
	staticBox packr.Box
	assetBox  packr.Box
)

func init() {
	staticBox = packr.NewBox("./static")
	assetBox = packr.NewBox("./assets")
}

func main() {
	http.HandleFunc("/api/list", list)
	http.Handle("/api/set", restrictMethod(http.HandlerFunc(set), "POST"))
	http.Handle("/api/unset", restrictMethod(http.HandlerFunc(unset), "POST"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticBox)))
	http.HandleFunc("/", index)
	log.Println("detour is up and running at http://127.0.0.1:9090/")
	http.ListenAndServe("127.0.0.1:9090", nil)
}

//////// Handlers

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, _ := template.New("index").Parse(assetBox.String("index.html"))
	t.Execute(w, map[string]interface{}{})
}

func list(w http.ResponseWriter, r *http.Request) {
	mapping := mustMap(w)
	if mapping == nil {
		return
	}

	mustRespond(w, *mapping)
}

func set(w http.ResponseWriter, r *http.Request) {
	file, mapping := mustMapRW(w)
	if file == nil {
		return
	}
	defer file.Close()

	host := r.PostFormValue("host")
	ip := r.PostFormValue("ip")

	if host == "" {
		http.Error(w, "'host' is required'", http.StatusBadRequest)
		return
	}
	if ip == "" {
		http.Error(w, "'ip' is required'", http.StatusBadRequest)
		return
	}

	mapping.Add(hosts.IPAddress(ip), hosts.Hostname(host))

	err := hosts.UpsertHostBock(mapping, file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error writing hosts file.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	mustRespond(w, mapping)
}

func unset(w http.ResponseWriter, r *http.Request) {
	file, mapping := mustMapRW(w)
	if file == nil {
		return
	}
	defer file.Close()

	host := r.PostFormValue("host")
	if host == "" {
		http.Error(w, "'host' is required'", http.StatusBadRequest)
		return
	}

	mapping.Remove(hosts.Hostname(host))

	err := hosts.UpsertHostBock(mapping, file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error writing hosts file.", http.StatusInternalServerError)
		return
	}

	mustRespond(w, mapping)
}

////////// Utilities

func restrictMethod(next http.Handler, httpMethod string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethod && r.Header.Get("X-HTTP-Method-Override") != httpMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mustMap(w http.ResponseWriter) *hosts.HostMapping {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error accessing hosts file.", http.StatusInternalServerError)
		return nil
	}
	defer file.Close()

	mapping, err := hosts.ExtractHostBlock(file)
	if err != nil {
		file.Close()
		log.Println(err)
		http.Error(w, "Error reading hosts file.", http.StatusInternalServerError)
		return nil
	}

	return &mapping
}

func mustMapRW(w http.ResponseWriter) (*os.File, hosts.HostMapping) {
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error accessing hosts file.", http.StatusInternalServerError)
		return nil, hosts.HostMapping{}
	}

	mapping, err := hosts.ExtractHostBlock(file)
	if err != nil {
		file.Close()
		log.Println(err)
		http.Error(w, "Error reading hosts file.", http.StatusInternalServerError)
		return nil, hosts.HostMapping{}
	}

	file.Seek(0, 0)

	return file, mapping
}

func mustRespond(w http.ResponseWriter, mapping hosts.HostMapping) {
	enc := json.NewEncoder(w)
	err := enc.Encode(APIResponse{mapping})
	if err != nil {
		log.Println(err)
		http.Error(w, "Error rendering hosts map.", http.StatusInternalServerError)
	}
}
