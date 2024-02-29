package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	client *resty.Client
}

func NewClient() Client {
	return Client{
		client: resty.New().
			SetTimeout(1*time.Minute).
			SetHeader("Accpet", "application/json").
			SetLogger(log.StandardLogger()),
	}
}

func (c *Client) Get(url, query string) ([]byte, error) {
	res, err := c.client.R().
		EnableTrace().
		SetQueryString(query).
		Get(url)
	if err != nil {
		return nil, err
	}

	log.Info("requset success: [GET: ", url, "]")
	return res.Body(), nil
}

func (c *Client) Put(url string, param interface{}) error {
	res, err := c.client.R().
		EnableTrace().
		SetBody(param).
		Put(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf(string(res.Body()))
	}

	log.Info("requset success: [PUT: ", url, "]")
	return nil
}

func (c *Client) Post(url string, param interface{}) error {
	res, err := c.client.R().
		EnableTrace().
		SetBody(param).
		Post(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf(string(res.Body()))
	}

	log.Info("requset success: [POST: ", url, "]")
	return nil
}
