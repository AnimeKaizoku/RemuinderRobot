package parser

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrNoMatchFound = errors.New("no match found")

type Parser struct {
	Regexes map[string]*regexp.Regexp
}

func NewParser(regexes map[string]*regexp.Regexp) *Parser {
	return &Parser{Regexes: regexes}
}

func (p *Parser) Parse(input string) (string, map[string]string, error) {
	for i := range p.Regexes {
		matches := p.Regexes[i].FindStringSubmatch(input)
		names := p.Regexes[i].SubexpNames()

		if len(matches) > 0 {
			fmt.Printf("found match with regex: %s - %v\n", i, p.Regexes[i])
			return i, mapSubexpNames(matches, names), nil
			// for k := range mapNames {
			// 	if k != "" {
			// 		fmt.Printf("%s: %s\n", k, mapNames[k])
			// 	}
			// }
			// break
		}
	}

	return "", nil, ErrNoMatchFound
}

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i := range n {
		r[n[i]] = m[i]
	}
	return r
}
