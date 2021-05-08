package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type printer func(io.Reader) error

func stringPrinter(reader io.Reader) error {
	contains, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error("Error when read from reader for json printer: ", err)
		return err
	}
	fmt.Println(contains)
	return nil
}

func jsonPrinter(reader io.Reader) error {
	contains, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error("Error when read from reader for json printer: ", err)
		return err
	}
	var str bytes.Buffer
	err = json.Indent(&str, contains, "", "\t")
	if err != nil {
		log.Error("Error when indent json to string: ", err)
		return err
	}
	fmt.Println(str.String())
	return nil
}

func sendRequest(path string, values map[string]string, port int, output printer) error{
	apiUrl := fmt.Sprintf("http://127.0.0.1:%d", port)

	data := url.Values{}
	if values != nil {
		for k, v := range values {
			data.Set(k, v)
		}
	}

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = path

	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error("Error when sent request to server: ", err)
		return err
	}

	defer resp.Body.Close()
	err = output(resp.Body)
	if err != nil {
		log.Error("Error when read response from server: ", err)
	}
	return nil
}
