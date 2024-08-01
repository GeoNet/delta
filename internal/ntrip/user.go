package ntrip

import (
	"fmt"
	"sort"
	"strings"
)

const (
	userUsername int = iota
	userGroups
	userPassword
	userLast
)

// User represents a registered account user.
type User struct {
	Username string
	Groups   []string
	Password string
}

// decode will extract user entries from a strings slice, as expected in a csv formatted file.
func (u *User) decode(row []string) error {
	if l := len(row); l != userLast {
		return fmt.Errorf("incorrect \"user\" \"%s\": found %d items, expected %d", strings.Join(row, ","), l, userLast)
	}

	var groups []string
	for _, g := range strings.Split(row[userGroups], ":") {
		groups = append(groups, strings.TrimSpace(g))
	}

	sort.Strings(groups)

	*u = User{
		Username: strings.TrimSpace(row[userUsername]),
		Groups:   groups,
		Password: strings.TrimSpace(row[userPassword]),
	}

	return nil
}

// encode will format user entries into a strings slice, as expected in a csv formatted file.
func (u User) encode() []string {
	var row []string

	row = append(row, u.Username)
	row = append(row, strings.Join(u.Groups, ":"))
	row = append(row, u.Password)

	return row
}

// Users represents a list of account users.
type Users []User

func ReadUsers(path string) ([]User, error) {
	var users Users
	if err := ReadFile(path, &users); err != nil {
		return nil, err
	}

	check := make(map[string]User)
	for _, u := range users {
		if _, ok := check[u.Username]; ok {
			return nil, fmt.Errorf("duplicate user name: %s", u.Username)
		}
		check[u.Username] = u
	}

	sort.Slice(users, func(i, j int) bool { return users[i].Username < users[j].Username })

	return users, nil
}

// Header is a function to build a string slice as expected for the first line of a csv file.
func (u Users) Header() []string {
	return []string{"#Username", "Groups", "Password"}
}

// Fields returns the number of expected csv fields.
func (u Users) Fields() int {
	return userLast
}

// Decode is a function to extract information from a list of string slices, as expected for the entries in a csv file.
func (u *Users) Decode(data [][]string) error {
	for _, v := range data {
		var user User
		if err := user.decode(v); err != nil {
			return err
		}
		*u = append(*u, user)
	}
	return nil
}

// Encode is a function to build a list of string slices as expected for the entries in a csv file.
func (u Users) Encode() [][]string {
	var items [][]string

	for _, user := range u {
		items = append(items, user.encode())
	}

	return items
}
