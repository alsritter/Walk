package main

import (
	"fmt"
	"regexp"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

func main() {
	// 模拟一个文件
	txt := []string{
		`"here is a string."`,
		"12345",
		"// hello world",
		"1 * 3 == alsritter",
	}

	for _, line := range txt {
		regex := *regexp.MustCompile(regexPat)
		pos := 0
		endPos := len(line)

		for pos < endPos {
			if s := regex.FindString(line[pos:]); s != "" {
				fmt.Println(matchToken(s, regex))
				pos = pos + len(s)
			}
		}
	}
}

func matchToken(str string, regexp regexp.Regexp) string {
	res := regexp.FindAllStringSubmatch(str, -1) // -1 表示匹配全部
	for i := range res {
		m := res[i][1]
		if m != "" { // if not a space
			if comment := res[i][2]; comment != "" {
				// 注释
				return comment
			}
			if number := res[i][3]; number != "" {
				// 数值
				return number
			}
			if str := res[i][4]; str != "" {
				// 字符串
				return str
			}
			// 符号
			return m
		}
	}

	return "空白"
}
