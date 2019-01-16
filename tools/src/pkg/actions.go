package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	eos "github.com/eoscanada/eos-go"
)

func LoadActions(filepath string) (out []*eos.Action, err error) {
	fi, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		var ms *eos.Action
		err = json.Unmarshal(line, &ms)
		if nil != err {
			log.Printf("Error: %s\n", err)
		}

		out = append(out, ms)
	}
	return
}

func GetCachedActions(action string) (out []*eos.Action, err error) {
	totalName := fmt.Sprintf("%s-actions.json", action)
	totalActions, err := LoadActions(totalName)
	if err != nil {
		log.Printf("Get %s Error: %s\n", totalName, err)
	}

	log.Printf("Total actions %d", len(totalActions))

	successName := fmt.Sprintf("%s-success.json", action)
	successActions, err := LoadActions(successName)
	if err != nil {
		log.Printf("Get %s Error: %s\n", successName, err)
	}

	log.Printf("Total success actions %d", len(successActions))
	fl, err := os.Create(fmt.Sprintf("%s-actions-todo.json", action))
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	successActionStrs := make([]string, 0)
	for _, action := range successActionStrs {
		successActionStrs = append(successActionStrs, string(action))
	}

	for _, action := range totalActions {
		data, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("binary marshalling: %s", err)
		}

		actionStr := string(data)
		if ok, _ := Contain(actionStr, successActionStrs); ok == false {
			out = append(out, action, nil)

			fmt.Println(string(data))
			_, err = fl.Write(data)
			if err != nil {
				return nil, err
			}
			_, _ = fl.Write([]byte("\n"))
		}
	}
	return
}

func Contain(str string, list []string) (bool, error) {
	for _, s := range list {
		if s == str {
			return true, nil
		}
	}
	return false, nil
}
