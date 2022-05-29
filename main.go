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
	"strconv"
	"strings"
)

const baseEndpoint = "https://admin.instantproxies.com/"
const loginPhp = "login_do.php"
const loginEndpoint = baseEndpoint + loginPhp
const mainPhp = "main.php"

func main() {
	fmt.Println("hello world")
}

// ------------------------------------
//   Proxy
// ------------------------------------

type Proxy struct {
	IP   net.IP
	Port uint16
}

func (proxyAddr *Proxy) String() string {
	return fmt.Sprintf("%s:%d", proxyAddr.IP.String(), proxyAddr.Port)
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

type simpleHTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type Client struct {
	UserId           string
	Password         string
	httpClient       simpleHTTPClient
	initSuccessfully bool
}

func NewClient() *Client {
	client := &Client{}

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
	payload.Add("username", client.UserId)
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
	res, networkErr := client.httpClient.Get(mainPhp)

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

/**
func (client *Client) TestProxies() []bool {

}

func (client *Client) GetAuthorizedIPs() []net.IP {

}

func (client *Client) AddAuthorizedIP(ip net.IP) {

}

func (client *Client) RemoveAuthorizedIP(ip net.IP) {

}

func (client *Client) SetAuthorizedIPs(ips []net.IP) {

}

func (client *Client) GetLocalEnvPublicIP() net.IP {

}
*/
