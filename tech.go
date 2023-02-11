package tech_detecter

import (
	"fmt"
	"github.com/XTeam-Wing/tech-detecter/cel"
	"github.com/XTeam-Wing/tech-detecter/model"
	"github.com/XTeam-Wing/tech-detecter/utils"
	"github.com/google/cel-go/common/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TechDetecter struct {
	// Apps is organized as <name, fingerprint>
	FinerPrint []model.FingerPrint
}

func (t *TechDetecter) Init(rulePath string) error {
	if !utils.Exists(rulePath) {
		return os.ErrNotExist
	}
	if utils.IsDir(rulePath) {
		files := utils.ReadDir(rulePath)
		for _, file := range files {
			rule, err := utils.ParseYaml(rulePath, file)
			if err != nil {
				fmt.Println(fmt.Sprintf("file %s error:%s", file, err))
				continue
			}
			t.FinerPrint = append(t.FinerPrint, rule)
		}
	}
	return nil
}

func (t *TechDetecter) Detect(response *http.Response) (string, error) {

	options := cel.InitCelOptions()
	env, err := cel.InitCelEnv(&options)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)
	headerInfo := ""
	for k, v := range response.Header {
		headerInfo += fmt.Sprintf("%v: %v\n", k, v[0])
	}
	var product []string
	for _, r := range t.FinerPrint {
		var matches string

		for i, match := range r.Matches {
			if i < len(r.Matches)-1 {
				matches = matches + "(" + match + ") || "
			} else {
				matches = matches + "(" + match + ")"
			}
		}
		ast, iss := env.Compile(matches)
		if iss.Err() != nil {
			log.Println(fmt.Sprintf("[X] product: %s rule Compile error", r.Infos))
			continue
		}
		prg, err := env.Program(ast)
		if err != nil {
			log.Println(fmt.Sprintf("[X] product: %s rule prg error:%s", r.Infos, err.Error()))
			continue
		}
		out, _, err := prg.Eval(map[string]interface{}{
			"body":     string(body),
			"title":    utils.GetTitle(string(body)),
			"header":   headerInfo,
			"server":   fmt.Sprintf("server: %v\n", response.Header["Server"]),
			"cert":     string(utils.GetCerts(response)),
			"banner":   headerInfo,
			"protocol": "",
			"port":     "",
		})
		if err != nil {
			log.Println(fmt.Sprintf("[X] product: %s rule Eval error:%s", r.Infos, err.Error()))

			continue
		}

		if out.(types.Bool) {
			product = append(product, r.Infos)
		}
	}
	return utils.SliceToSting(product), nil

}
