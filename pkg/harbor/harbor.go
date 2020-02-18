/*
@Time : 2020/2/13 15:53
@Author : Tux
@Description :
*/

package harbor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"clean-harbor/model"
)

// Client
type Client struct {
	Client  *http.Client
	BaseUrl string
}

// NewClient
func NewClient(username, password, baseUrl string) *Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				req.SetBasicAuth(username, password)
				return nil, nil
			},
		},
	}
	return &Client{
		Client:  client,
		BaseUrl: baseUrl,
	}
}

// getProjectID 更加projectName 获取projectID
func (c *Client) getProjectID(projectName string) (projectId int, err error) {
	resp, err := c.Client.Get(c.BaseUrl + "/api/projects?name=" + projectName)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code is:%v", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var projects []model.Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return
	}
	for _, p := range projects {
		if p.Name == projectName {
			return p.ID, nil
		}
	}
	// 若前面出错则验证
	if string(body) == "null" {
		err = errors.New("body is null,maybe login error")
		return
	}

	return 0, errors.New("not found")
}

// func (c *Client) GetRepo(projectId int) (repoNames []string, err error)
// GetRepoNames 内部调用getProjectID,并根据getProjectID获取project下面的repoNames,example:
// [cloud/demojava cloud/ecommerce-cloud-api-service cloud/ecommerce-cloud-automationtest-service]
func (c *Client) GetRepoNames(projectName string) (repoNames []string, err error) {
	projectId, err := c.getProjectID(projectName)
	if err != nil {
		return
	}
	resp, err := c.Client.Get(c.BaseUrl + "/api/repositories?project_id=" + strconv.Itoa(projectId))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var repos []model.Repo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return
	}
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

// GetRepoTags
func (c *Client) GetRepoTags(repo string) (tags model.Tags, err error) {
	resp, err := c.Client.Get(c.BaseUrl + "/api/repositories/" + repo + "/tags")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &tags)
	if err != nil {
		return
	}
	// 排序
	sort.Sort(tags)
	if !sort.IsSorted(tags) {
		return nil, errors.New("tags not sorted")
	}
	return
}

// DeleteRepoTag delete tags with repo name and tag.
func (c *Client) DeleteRepoTag(repo string, tag string) (err error) {
	request, err := http.NewRequest("DELETE", c.BaseUrl+"/api/repositories/"+repo+"/tags/"+tag, nil)
	if err != nil {
		return
	}
	resp, err := c.Client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("resp code=%v", resp.StatusCode)
		return
	}
	return
}
