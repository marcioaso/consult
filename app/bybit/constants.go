package bybit

import "fmt"

var DefaultHeaders = map[string]string{
	"Accept":             "application/json",
	"Content-Type":       "application/json",
	"Sec-CH-UA":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
	"Sec-CH-UA-Platform": `"Windows"`,
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

const api = "https://api2.bybit.com"

func getUrl(path string) string {
	var url = fmt.Sprintf("%v%v", api, path)
	return url
}
