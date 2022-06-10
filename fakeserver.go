// A dead simple replica of InstantProxy server
//
// Login detail:
//      Username: username
//      Password: password
//
// Usage:
// 		server := NewFakeServer{3000, nil}
// 		server.StartServing() // returns immediately
// 		... do something
// 		server.Stop()
// Or:
// 		server := NewFakeServer{3000, nil}
// 		server.ServeForever() // only stops on Ctrl+C signal
//

package client

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FakeServer struct {
	AuthIPs    []net.IP
	Port       int
	httpServer *http.Server
}

func NewFakeServer(port int, authIPs []net.IP) *FakeServer {
	if port == 0 {
		port = 3000
	}
	if authIPs == nil {
		authIPs = []net.IP{}
	}
	fakeServer := &FakeServer{
		AuthIPs: authIPs,
		Port:    port,
	}
	fakeServer.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(fakeServer.handle),
	}
	return fakeServer
}

// StartServing starts the server in another go routine and return immediately
func (server *FakeServer) StartServing() {
	log.Printf("Listening on port %d...\n", server.Port)
	go server.serve()
}

// ServeForever starts the server in the caller routine and blocks until Stop is called
func (server *FakeServer) ServeForever() {
	log.Printf("Listening on port %d...\n", server.Port)
	server.serve()
}

// Stop the server
func (server *FakeServer) Stop() {
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
}

func (server *FakeServer) serve() {
	if err := server.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

func (server *FakeServer) handle(w http.ResponseWriter, r *http.Request) {
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
			renderMainPHP(w, server.AuthIPs)
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
				server.AuthIPs = newIPs
				renderMainPHP(w, server.AuthIPs)
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

// ------------------------
// Utils functions
// ------------------------

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

func renderMainPHP(w http.ResponseWriter, authIPs []net.IP) {
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
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		return false
	}
	return r.Form.Get("username") == "username" &&
		r.Form.Get("password") == "password" &&
		r.Form.Get("button") == "Sign In"
}

func validateIPs(r *http.Request) []net.IP {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		return nil
	}
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
