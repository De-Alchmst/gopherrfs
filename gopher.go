package main

import (
	"fmt"
	"errors"

	"git.mills.io/prologic/go-gopher"
)

type API struct {}

func (API) Read(address string, modifiers []string) ([]byte, error) {
	if len(modifiers) > 0 {
		return []byte(""), errors.New(fmt.Sprintf("unsuppordet modifier: %s", modifiers[0]))
	}

	res, err := gopher.Get("gopher://" + address)
	if err != nil {
		return []byte(""), err
	}

	txt, err := res.Dir.ToText()
	if err != nil {
		return []byte(""), err
	}

	return txt, nil
}

func (API) Write(address string, modifiers []string, data []byte) ([]byte, error) {
	return []byte(""), errors.New("not applicable")
}
