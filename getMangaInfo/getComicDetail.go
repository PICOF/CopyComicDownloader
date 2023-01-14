package getMangaInfo

import (
	"bytes"
	"encoding/json"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sese/AESDecrypt"
)

var PublicCli *PublicClient

const KEY = "xxxmanga.woo.key"
const PageEachRow = 6

type PublicClient struct {
	Client  cycletls.CycleTLS
	Options cycletls.Options
}
type Chapter struct {
	Type int    `json:"type"`
	Name string `json:"name"`
	Id   string `json:"id"`
}
type CopyMangaInfo struct {
	Build struct {
		PathWord string `json:"path_word"`
		Type     []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"build"`
	Groups struct {
		Default struct {
			PathWord    string    `json:"path_word"`
			Count       int       `json:"count"`
			Name        string    `json:"name"`
			Chapters    []Chapter `json:"chapters"`
			LastChapter struct {
				Index           int    `json:"index"`
				Uuid            string `json:"uuid"`
				Count           int    `json:"count"`
				Ordered         int    `json:"ordered"`
				Size            int    `json:"size"`
				Name            string `json:"name"`
				ComicId         string `json:"comic_id"`
				ComicPathWord   string `json:"comic_path_word"`
				GroupPathWord   string `json:"group_path_word"`
				Type            int    `json:"type"`
				ImgType         int    `json:"img_type"`
				News            string `json:"news"`
				DatetimeCreated string `json:"datetime_created"`
				Prev            string `json:"prev"`
			} `json:"last_chapter"`
		} `json:"default"`
	} `json:"groups"`
}
type ComicPic struct {
	Url string `json:"url"`
}

func GetComicDetails(query string) CopyMangaInfo {
	get, err := PublicCli.Client.Do("http://www.copymanga.site/comicdetail/"+query+"/chapters", PublicCli.Options, "GET")
	if err != nil {
		log.Println("Failed to get,error:", err)
		return CopyMangaInfo{}
	}
	if get.JSONBody()["code"].(float64) != 200 {
		log.Println("Failed to get,error:", get.JSONBody()["message"].(string))
		return CopyMangaInfo{}
	}
	results := get.JSONBody()["results"].(string)
	res, err := AESDecrypt.AesDecrypt(results[16:], []byte(KEY), []byte(results[:16]))
	if err != nil {
		log.Println("Failed to decrypt,error:", err)
		return CopyMangaInfo{}
	}
	var info CopyMangaInfo
	err = json.Unmarshal(res, &info)
	if err != nil {
		log.Println("Failed to unmarshal json,error:", err)
		return CopyMangaInfo{}
	}
	return info
}

func (info CopyMangaInfo) GetList() map[string]string {
	var results = make(map[string]string)
	var types = make(map[int]string)
	for _, v := range info.Build.Type {
		types[v.Id] = v.Name
	}
	for i, v := range info.Groups.Default.Chapters {
		results[types[v.Type]] += v.Name + "\t"
		if (i+1)%PageEachRow == 0 {
			results[types[v.Type]] += "\n"
		}
	}
	return results
}
func (info CopyMangaInfo) GetChaptersInfo(chapterId string) []ComicPic {
	get, err := PublicCli.Client.Do("http://www.copymanga.site/comic/"+info.Build.PathWord+"/chapter/"+chapterId, PublicCli.Options, "GET")
	if err != nil {
		log.Println("Failed to get,error:", err)
		return nil
	}
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(get.Body)))
	if err != nil {
		log.Println("Failed to get new goquery document,error:", err)
		return nil
	}
	data, exist := reader.Find(".imageData").Attr("contentkey")
	if !exist {
		log.Println("Cannot find info about chapter,error:", err)
		return nil
	}
	res, err := AESDecrypt.AesDecrypt(data[16:], []byte(KEY), []byte(data[:16]))
	var result []ComicPic
	err = json.Unmarshal(res, &result)
	if err != nil {
		log.Println("Failed to unmarshal json,error:", err)
		return nil
	}
	return result
}
