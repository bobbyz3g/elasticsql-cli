package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

type query struct {
	Stmt string `json:"query"`
}

type CLI struct {
	url    string
	auth   []string
	client *http.Client
}

func NewCLI() *CLI {
	return &CLI{
		url:    "",
		client: http.DefaultClient,
	}
}

func (c *CLI) Ping(url string, auth []string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(auth[0], auth[1])
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 300 {
		return errors.New(fmt.Sprintf("response not ok: %s", resp.Status))
	}

	c.url = fmt.Sprintf("%s/_sql?format=txt", url)
	c.auth = auth
	fmt.Println("Success connect to", url)
	return nil
}

func (c *CLI) Run() error {
	p := prompt.New(
		c.executor,
		c.completer,
	)
	p.Run()
	return nil
}

var suggestions = []prompt.Suggest{
	// Command
	{"exit", "Exit CLI"},
	{"show tables", "Show all tables"},
	{"select", "Query data"},
	{"desc", "Describe table"},
}

func (c *CLI) completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func (c *CLI) executor(in string) {
	stmt := strings.TrimRight(strings.TrimSpace(in), ";")
	block := strings.Split(stmt, " ")

	switch block[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "select", "desc", "show":
		c.exec(stmt)
	default:
		fmt.Println("Unknown command")
	}
}

func (c *CLI) buildBody(stmt string) (io.Reader, error) {
	q := query{Stmt: stmt}
	data, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func (c *CLI) exec(stmt string) {
	body, err := c.buildBody(stmt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req, err := http.NewRequest("POST", c.url, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(data))
}
