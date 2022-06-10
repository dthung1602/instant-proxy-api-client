package client

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

const adminEndpoint = "https://admin.instantproxies.com/"
const loginPhp = "login_do.php"
const mainPhp = "main.php"
const checkIPEndpoint = "https://checkip.instantproxies.com/"
const checkProxiesEndpoint = "https://instantproxies.com/proxytester/test.json.php"

// ------------------------------------
//   API client
// ------------------------------------

type Client struct {
	UserName         string
	Password         string
	Endpoint         string
	httpClient       SimpleHTTPClient
	initSuccessfully bool
}

func NewClient(username string, password string, endpoint string) *Client {
	if endpoint == "" {
		endpoint = adminEndpoint
	}
	client := &Client{
		UserName: username,
		Password: password,
		Endpoint: endpoint,
	}
	jar, _ := cookiejar.New(nil)
	client.httpClient = &http.Client{
		Jar: jar,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse // no follow redirect
		},
	}
	return client
}

func (client *Client) Authenticate() error {
	if client.initSuccessfully {
		return nil
	}

	payload := url.Values{}
	payload.Add("username", client.UserName)
	payload.Add("password", client.Password)
	payload.Add("button", "Sign In")

	res, networkErr := client.httpClient.PostForm(client.loginEndpoint(), payload)
	defer res.Body.Close()
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

func (client *Client) loginEndpoint() string {
	return client.Endpoint + loginPhp
}

func (client *Client) mainEndpoint() string {
	return client.Endpoint + mainPhp
}

// ------------------------------------
//   API client get data
// ------------------------------------

func (client *Client) getMainPhpText() (string, error) {
	res, networkErr := client.httpClient.Get(client.mainEndpoint())
	defer res.Body.Close()
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
	initError := client.Authenticate()
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
	initError := client.Authenticate()
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
	initError := client.Authenticate()
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
	return client.RemoveAuthorizedIPs([]net.IP{ip})
}

func (client *Client) RemoveAuthorizedIPs(ips []net.IP) error {
	initError := client.Authenticate()
	if initError != nil {
		return initError
	}

	authorizedIPs, getIPErr := client.GetAuthorizedIPs()
	if getIPErr != nil {
		return getIPErr
	}

	for _, ip := range ips {
		idx := -1
		for i, aip := range authorizedIPs {
			if aip.Equal(ip) {
				idx = i
				break
			}
		}
		if idx != -1 {
			authorizedIPs = append(authorizedIPs[:idx], authorizedIPs[idx+1:]...)
		}
	}
	return client.SetAuthorizedIPs(authorizedIPs)
}

func (client *Client) SetAuthorizedIPs(ips []net.IP) error {
	initError := client.Authenticate()
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

	response, networkErr := client.httpClient.PostForm(client.mainEndpoint(), payload)
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

func (client *Client) GetOwnedPublicIP() (net.IP, error) {
	response, networkErr := client.httpClient.Get(checkIPEndpoint)
	defer response.Body.Close()
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

func (client *Client) TestOwnedProxies() ([]bool, error) {
	proxies, err := client.GetProxies()
	if err != nil {
		return nil, err
	}
	return client.TestProxies(proxies), nil
}

func (client *Client) TestProxies(proxies []*Proxy) []bool {
	pool := NewWorkerPool[bool](4)
	for _, proxy := range proxies {
		pool.Submit(func() (bool, error) {
			return client.testProxy(proxy)
		})
	}
	pool.Close()
	pool.Wait()
	return pool.ResultValues()
}

func (client *Client) testProxy(proxy *Proxy) (bool, error) {
	fmt.Printf("Testing %v\n", proxy)
	query := url.Values{}
	query.Add("proxy", proxy.String())
	checkUrl := checkProxiesEndpoint + "?" + query.Encode()

	res, networkErr := client.httpClient.Get(checkUrl)
	defer res.Body.Close()
	if networkErr != nil {
		return false, networkErr
	}
	if res.StatusCode != 200 {
		return false, fmt.Errorf("expected HTTP status 200, got %d", res.StatusCode)
	}

	bodyRaw, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return false, readErr
	}
	body := strings.Trim(string(bodyRaw), " \n\r\t")
	passed := strings.HasSuffix(body, "200,PASSED")
	fmt.Printf("Result testing %v: %v\n", proxy, passed)
	return passed, nil
}
