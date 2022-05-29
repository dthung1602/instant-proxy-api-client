package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// ------------------------------------
//   Testing helpers
// ------------------------------------

func getFullResourcePath(fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "resources", "testing", fileName)
}

func getFileContent(fileName string) string {
	rawContent, _ := ioutil.ReadFile(getFullResourcePath(fileName))
	return string(rawContent)
}

func getMockHTTPResponse(fileName string, statusCode int, header http.Header) *http.Response {
	file, _ := os.Open(getFullResourcePath(fileName))

	response := &http.Response{
		StatusCode: statusCode,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     header,
		Body:       file,
		Close:      false,
	}

	return response
}

type MockHTTPClient struct {
	url  string
	data url.Values
	resp *http.Response
	err  error
}

func (mock *MockHTTPClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	mock.url = url
	mock.data = data
	return mock.resp, mock.err
}

func (mock *MockHTTPClient) Get(url string) (resp *http.Response, err error) {
	mock.url = url
	return mock.resp, mock.err
}

// ------------------------------------
//   Test make proxy
// ------------------------------------

func TestMakeProxy(t *testing.T) {
	proxyString := "1.0.123.255:12345"
	expected := &Proxy{
		IP:   net.ParseIP("1.0.123.255"),
		Port: 12345,
	}

	proxyAdrr, err := MakeProxy(proxyString)

	if err != nil {
		t.Errorf("found error %v", err)
	}

	if !reflect.DeepEqual(*proxyAdrr, *expected) {
		t.Errorf("proxy adrr is not equal")
	}
}

func TestMakeProxyWithPortError(t *testing.T) {
	proxyString := "1.0.123.255:-1"
	_, err := MakeProxy(proxyString)

	if err == nil {
		t.Errorf("error should be returned")
	}

	proxyString = "1.0.123.255:65536"
	_, err = MakeProxy(proxyString)

	if err == nil {
		t.Errorf("error should be returned")
	}
}

func TestMakeProxies(t *testing.T) {
	proxyStrings := []string{
		"1.0.123.255:12345",
		"212.123.84.94:1250",
	}
	expected := []*Proxy{
		{
			IP:   net.ParseIP("1.0.123.255"),
			Port: 12345,
		},
		{
			IP:   net.ParseIP("212.123.84.94"),
			Port: 1250,
		},
	}

	proxies, err := MakeProxies(proxyStrings)

	if err != nil {
		t.Errorf("error not expected: %v", err)
		return
	}

	if !reflect.DeepEqual(proxies, expected) {
		t.Error("proxies are not equal")
		return
	}
}

// ------------------------------------
//   Testing client
// ------------------------------------

func TestInitHTTPClient(t *testing.T) {
	mock := &MockHTTPClient{}
	headers := http.Header{
		"Location": []string{"main.php"},
	}
	mock.resp = getMockHTTPResponse("main.html", 302, headers)

	client := Client{
		UserId:           "123456",
		Password:         "secret",
		httpClient:       mock,
		initSuccessfully: false,
	}

	err := client.initHTTPClient()

	if err != nil {
		t.Errorf("init returns error %v", err)
		return
	}

	if !client.initSuccessfully {
		t.Errorf("init not is not successful")
		return
	}

	expectedFormData := url.Values{
		"username": []string{"123456"},
		"password": []string{"secret"},
		"button":   []string{"Sign+In"},
	}
	if !reflect.DeepEqual(expectedFormData, mock.data) {
		t.Error("http client doesn't received expected form data")
		return
	}

	expectedLoginEndpoint := "https://admin.instantproxies.com/login_do.php"
	if !(expectedLoginEndpoint == mock.url) {
		t.Error("login end point is incorrect")
		return
	}
}

// ------------------------------------
//   Testing extract data
// ------------------------------------

func TestGetProxiesByElementId(t *testing.T) {
	html := getFileContent("main.html")

	testCases := []struct {
		id             string
		expectedString []string
	}{
		{
			id: "proxies-textarea",
			expectedString: []string{
				"67.123.80.92:8800",
				"145.37.250.71:8800",
				"87.101.82.36:8800",
				"87.101.75.254:8800",
			},
		},
		{
			id: "authips-textarea",
			expectedString: []string{
				"48.201.98.5",
				"222.24.125.30",
				"119.19.82.29",
			},
		},
	}

	for _, testCase := range testCases {
		proxies, parseErr := getTextAreaInnerText(html, testCase.id)

		if parseErr != nil {
			t.Errorf("error getting proxies: %v", parseErr)
			return
		}

		if !reflect.DeepEqual(testCase.expectedString, proxies) {
			t.Errorf("got unexpteced proxies value")
			return
		}
	}
}

func TestClient_GetProxies(t *testing.T) {

}
