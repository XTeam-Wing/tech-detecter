package tech_detecter

import (
	"fmt"
	"github.com/XTeam-Wing/tech-detecter/cel"
	"github.com/XTeam-Wing/tech-detecter/model"
	"github.com/XTeam-Wing/tech-detecter/utils"
	"github.com/google/cel-go/common/types"
	"io/ioutil"
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
				return err
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
			return "", err
		}
		prg, err := env.Program(ast)
		if err != nil {
			return "", err
		}
		out, _, err := prg.Eval(map[string]interface{}{
			"body":   string(body),
			"title":  utils.GetTitle(string(body)),
			"header": headerInfo,
			"server": fmt.Sprintf("server: %v\n", response.Header["Server"]),
			"cert":   utils.GetCerts(response),
			"banner": "",
			//"protocol": "",
		})
		if err != nil {
			return "", err
		}

		if out.(types.Bool) {
			product = append(product, r.Infos)
		}
	}
	return utils.SliceToSting(product), nil

}
