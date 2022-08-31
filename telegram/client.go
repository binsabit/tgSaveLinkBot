package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	getUpdatesMthd  = "/getUpdates"
	sendMessageMthd = "/sendMessage"
)

type Client struct {
	host     string
	basepath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basepath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(t string) string {
	return "bot" + t
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   c.basepath + method,
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	//encoding incomming query to query of new request
	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (c *Client) GetUpdates(offset int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))

	response, err := c.doRequest(getUpdatesMthd, q)
	if err != nil {
		return nil, err
	}

	var updateResponse UpdateResponse
	err = json.Unmarshal(response, &updateResponse)
	if err != nil {
		return nil, err
	}

	return updateResponse.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMthd, q)
	if err != nil {
		return err
	}
	fmt.Println("sending Message")
	return nil
}
