package x

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/0xDeSchool/gap/app"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func PostUrlForm(path string, body map[string]string, result any) error {
	c, ok := app.GetOptional[http.Client]()
	if !ok {
		c = http.DefaultClient
	}
	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}
	res, err := c.Post(path, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		cont, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(cont))
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(result)
}

func Post(url, contentType string, body any, result any) error {
	c, ok := app.GetOptional[http.Client]()
	if !ok {
		c = http.DefaultClient
	}
	content, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err := c.Post(url, contentType, bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		cont, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(cont))
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(result)
}

func Get(url string, headers map[string]string, result any) error {
	c, ok := app.GetOptional[http.Client]()
	if !ok {
		c = http.DefaultClient
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}
	res, err := c.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//conten, err := io.ReadAll(res.Body)
	//if err != nil {
	//	return err
	//}
	//log.Info(string(conten))
	return json.NewDecoder(res.Body).Decode(result)
}
