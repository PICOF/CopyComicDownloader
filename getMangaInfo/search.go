package getMangaInfo

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
)

const nameWidth = 60
const authorWidth = 40

type Author struct {
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	PathWord string `json:"path_word"`
}
type ComicInfo struct {
	Name     string   `json:"name"`
	Alias    string   `json:"alias"`
	PathWord string   `json:"path_word"`
	Cover    string   `json:"cover"`
	Author   []Author `json:"author"`
	Popular  int      `json:"popular"`
}
type ComicResult struct {
	Result struct {
		List   []ComicInfo `json:"list"`
		Total  int         `json:"total"`
		Limit  int         `json:"limit"`
		Offset int         `json:"offset"`
	} `json:"results"`
}

func SearchComic(limit int, name string) ComicResult {
	get, err := PublicCli.Client.Do("https://www.copymanga.site/api/kb/web/searchs/comics?limit="+strconv.Itoa(limit)+"&q="+url.QueryEscape(name), PublicCli.Options, "GET")
	if err != nil {
		log.Println("Failed to get,error:", err)
		return ComicResult{}
	}
	if get.JSONBody()["code"].(float64) != 200 {
		log.Println("Failed to get,error:", get.JSONBody()["message"].(string))
		return ComicResult{}
	}
	results := get.Body
	var comicResult ComicResult
	err = json.Unmarshal([]byte(results), &comicResult)
	if err != nil {
		log.Println("Failed to unmarshal json,error:", err)
		return ComicResult{}
	}
	return comicResult
}

func (results ComicResult) GetSearchInfo() (info string) {
	info += "搜索结果共" + strconv.Itoa(int(math.Min(float64(results.Result.Limit), float64(results.Result.Total)))) + "条（受 limit 影响）：\n"
	info += fmt.Sprintf("%-4s%-"+strconv.Itoa(nameWidth-3)+"s%-"+strconv.Itoa(authorWidth-2)+"s%s\n", "序号", "漫画名", "作者", "热度")
	var author string
	var a, b int
	for i, v := range results.Result.List {
		author = authorFormat(v.Author)
		a = nameWidth - (len(v.Name)-len([]rune(v.Name)))/2
		b = authorWidth - (len(author)-len([]rune(author)))/2
		info += fmt.Sprintf("%-6d%-"+strconv.Itoa(a)+"s%-"+strconv.Itoa(b)+"s%d\n", i, v.Name, author, v.Popular)
	}
	return
}
func authorFormat(author []Author) (ret string) {
	for i, v := range author {
		if i == 0 {
			ret += v.Name
		} else {
			ret += "、" + v.Name
		}
	}
	return
}
