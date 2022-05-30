package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const baseEndpoint = "https://admin.instantproxies.com/"
const loginPhp = "login_do.php"
const loginEndpoint = baseEndpoint + loginPhp
const mainPhp = "main.php"
const mainEndpoint = baseEndpoint + mainPhp
const checkIPEndpoint = "https://checkip.instantproxies.com/"

func main() {
	fmt.Println("Start manin")
	client := NewClient("123456", "secret")
	myIP, myIPErr := client.GetLocalEnvPublicIP()
	fmt.Printf("MY IP IS %v\n", myIP)
	fmt.Printf("ERROR MY IP %v\n", myIPErr)
	proxies, err := client.GetProxies()
	ips, ipErr := client.GetAuthorizedIPs()
	fmt.Println()
	fmt.Printf("PROXIES ERROR: %v", err)
	fmt.Println()
	fmt.Println()
	fmt.Printf("PROXIES: %v", proxies)
	fmt.Println()
	fmt.Println()
	fmt.Printf("IPS: %v", ips)
	fmt.Println()
	fmt.Println()
	fmt.Printf("IP ERROR: %v", ipErr)
	fmt.Println()
	fmt.Println()
	fmt.Println("End manin")
}

// ------------------------------------
//   Proxy
// ------------------------------------

type Proxy struct {
	IP   net.IP
	Port uint16
}

func (proxy *Proxy) String() string {
	return fmt.Sprintf("%s:%d", proxy.IP.String(), proxy.Port)
}

func MakeProxy(str string) (*Proxy, error) {
	str = strings.Trim(str, " \n\r\t")
	parts := strings.Split(str, ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy string '%s'", str)
	}

	ip := net.ParseIP(parts[0])
	port, portErr := strconv.Atoi(parts[1])

	if port < 0 || port > 65535 {
		portErr = errors.New("port out of range")
	}

	if ip == nil || portErr != nil {
		return nil, fmt.Errorf("invalid proxy string '%s'", str)
	}

	return &Proxy{ip, uint16(port)}, nil
}

func MakeProxies(strings []string) ([]*Proxy, error) {
	proxies := make([]*Proxy, len(strings))
	for i, line := range strings {
		proxy, parseErr := MakeProxy(line)
		if parseErr != nil {
			return nil, parseErr
		}
		proxies[i] = proxy
	}
	return proxies, nil
}

// ------------------------------------
//   API client
// ------------------------------------

type SimpleHTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type Client struct {
	UserName         string
	Password         string
	httpClient       SimpleHTTPClient
	initSuccessfully bool
}

func NewClient(username string, password string) *Client {
	client := &Client{
		UserName: username,
		Password: password,
	}
	jar, _ := cookiejar.New(nil)
	client.httpClient = &http.Client{
		Jar:           jar,
		CheckRedirect: dontRedirect,
	}
	return client
}

func dontRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func (client *Client) initHTTPClient() error {
	if client.initSuccessfully {
		return nil
	}

	payload := url.Values{}
	payload.Add("username", client.UserName)
	payload.Add("password", client.Password)
	payload.Add("button", "Sign+In")

	res, networkErr := client.httpClient.PostForm(loginEndpoint, payload)
	if networkErr != nil {
		return networkErr
	}

	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	bodyStr := string(bodyBytes)

	if strings.Contains(bodyStr, "Invalid username or password") {
		return errors.New("invalid username or password")
	}

	if res.StatusCode != 302 {
		return fmt.Errorf("expected 302 HTTP status for success login, got %d", res.StatusCode)
	}

	redirectLocation := res.Header.Get("Location")
	if redirectLocation != mainPhp {
		return fmt.Errorf("expect redirect location to be '%s', got '%s'", mainPhp, redirectLocation)
	}

	client.initSuccessfully = true
	return nil
}

// ------------------------------------
//   API client get data
// ------------------------------------

func (client *Client) getMainPhpText() (string, error) {
	res, networkErr := client.httpClient.Get(mainEndpoint)
	if networkErr != nil {
		return "", networkErr
	}

	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}

	return string(bodyBytes), nil
}

