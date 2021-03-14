package deepl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mcaubrey/sakigake_rss_local/internal/secret"
)

func Translate(text string, source string, target string) (string, error) {
	var dlRes map[string]interface{}
	t := getDeepLURL(strings.ReplaceAll(text, "\n", ""), target, source)

	log.Print(t)

	_ = getJson(t, &dlRes)
	body, err := json.Marshal(dlRes)
	if err != nil {
		return "", err
	}

	jsonMap := deepLResponse{}
	jErr := json.Unmarshal(body, &jsonMap)
	if jErr != nil {
		return "", err
	}

	log.Print(jsonMap)

	result := ""
	for _, s := range jsonMap.Translations {
		result = result + s.Text
	}
	return result, err
}

type deepLResponse struct {
	Translations []deepLTranslations `json:"translations"`
}

type deepLTranslations struct {
	Text string `json:"text"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(t string, target interface{}) error {
	r, err := myClient.Get(t)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getDeepLURL(text string, target string, source string) (t string) {
	t = fmt.Sprintf("https://api.deepl.com/v2/translate?auth_key=%v&split_sentences=1&text=%v", secret.SecretAPIKey, text)

	if source != "" {
		t = t + "&source_lang=" + source
	}

	if target != "" {
		t = t + "&target_lang=" + target
	}

	return t + "&tag_handling=xml"
}
