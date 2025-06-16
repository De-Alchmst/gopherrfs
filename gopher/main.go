package gopher

import (
	"git.mills.io/prologic/go-gopher"
)

func FetchData(path string) ([]byte, error) {
	res, err := gopher.Get("gopher://"+path)
	if err != nil {
		return []byte(""), err
	}

	txt, err := res.Dir.ToText()
	if err != nil {
		return []byte(""), err
	}

	return txt, nil
}
