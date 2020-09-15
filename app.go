package main

import (
	"encoding/json"
	"fmt"
	"github.com/oniontree-org/go-oniontree"
	"github.com/urfave/cli/v2"
	"net/http"
	"net/url"
)

const Version = "0.1"

type Application struct {
	ot  *oniontree.OnionTree
	app *cli.App
}

func (a *Application) handleOnionTreeOpen() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ot, err := oniontree.Open(c.String("C"))
		if err != nil {
			return fmt.Errorf("failed to open OnionTree repository: %s", err)
		}
		a.ot = ot
		return nil
	}
}

func (a *Application) handleSyncCommand() cli.ActionFunc {
	newService := func(id string) *oniontree.Service {
		s := oniontree.NewService(id)
		s.Name = "Invidious"
		s.Description = "Invidious is an alternative front-end to YouTube"
		return s
	}
	return func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return fmt.Errorf("Missing a service ID")
		}

		client := &http.Client{
			Timeout: c.Duration("timeout"),
		}

		req, err := http.NewRequest("GET", c.String("url"), nil)
		if err != nil {
			return fmt.Errorf("failed to create new request: %s", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to establish connection with the API: %s", err)
		}
		defer resp.Body.Close()

		// Decode JSON payload
		var arr []interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&arr); err != nil {
			return fmt.Errorf("failed to decode response data: %s", err)
		}

		service, err := a.ot.GetService(id)
		if err != nil {
			if _, ok := err.(*oniontree.ErrIdNotExists); !ok {
				return fmt.Errorf("failed to get existing service data: %s", err)
			}
			service = newService(id)
		}

		normalizeURL := func(u string) (string, error) {
			nu, err := url.ParseRequestURI(u)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s://%s", nu.Scheme, nu.Host), nil
		}
		urls := []string{}

		for _, obj := range arr {
			instance := obj.([]interface{})[1].(map[string]interface{})
			if instance["type"] != "onion" {
				continue
			}
			u, err := normalizeURL(instance["uri"].(string))
			if err != nil {
				return fmt.Errorf("failed to normalize URL: %s", err)
			}
			urls = append(urls, u)
		}

		if c.Bool("replace") {
			service.SetURLs(urls)
		} else {
			service.AddURLs(urls)
		}

		if err := a.ot.AddService(service); err != nil {
			if _, ok := err.(*oniontree.ErrIdExists); !ok {
				return fmt.Errorf("failed to add new service: %s", err)
			}
			if err := a.ot.UpdateService(service); err != nil {
				return fmt.Errorf("failed to update service: %s", err)
			}
		}
		return nil
	}
}

func (a *Application) Run(args []string) error {
	return a.app.Run(args)
}
