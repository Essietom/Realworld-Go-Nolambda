package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type StringSet map[string]bool

func NewStringSetFromSlice(slice []string) StringSet {
	s := make(StringSet, len(slice))
	for _, value := range slice {
		s[value] = true
	}
	return s
}

func (s StringSet) Difference(other StringSet) StringSet {
	result := make(StringSet)
	for value := range s {
		if !other[value] {
			result[value] = true
		}
	}
	return result
}

func (s StringSet) ToSlice() []string {
	result := make([]string, 0, len(s))
	for value := range s {
		result = append(result, value)
	}
	return result
}

func ParseBody(r *http.Request, x interface{}) error{
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return err
		}
	}
	return nil
}