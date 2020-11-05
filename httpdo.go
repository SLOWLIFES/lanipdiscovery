package lanipdiscovery

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type HttpDo struct {
	Url      string
	Data     map[string]string
	FileKey  string
	FileName string
	File     io.Reader
}

func (that HttpDo) Request(method string) ([]byte, error) {
	var r http.Request
	r.ParseForm()
	for e := range that.Data {
		r.Form.Add(e, that.Data[e])
	}

	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest(method, that.Url, strings.NewReader(bodystr))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that HttpDo) Post() ([]byte, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if that.FileKey != "" {
		part, err := writer.CreateFormFile(that.FileKey, that.FileName)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, that.File)
	}

	for e := range that.Data {
		_ = writer.WriteField(e, that.Data[e])
	}
	err := writer.Close()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", that.Url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that HttpDo) Get() ([]byte, error) {
	return that.Request("GET")
}

func Get(url string, p string) string {

	if strings.Index(url, "https://") == -1 && strings.Index(url, "http://") == -1 {
		return ""
	}

	if p != "" {
		url = fmt.Sprintf("%s?%s", url, p)
	}
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}

	return string(body)
}

func Post(url string, p string) string {
	if strings.Index(url, "https://") == -1 && strings.Index(url, "http://") == -1 {
		return ""
	}

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(p))

	if err != nil {
		// handle error
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}

	return string(body)

}

func PostRaw(url string, data string) (string, error) {
	if strings.Index(url, "https://") == -1 && strings.Index(url, "http://") == -1 {
		return "", errors.New("url error")
	}

	body := &bytes.Buffer{}
	body.Write([]byte(data))
	request, err := http.NewRequest("POST", url, body)

	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", err
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	return string(result), err
}
