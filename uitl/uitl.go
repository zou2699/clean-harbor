/*
@Time : 2020/2/13 15:43
@Author : Tux
@Description :
*/

package uitl

import (
	"sort"
	"strings"

	"clean-harbor/model"
)

// 获取tag.Name[:3]作为prefix,每个prefix保留keepNum个,返回待清理的Tags
func FilterByPrefix(tags model.Tags, keepNum int) (filterTags model.Tags, err error) {
	mapTags := make(map[string]model.Tags)
	var prefix string
	for _, tag := range tags {
		// todo
		if strings.Contains(tag.Name, "test") {
			prefix = tag.Name[:5]
		} else {
			prefix = tag.Name[:4]
		}
		// fmt.Println("prefix:", prefix)
		if strings.HasPrefix(tag.Name, prefix) {
			mapTags[prefix] = append(mapTags[prefix], tag)
		}
	}

	for _, v := range mapTags {
		if len(v) > keepNum {
			sort.Sort(v)
			filterTags = append(filterTags, v[keepNum:]...)
		}
	}
	return filterTags, err
}
