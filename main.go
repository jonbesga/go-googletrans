package googletrans

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Translator struct {
}

type Translated struct {
	Text string
}

func (t Translator) translate(text string, dest string, src string) ([]interface{}, error) {
	client := &http.Client{}
	url := "https://translate.googleapis.com/translate_a/single"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	for k, v := range build_params("gtx", text, src, dest, "xxxx") {
		switch val := v.(type) {
		case string:
			q.Add(k, val)
		case []string:
			for _, item := range val {
				q.Add(k, item)
			}
		}
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received non-200 response: " + resp.Status)
	}
	var parsedData []interface{}

	if err := json.Unmarshal(body, &parsedData); err != nil {
		return nil, err
	}
	return parsedData, nil
}

func (t Translator) Translate(text string, dest string, src string) (Translated, error) {
	dest = strings.SplitN(strings.ToLower(dest), "_", 2)[0]
	src = strings.SplitN(strings.ToLower(src), "_", 2)[0]

	if src != "auto" {
		if _, exists := LANGUAGES[src]; !exists {
			if special, ok := SPECIAL_CASES[src]; ok {
				src = special
			} else if langCode, ok := LANGCODES[src]; ok {
				src = langCode
			} else {
				return Translated{}, errors.New("invalid source language")
			}
		}
	}

	if _, exists := LANGUAGES[dest]; !exists {
		if special, ok := SPECIAL_CASES[dest]; ok {
			dest = special
		} else if langCode, ok := LANGCODES[dest]; ok {
			dest = langCode
		} else {
			return Translated{}, errors.New("invalid destination language")
		}
	}

	data, _ := t.translate(text, dest, src)
	translated := data[0].([]interface{})[0].([]interface{})[0]
	return Translated{Text: translated.(string)}, nil
}
