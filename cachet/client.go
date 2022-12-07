package cachet

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andygrunwald/cachet"
)

type Config struct {
	Addr  string
	Token string
}

type client struct {
	cli *cachet.Client

	componentID map[string]int

	m sync.Mutex
}

var globalClient *client

func Init(cfg Config) error {
	cli, err := cachet.NewClient(cfg.Addr, &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 8,
		},
		Timeout: time.Second,
	})
	if err != nil {
		return err
	}
	cli.Authentication.SetTokenAuth(cfg.Token)

	resp, _, err := cli.Components.GetAll(&cachet.ComponentsQueryParams{})
	if err != nil {
		return err
	}

	componentID := make(map[string]int)
	for _, c := range resp.Components {
		componentID[c.Name] = c.ID
	}

	globalClient = &client{cli: cli,
		componentID: componentID,
	}

	return nil
}

func (c *client) create(component string) error {
	if _, ok := c.componentID[component]; ok {
		return nil
	}

	resp, _, err := c.cli.Components.Create(&cachet.Component{
		Name: component,
	})
	if err != nil {
		return fmt.Errorf("create component err: %s", err)
	}

	c.componentID[component] = resp.ID
	return nil
}

func (c *client) UpdateStatus(component string, state State) error {
	c.m.Lock()
	defer c.m.Unlock()

	if err := c.create(component); err != nil {
		return fmt.Errorf("create component %s err: %s", component, err)
	}

	id := c.componentID[component]
	_, _, err := c.cli.Components.Update(id, &cachet.Component{
		Status: int(state),
	})
	if err != nil {
		return fmt.Errorf("update component state err: %s", err)
	}
	return nil
}

func UpdateStatus(component string, state State) error {
	return globalClient.UpdateStatus(component, state)
}
