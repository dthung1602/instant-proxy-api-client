package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

type ProxyAddr struct {
	IP   net.IP
	Port uint16
}

func (proxyAddr *ProxyAddr) String() string {
	return fmt.Sprintf("%s:%d", proxyAddr.IP.String(), proxyAddr.Port)
}

func MakeProxyAddr(str string) (*ProxyAddr, error) {
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

	return &ProxyAddr{ip, uint16(port)}, nil
}

type Client struct {
	UserId     string
	Password   string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	client := &Client{}

	jar, _ := cookiejar.New(nil)

	client.httpClient = &http.Client{
		Jar:           jar,
		CheckRedirect: dontRedirect,
	}

	payload := url.Values{}
	payload.Add("username", client.UserId)
	payload.Add("password", client.Password)
	payload.Add("button", "Sign+In")

	res, networkErr := client.httpClient.PostForm(loginEndpoint, payload)

	if networkErr != nil {
		return nil, networkErr
	}

	bodyBytes, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		return nil, readErr
	}

	bodyStr := string(bodyBytes)

	if strings.Contains(bodyStr, "Invalid username or password") {
		return nil, errors.New("invalid username or password")
	}

	if res.StatusCode != 302 {
		return nil, fmt.Errorf("expected 302 HTTP status for success login, got %d", res.StatusCode)
	}

	redirectLocation := res.Header.Get("Location")
	if redirectLocation != mainPhp {
		return nil, fmt.Errorf("expect redirect location to be '%s', got '%s'", mainPhp, redirectLocation)
	}

	return client, nil
}

func dontRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (client Client) getMainPhpText() (string, error) {
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

func getProxiesByElementId(html string, id string) ([]ProxyAddr, error) {
	// todo consider an html parser?

	lines := strings.Split(html, "\n")
	matchText := "id=\"" + id + "\""
	for i, line := range lines {
		if !strings.Contains(line, matchText) {
			continue
		}

		firstLineComponents := strings.Split(line, ">")
		if len(firstLineComponents) != 2 {
			return nil, errors.New("cannot extract proxies from HTML")
		}

		proxyStrings := []string{firstLineComponents[1]}

		// TODO optimize this
		for _, ln := range lines[i+1:] {
			if strings.Contains(ln, "<") {
				lastLineComponents := strings.Split(ln, "<")
				if len(lastLineComponents) != 2 {
					return nil, errors.New("cannot extract proxies from HTML")
				}
				proxyStrings = append(proxyStrings, lastLineComponents[0])
			} else {
				proxyStrings = append(proxyStrings, ln)
			}
		}

		proxies := make([]ProxyAddr, len(proxyStrings))

		for i, str := range proxyStrings {
			proxy, proxyErr := MakeProxyAddr(str)
			if proxyErr != nil {
				return nil, proxyErr
			}
			proxies[i] = *proxy
		}

		return proxies, nil
	}

	return nil, fmt.Errorf("cannot find element with id '%s'", id)
}

func (client Client) GetProxies() ([]ProxyAddr, error) {
	html, reqErr := client.getMainPhpText()

	if reqErr != nil {
		return nil, reqErr
	}

	proxies, parseErr := getProxiesByElementId(html, "proxies-textarea")

	if parseErr != nil {
		return nil, parseErr
	}

	return proxies, nil
}

/**
func (client Client) TestProxies() []bool {

}

func (client Client) GetAuthorizedIPs() []net.IP {

}

func (client Client) AddAuthorizedIP(ip net.IP) {

}

func (client Client) RemoveAuthorizedIP(ip net.IP) {

}

func (client Client) SetAuthorizedIPs(ips []net.IP) {

}

func (client Client) GetLocalEnvPublicIP() net.IP {

}
*/
