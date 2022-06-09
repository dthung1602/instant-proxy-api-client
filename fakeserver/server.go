package main

import (
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var authIPs []net.IP

func main() {
	initState()
	serve()
}

func initState() {
	authIPs = []net.IP{
		net.ParseIP("48.201.98.5"),
		net.ParseIP("222.24.125.30"),
		net.ParseIP("119.19.82.29"),
	}
}

func serve() {
	log.Println("Listening on port 3000...")
	http.Handle("/", http.HandlerFunc(serveFile))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DONE!")
}

func getFullPath(file string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "resources", "testing", file)
}

func sendHTMLFile(w http.ResponseWriter, file string) {
	fileName := getFullPath(file)
	reader, _ := os.Open(fileName)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, reader)
}

func renderMainPHP(w http.ResponseWriter) {
	temp, _ := template.ParseFiles(getFullPath("main.template"))
	authips := ""
	for _, ip := range authIPs {
		if ip != nil {
			authips += ip.String() + "\n"
		}
	}
	log.Println("IPS: ", authips)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := temp.Execute(w, struct{ Authips string }{Authips: authips})
	log.Println(err)
}

func validatePassword(r *http.Request) bool {
	r.ParseForm()
	log.Printf("%v", r.Form)
	return r.Form.Get("username") == "username" &&
		r.Form.Get("password") == "password" &&
		r.Form.Get("button") == "Sign In"
}

func validateIPs(r *http.Request) []net.IP {
	r.ParseForm()
	ipstrs := strings.Split(r.Form.Get("authips"), "\n")
	var newIPs []net.IP
	for _, ipstr := range ipstrs {
		ipstr = strings.Trim(ipstr, " \n\t\r")
		ip := net.ParseIP(ipstr)
		newIPs = append(newIPs, ip)
		if ip == nil && ipstr != "" {
			return nil
		}
	}
	if r.Form.Get("cmd") != "Submit Update" {
		return nil
	}
	return newIPs
}

func isAuthenticated(r *http.Request) bool {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "auth" && cookie.Value == "secret" {
			return true
		}
	}
	return false
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n---> Processing request <---")

	if r.Method == "GET" && (r.URL.Path == "/" || r.URL.Path == "" || r.URL.Path == "/index.php") {
		log.Printf("GET %s\n", r.URL.Path)
		w.Header().Set("Location", "login.php")
		w.WriteHeader(302)
		return
	}

	if r.Method == "GET" && r.URL.Path == "/login.php" {
		println("GET /login.php")
		sendHTMLFile(w, "login.php")
		return
	}

	if r.Method == "POST" && r.URL.Path == "/login_do.php" {
		log.Println("POST /login_do.php")
		if !validatePassword(r) {
			log.Println("Fail login")
			sendHTMLFile(w, "login_do.fail.php")
		} else {
			log.Println("Success login")
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: "secret",
			})
			w.Header().Set("Location", "main.php")
			w.WriteHeader(302)
		}
		return
	}

	if r.Method == "GET" && r.URL.Path == "/main.php" {
		log.Println("GET /main.php")
		if isAuthenticated(r) {
			log.Println("Auth")
			renderMainPHP(w)
		} else {
			log.Println("Not auth")
			w.Header().Set("Location", "index.php")
			w.WriteHeader(302)
		}
		return
	}

	if r.Method == "POST" && r.URL.Path == "/main.php" {
		log.Println("POST /main.php")
		if isAuthenticated(r) {
			log.Println("Auth")
			newIPs := validateIPs(r)
			if newIPs != nil {
				log.Println("Success change auth ip")
				authIPs = newIPs
				renderMainPHP(w)
			} else {
				log.Println("Fail change auth ip")
				w.WriteHeader(500)
			}
		} else {
			log.Println("Not auth")
			w.Header().Set("Location", "index.php")
			w.WriteHeader(302)
		}
		return
	}

	log.Printf("%s %s\n", r.Method, r.URL.Path)
	log.Println("404 Not found")
	w.WriteHeader(404)
}