func getTextAreaInnerText(html string, id string) ([]string, error) {
	// todo consider an html parser?

	regex, regexErr := regexp.Compile(fmt.Sprintf("(?s)<textarea id=\"%s\".*?>(.*?)</textarea>", id))
	if regexErr != nil {
		return nil, regexErr
	}

	match := regex.FindStringSubmatch(html)
	if len(match) != 2 {
		return nil, fmt.Errorf("cannot find element '%s'", id)
	}

	replacer := strings.NewReplacer("\t", "", "\r", "", "\n\n", "\n")
	textAreaText := strings.Trim(replacer.Replace(match[1]), "\n")
	lines := strings.Split(textAreaText, "\n")

	return lines, nil
}

func (client *Client) GetProxies() ([]*Proxy, error) {
	initError := client.initHTTPClient()
	if initError != nil {
		return nil, initError
	}

	html, reqErr := client.getMainPhpText()
	if reqErr != nil {
		return nil, reqErr
	}

	lines, err := getTextAreaInnerText(html, "proxies-textarea")
	if err != nil {
		return nil, err
	}

	return MakeProxies(lines)
}

func (client *Client) GetAuthorizedIPs() ([]net.IP, error) {
	initError := client.initHTTPClient()
	if initError != nil {
		return nil, initError
	}

	html, reqErr := client.getMainPhpText()
	if reqErr != nil {
		return nil, reqErr
	}

	lines, err := getTextAreaInnerText(html, "authips-textarea")
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, len(lines))
	for i, line := range lines {
		ips[i] = net.ParseIP(line)
		if ips[i] == nil {
			return nil, fmt.Errorf("can not parse IP %s", line)
		}
	}
	return ips, nil
}

// ------------------------------------
//   API client set data
// ------------------------------------

func (client *Client) AddAuthorizedIP(ip net.IP) error {
	return client.AddAuthorizedIPs([]net.IP{ip})
}

func (client *Client) AddAuthorizedIPs(ips []net.IP) error {
	initError := client.initHTTPClient()
	if initError != nil {
		return initError
	}

	authorizedIPs, getIPErr := client.GetAuthorizedIPs()
	if getIPErr != nil {
		return getIPErr
	}
	// TODO make ip unique?
	authorizedIPs = append(authorizedIPs, ips...)
	return client.SetAuthorizedIPs(authorizedIPs)
}

func (client *Client) RemoveAuthorizedIP(ip net.IP) error {
	initError := client.initHTTPClient()
	if initError != nil {
		return initError
	}

	authorizedIPs, getIPErr := client.GetAuthorizedIPs()
	if getIPErr != nil {
		return getIPErr
	}
	// TODO make ip unique?
	idx := sort.Search(len(authorizedIPs), func(i int) bool {
		return authorizedIPs[i].Equal(ip)
	})
	authorizedIPs = append(authorizedIPs[:idx], authorizedIPs[idx+1:]...)
	return client.SetAuthorizedIPs(authorizedIPs)
}

func (client *Client) SetAuthorizedIPs(ips []net.IP) error {
	initError := client.initHTTPClient()
	if initError != nil {
		return initError
	}

	authorizedIPStrs := make([]string, len(ips))
	for i, ip := range ips {
		authorizedIPStrs[i] = ip.String()
	}

	payload := url.Values{}
	payload.Add("cmd", "Submit Update")
	payload.Add("authips", strings.Join(authorizedIPStrs, "\n"))

	response, networkErr := client.httpClient.PostForm(mainEndpoint, payload)
	if networkErr != nil {
		return networkErr
	}

	// todo check status & content of response
	if response.StatusCode != 200 {
		return fmt.Errorf("expected status 200, got %v", response.StatusCode)
	}

	return nil
}

// ------------------------------------
//   API util
// ------------------------------------

func (client *Client) GetLocalEnvPublicIP() (net.IP, error) {
	response, networkErr := client.httpClient.Get(checkIPEndpoint)
	if networkErr != nil {
		return nil, networkErr
	}

	bodyRaw, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	ipString := strings.Trim(string(bodyRaw)[3:], " \n\r\t")
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("can not parse IP '%s'", ipString)
	}
	return ip, nil
}

/**
func (client *Client) TestProxies() []bool {

}
*/
