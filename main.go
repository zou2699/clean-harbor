/*
@Time : 2020/2/13 15:41
@Author : Tux
@Description :
*/

package main

import (
	"flag"
	"fmt"
	"log"

	"clean-harbor/pkg/harbor"
	"clean-harbor/uitl"
)

var (
	url         string
	user        string
	password    string
	projectName string
	keepNum     int
	help        bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help message")
	flag.StringVar(&url, "url", "", "harbor地址")
	flag.StringVar(&user, "user", "", "harbor账号")
	flag.StringVar(&password, "password", "", "harbor密码")
	flag.StringVar(&projectName, "projectName", "", "projectName")
	flag.IntVar(&keepNum, "keepNum", 5, "每个repo保留的tag个数")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if url == "" || user == "" || password == "" || projectName == "" {
		flag.Usage()
		return
	}
	harborClient := harbor.NewClient(user, password, url)

	// 获取cloud下面所有repo的名字
	repoNames, err := harborClient.GetRepoNames(projectName)
	if err != nil {
		log.Fatalf("GetRepoNames: %s\n", err)
	}

	var size int64
	for _, repoName := range repoNames {
		// 根据repo的获取其下所有的tag
		tags, err := harborClient.GetRepoTags(repoName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("repo: %s 当前的tag数量为: %d\n", repoName, len(tags))

		// 自动获取prefix,并根据prefix分类排序过滤
		filterTags, err := uitl.FilterByPrefix(tags, keepNum)
		if err != nil {
			panic(err)
		}
		// fmt.Println(filterTags)
		for _, tag := range filterTags {
			fmt.Printf("删除image: %s:%s, 创建时间为: %s\n", repoName, tag.Name, tag.Created)
			err := harborClient.DeleteRepoTag(repoName, tag.Name)
			if err != nil {
				fmt.Printf("image: %s:%s DeleteRepoTag: %s\n", repoName, tag.Name, err)
				continue
			}
			size += tag.Size
		}
	}

	fmt.Printf("本次共清理: %.2f MB\n", float64(size)/1024/1024)

}
