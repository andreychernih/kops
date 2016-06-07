package fi

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"strconv"
	"strings"
)

// This file parses /etc/passwd and /etc/group to get information about users & groups
// Go has built-in user functionality, and group functionality is merged but not yet released
// TODO: Replace this file with e.g. user.LookupGroup once 42f07ff2679d38a03522db3ccd488f4cc230c8c2 lands in go 1.7

type User struct {
	Name    string
	Uid     int
	Gid     int
	Comment string
	Home    string
	Shell   string
}

func parseUsers() (map[string]*User, error) {
	users := make(map[string]*User)

	path := "/etc/passwd"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading user file %q", path)
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Split(line, ":")

		if len(tokens) < 7 {
			glog.Warning("Ignoring malformed /etc/passwd line (too few tokens): %q", line)
			continue
		}

		uid, err := strconv.Atoi(tokens[2])
		if err != nil {
			glog.Warning("Ignoring malformed /etc/passwd line (bad uid): %q", line)
			continue
		}
		gid, err := strconv.Atoi(tokens[3])
		if err != nil {
			glog.Warning("Ignoring malformed /etc/passwd line (bad gid): %q", line)
			continue
		}

		u := &User{
			Name: tokens[0],
			// Password
			Uid:     uid,
			Gid:     gid,
			Comment: tokens[4],
			Home:    tokens[5],
			Shell:   tokens[6],
		}
		users[u.Name] = u
	}
	return users, nil
}

func LookupUser(name string) (*User, error) {
	users, err := parseUsers()
	if err != nil {
		return nil, fmt.Errorf("error reading users: %v", err)
	}
	return users[name], nil
}

func LookupUserById(uid int) (*User, error) {
	users, err := parseUsers()
	if err != nil {
		return nil, fmt.Errorf("error reading users: %v", err)
	}
	for _, v := range users {
		if v.Uid == uid {
			return v, nil
		}
	}
	return nil, nil
}

type Group struct {
	Name string
	Gid  int
	//Members []string
}

func parseGroups() (map[string]*Group, error) {
	groups := make(map[string]*Group)

	path := "/etc/group"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading group file %q", path)
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Split(line, ":")

		if len(tokens) < 4 {
			glog.Warning("Ignoring malformed /etc/group line (too few tokens): %q", line)
			continue
		}

		gid, err := strconv.Atoi(tokens[2])
		if err != nil {
			glog.Warning("Ignoring malformed /etc/group line (bad gid): %q", line)
			continue
		}

		g := &Group{
			Name: tokens[0],
			// Password: tokens[1]
			Gid: gid,
			// Members: strings.Split(tokens[3], ",")
		}
		groups[g.Name] = g
	}
	return groups, nil
}

func LookupGroup(name string) (*Group, error) {
	groups, err := parseGroups()
	if err != nil {
		return nil, fmt.Errorf("error reading groups: %v", err)
	}
	return groups[name], nil
}

func LookupGroupById(gid int) (*Group, error) {
	users, err := parseGroups()
	if err != nil {
		return nil, fmt.Errorf("error reading groups: %v", err)
	}
	for _, v := range users {
		if v.Gid == gid {
			return v, nil
		}
	}
	return nil, nil
}