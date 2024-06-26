package graphql

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (gql *subscriberGql) SendRequest(server string, payload []byte) (string, error) {
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close body: %v", err)
		}
	}()

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// MarshalGQLData processes event data into proper format acceptable by graphql
func (gql *subscriberGql) MarshalGQLData(gqlData interface{}) (string, error) {
	data, err := json.Marshal(gqlData)
	if err != nil {
		return "", err
	}

	// process the marshalled data to make it graphql compatible
	processed := strconv.Quote(string(data))
	processed = strings.Replace(processed, `\"`, `\\\"`, -1)
	return processed, nil
}
