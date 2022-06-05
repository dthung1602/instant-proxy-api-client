package client

import (
	"io/ioutil"
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
//   Testing client
// ------------------------------------

func TestAuthenticate(t *testing.T) {
	mock := &MockHTTPClient{}
	headers := http.Header{
		"Location": []string{"main.php"},
	}
	mock.resp = getMockHTTPResponse("main.html", 302, headers)

	client := Client{
		UserName:         "123456",
		Password:         "secret",
		httpClient:       mock,
		initSuccessfully: false,
	}

	err := client.Authenticate()

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

func TestGetMainPhpText(t *testing.T) {
	mock := &MockHTTPClient{
		resp: getMockHTTPResponse("main.html", 200, http.Header{}),
	}
	client := Client{
		httpClient:       mock,
		initSuccessfully: true,
	}

	text, err := client.getMainPhpText()

	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	expectedText := getFileContent("main.html")
	if text != expectedText {
		t.Error("unexpected text")
	}

	if mock.url != "https://admin.instantproxies.com/main.php" {
		t.Errorf("unexpected url called: %v", mock.url)
		return
	}
}

func TestGetTextAreaInnerText(t *testing.T) {
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
