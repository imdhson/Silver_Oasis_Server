package modules

import (
	"errors"
	"fmt"
	"strings"
)

func (i *Vars_on_html) VarsOnHTML(wwwfile []byte) []byte { //<go>변수명</go> 사이에 실제 변수값을 넣어서 써주는 함수입니다.
	wwwfile_str := string(wwwfile)
	for {
		seek_start := strings.Index(wwwfile_str, "<go>") //<go>를 찾아서 첫 인덱스를 반환
		if seek_start != -1 {                            //<go>를 찾았으면
			seek_end := strings.Index(wwwfile_str[seek_start:], "</go>") //</go>가 어디인지 파악
			if seek_end < 0 {                                            //올바르지 않은 문서형식
				err := errors.New("HtmlInlineReplace에서 오류 발생. <go>변수이름</go>를 올바르게 사용했는지 확인필요")
				Critical(err)
			}
			varname := wwwfile_str[seek_start+4 : seek_start+seek_end]
			vardata := i.VarsMap[varname]
			//fmt.Println(seek_start, seek_end)
			//fmt.Printf("[%v]:%v\n", varname, vardata)
			maked := wwwfile_str[:seek_start] + vardata.(string) + wwwfile_str[seek_start+seek_end+5:] //</go>가 5글자여서 5글자 더함
			wwwfile_str = maked
		} else {
			return []byte(wwwfile_str)
		}
	}
}

type Vars_on_html struct {
	VarsMap map[string]interface{}
}

func (i *Vars_on_html) Init() {
	i.VarsMap = make(map[string]interface{})
}

func (i *Vars_on_html) AddVar(name string, data interface{}) {
	i.VarsMap[name] = data
}
func (i *Vars_on_html) Display() {
	for index, data := range i.VarsMap {
		fmt.Printf("[%v]: %v %T\n", index, data, data)
	}
}
