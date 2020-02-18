/*
@Time : 2020/2/13 15:42
@Author : Tux
@Description :
*/

package model

import (
	"time"
)

// Project /api/projects?name=cloud
type Project struct {
	Name string `json:"name"`
	ID   int    `json:"project_id"`
}

// Repo /api/repositories?project_id=2
type Repo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Tag /api/repositories/cloud/demojava/tags  test-bcd2e9d
type Tag struct {
	Size    int64     `json:"size"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// Tags implement the sort interface
type Tags []Tag

func (t Tags) Len() int           { return len(t) }
func (t Tags) Less(i, j int) bool { return t[i].Created.After(t[j].Created) }
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
