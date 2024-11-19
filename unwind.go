package go_cypherdsl

import (
	"errors"
	"fmt"
	"strings"
)

type UnwindConfig struct {
	Slice interface{}
	As    string
}

func (u *UnwindConfig) ToString() (string, error) {
	if u.Slice == nil {
		return "", errors.New("slice in unwind can not be empty")
	}

	if u.As == "" {
		return "", errors.New("AS has to be defined")
	}

	switch v := u.Slice.(type) {
	case []interface{}:
		query := "["

		for _, i := range v {
			str, err := cypherizeInterface(i)
			if err != nil {
				return "", err
			}

			query += fmt.Sprintf("%s,", str)
		}

		query = strings.TrimSuffix(query, ",")
		return query + fmt.Sprintf("] AS %s", u.As), nil
	default:
		query, err := cypherizeInterface(v)
		if err != nil {
			return "", err
		}
		return query + " AS " + u.As, nil
	}
}
