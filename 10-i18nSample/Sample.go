// Day10 | 讓我們的 Go 更國際化 - i18n 的應用
// https://ithelp.ithome.com.tw/articles/10235509
package main

import (
	"fmt"

	. "github.com/codingXiang/gogo-i18n"
	"golang.org/x/text/language"
)

func main() {
	// 可以透過 SetUseLanguage 方法進行更換，以下將預設語言更換為 英文
	GGi18n = NewGoGoi18n(language.TraditionalChinese)
	GGi18n.SetFileType("yaml")
	GGi18n.LoadTranslationFile("./i18n",
		language.TraditionalChinese,
		language.English)
	msg := GGi18n.GetMessage("welcome", map[string]interface{}{
		"username": "阿翔",
	})
	fmt.Println(msg)

	GGi18n.SetUseLanguage(language.English)
	msg = GGi18n.GetMessage("welcome", map[string]interface{}{
		"username": "阿翔",
	})
	fmt.Println(msg)
}
