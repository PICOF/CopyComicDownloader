package getMangaInfo

import (
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func init() {
	var client = PublicClient{}
	client.Client = cycletls.Init()
	client.Options = cycletls.Options{
		Timeout: 100,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Ja3:       "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0",
		UserAgent: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
	}
	PublicCli = &client
}
