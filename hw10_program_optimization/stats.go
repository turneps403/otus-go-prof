package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + domain
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	var result users
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		var user User
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return result, err
		}
		result[i] = user
	}
	return result, nil
}

func countDomains(us users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, u := range us {
		dom := eDomain(u.Email)
		if strings.HasSuffix(dom, domain) {
			result[dom]++
		}
	}
	return result, nil
}

func eDomain(e string) string {
	for i := 0; i < len(e); i++ {
		if e[i] == '@' {
			return strings.ToLower(e[i+1:])
		}
	}
	return ""
}
