package imports

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ImportMap map[string]*string

func NewImportMap() *ImportMap {
	am := make(ImportMap)
	tmp := "jModels"
	am["github.com/andrskom/jrpc2hh/models"] = &tmp
	tmp1 := "json"
	am["encoding/json"] = &tmp1
	return &am
}

func (am *ImportMap) Register(i string) {
	if _, ok := (*am)[i]; !ok {
		(*am)[i] = nil
	}
}

func (am *ImportMap) GenerateAlias() {
	crossMap := make(map[string]*bool)
	for i, _ := range *am {
		var p string
		if (*am)[i] != nil {
			continue
		}
		if strings.Contains(i, "/") {
			path := strings.Split(i, "/")
			p = path[len(path)-1]
		} else {
			p = i
		}
		if _, ok := crossMap[p]; !ok {
			crossMap[p] = nil
		} else {
			tmp := fmt.Sprintf("%s_%d", p, randomNumberGenerator().Int())
			(*am)[i] = &tmp
		}
	}
}

func (am *ImportMap) GetFormattedAlias(name string) string {
	if (*am)[name] == nil {
		path := strings.Split(name, "/")
		return path[len(path)-1]
	} else {
		return *((*am)[name])
	}
}

func randomNumberGenerator() *rand.Rand {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1
}
