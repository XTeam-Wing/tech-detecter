package utils

import (
	"encoding/json"
	"fmt"
	"github.com/XTeam-Wing/tech-detecter/model"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strings"
)

func ParseYaml(path, filename string) (model.FingerPrint, error) {
	var fingerPrint = model.FingerPrint{}

	config := viper.New()
	config.AddConfigPath(path)
	config.SetConfigName(filename)
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		return fingerPrint, err
	}

	matches := config.Get("matches").(string)
	slice := LinesToSlice(matches)
	for _, line := range slice {
		if line != "" {
			fingerPrint.Matches = append(fingerPrint.Matches, line)
		}
	}
	data, err := json.Marshal(config.Get("info"))
	if err != nil {
		return fingerPrint, err
	}
	var info model.TechInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return fingerPrint, err
	}
	fingerPrint.Infos = info.Product
	return fingerPrint, nil

}

func LinesToSlice(str string) []string {
	toSlice := strings.Split(str, "\n")
	return toSlice
}

func GetCerts(resp *http.Response) []byte {
	var certs []byte
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		var str string
		if js, err := json.Marshal(cert); err == nil {
			certs = js
		}
		str = string(certs) + cert.Issuer.String() + cert.Subject.String()
		certs = []byte(str)
	}
	return certs
}

func GetTitle(content string) string {
	reTitle := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	matchResults := reTitle.FindAllString(content, -1)
	var nilString = ""
	var matches = []string{"<title>", "</title>"}
	return StringReplace(SliceToSting(matchResults), matches, nilString)
}

func StringReplace(old string, matches []string, new string) string {
	for _, math := range matches {
		old = strings.Replace(old, math, new, -1)
	}
	return old
}
func SliceToSting(slice []string) string {
	toString := fmt.Sprintf(strings.Join(slice, ","))
	return toString
}
