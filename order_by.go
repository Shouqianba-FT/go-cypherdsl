package go_cypherdsl

import (
	"errors"
	"fmt"
)

type OrderByConfig struct {
	Type string
	Name string
	Desc bool
}

func (o *OrderByConfig) ToString() (string, error) {
	if o.Name == "" {
		return "", errors.New("member have to be defined")
	}
	query := ""
	if o.Type != "" {
		query = fmt.Sprintf("%s.%s", o.Type, o.Name)
	} else {
		query = o.Name
	}

	if o.Desc {
		query += " DESC"
	}
	return query, nil
}
