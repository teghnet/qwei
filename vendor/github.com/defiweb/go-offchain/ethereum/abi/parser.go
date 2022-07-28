// Copyright © 2022 The go-offchain Authors
// Copyright © 2022 The go-ethereum Authors
// This file is an integral part of a package and it is a re-implementation of parts of the go-ethereum library.
//
// This package is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This package is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"fmt"
	"strings"
)

type SignatureMarshaling struct {
	Name    string               `json:"name"`
	Type    string               `json:"type"`
	Inputs  []ArgumentMarshaling `json:"inputs"`
	Outputs []ArgumentMarshaling `json:"outputs"`
}

type ArgumentMarshaling struct {
	Name         string
	Type         string
	InternalType string
	Components   []ArgumentMarshaling
	Indexed      bool
}

// parse converts a method selector into a struct that can be JSON encoded
// and consumed by other functions in this package.
// Note, although uppercase letters are not part of the ABI spec, this function
// still accepts it as the general format is valid.
func parse(sig string) (SignatureMarshaling, error) {
	if len(sig) == 0 {
		return SignatureMarshaling{}, fmt.Errorf("empty token")
	}
	fir, las, err := find(sig)
	if err != nil {
		return SignatureMarshaling{}, fmt.Errorf("failed to parse selector '%s': %w", sig, err)
	}

	name, in, out := strings.Trim(sig[:fir], " "), strings.Trim(sig[fir:las], " "), ""

	inputs, err := assemble(in)
	if err != nil {
		return SignatureMarshaling{}, fmt.Errorf("failed to assemble input args: %w", err)
	}

	if len(sig) > las {
		sig = sig[las:]
		fir, las, err = find(sig)
		if err != nil {
			return SignatureMarshaling{}, fmt.Errorf("failed to find returns '%s': %w", sig, err)
		}
		out = strings.Trim(sig[fir:las], " ")
	}

	outputs, err := assemble(out)
	if err != nil {
		return SignatureMarshaling{}, fmt.Errorf("failed to assemble output args: %w", err)
	}

	return SignatureMarshaling{name, "function", inputs, outputs}, nil
}

func find(s string) (first int, last int, err error) {
	depth := 0
	for i, c := range s {
		if c == '(' {
			if depth == 0 {
				first = i
			}
			depth++
		} else if c == ')' {
			if depth == 0 {
				return 0, 0, fmt.Errorf("closing paren before opening '%s'", s)
			}
			depth--
			if depth == 0 {
				return first, i + 1, nil
			}
		}
	}
	if depth > 0 {
		return 0, 0, fmt.Errorf("not enough closing parens '%s'", s)
	}
	return 0, 0, fmt.Errorf("no parens '%s'", s)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentifierSymbol(c byte) bool {
	return c == '$' || c == '_'
}

func parseToken(unescapedSelector string, isIdent bool) (string, string, error) {
	if len(unescapedSelector) == 0 {
		return "", "", fmt.Errorf("empty token")
	}
	firstChar := unescapedSelector[0]
	position := 1
	if !isAlpha(firstChar) && (!isIdent || !isIdentifierSymbol(firstChar)) {
		return "", "", fmt.Errorf("invalid token start: %c", firstChar)
	}
	for position < len(unescapedSelector) {
		char := unescapedSelector[position]
		if !(isAlpha(char) || isDigit(char) || (isIdent && isIdentifierSymbol(char))) {
			break
		}
		position++
	}
	return unescapedSelector[:position], unescapedSelector[position:], nil
}

func parseElementaryType(unescapedSelector string) (string, string, error) {
	parsedType, rest, err := parseToken(unescapedSelector, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse elementary type: %v", err)
	}
	for len(rest) > 0 && rest[0] == '[' {
		parsedType = parsedType + string(rest[0])
		rest = rest[1:]
		for len(rest) > 0 && isDigit(rest[0]) {
			parsedType = parsedType + string(rest[0])
			rest = rest[1:]
		}
		if len(rest) == 0 || rest[0] != ']' {
			return "", "", fmt.Errorf("failed to parse array: expected ']', got %c", unescapedSelector[0])
		}
		parsedType = parsedType + string(rest[0])
		rest = rest[1:]
	}
	return parsedType, rest, nil
}

func parseCompositeType(unescapedSelector string) ([]interface{}, string, error) {
	if len(unescapedSelector) == 0 {
		return nil, "", fmt.Errorf("empty composite type")
	}
	if unescapedSelector[0] != '(' {
		return nil, "", fmt.Errorf("expected '(', got %c", unescapedSelector[0])
	}
	parsedType, rest, err := parseType(unescapedSelector[1:])
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse type: %v", err)
	}
	result := []interface{}{parsedType}
	for len(rest) > 0 && rest[0] != ')' {
		parsedType, rest, err = parseType(rest[1:])
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse type: %v", err)
		}
		result = append(result, parsedType)
	}
	if len(rest) == 0 || rest[0] != ')' {
		return nil, "", fmt.Errorf("expected ')', got '%s'", rest)
	}
	if len(rest) >= 3 && rest[1] == '[' && rest[2] == ']' {
		return append(result, "[]"), rest[3:], nil
	}
	return result, rest[1:], nil
}

func parseType(unescapedSelector string) (interface{}, string, error) {
	if len(unescapedSelector) == 0 {
		return nil, "", fmt.Errorf("empty type")
	}
	if unescapedSelector[0] == '(' {
		return parseCompositeType(unescapedSelector)
	} else {
		return parseElementaryType(unescapedSelector)
	}
}

func assemble(in string) ([]ArgumentMarshaling, error) {
	if len(in) == 0 || in == "()" {
		return []ArgumentMarshaling{}, nil
	}
	inputArgs, _, err := parseCompositeType(in)
	if err != nil {
		return []ArgumentMarshaling{}, err
	}

	return assembleArgs(inputArgs)
}

func assembleArgs(args []interface{}) ([]ArgumentMarshaling, error) {
	arguments := make([]ArgumentMarshaling, 0)
	for i, arg := range args {
		// generate dummy name to avoid unmarshal issues
		name := fmt.Sprintf("name%d", i)
		if s, ok := arg.(string); ok {
			arguments = append(arguments, ArgumentMarshaling{name, s, s, nil, false})
		} else if components, ok := arg.([]interface{}); ok {
			subArgs, err := assembleArgs(components)
			if err != nil {
				return nil, fmt.Errorf("failed to assemble components: %v", err)
			}
			tupleType := "tuple"
			if len(subArgs) != 0 && subArgs[len(subArgs)-1].Type == "[]" {
				subArgs = subArgs[:len(subArgs)-1]
				tupleType = "tuple[]"
			}
			arguments = append(arguments, ArgumentMarshaling{name, tupleType, tupleType, subArgs, false})
		} else {
			return nil, fmt.Errorf("failed to assemble args: unexpected type %T", arg)
		}
	}
	return arguments, nil
}
