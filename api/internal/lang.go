package internal

import (
	"os"
	"strings"

	gowords "github.com/saleh-rahimzadeh/go-words"
	core "github.com/saleh-rahimzadeh/go-words/core"
)

const stringSource string = `
MySQL Version: %s, Connection Count: %d_EN = MySQL Version: %s, Connection Count: %d
MySQL Version: %s, Connection Count: %d_ZH = MySQL版本: %s, 连接数: %d
  
Redis: Connected to %s_EN = Redis: Connected to %s
Redis: Connected to %s_ZH = Redis: 链接到 %s


Hello world_EN = Hello world
Hello world_ZH = 你好，世界


He"s world_EN = He"s world
He"s world_ZH = He"s 世界
`

type MyWords struct {
	Words gowords.Words
}

func (t *MyWords) Get(key string) string {
	v, ok := t.Words.Find(key)
	if ok {
		return v
	} else {
		return key
	}
}
func NewMyWords() *MyWords {
	return &W
}

func (t *MyWords) T(key string) string {
	v, ok := t.Words.Find(key)
	if ok {
		return v
	} else {
		return key
	}
}

func (t *MyWords) EN() {
	t.Words = wordsEN
}
func (t *MyWords) ZH() {
	t.Words = wordsZH
}

var W MyWords
var wordsEN gowords.Words
var wordsZH gowords.Words

func init() {
	const EN core.Suffix = "_EN"
	const ZH core.Suffix = "_ZH"

	var wrd gowords.WordsRepository
	wrd, err := gowords.NewWordsRepository(stringSource, core.Separator, core.Comment)
	if err != nil {
		panic(err)
	}

	wordsEN, err = gowords.NewWithSuffix(wrd, EN)
	if err != nil {
		panic(err)
	}

	wordsZH, err = gowords.NewWithSuffix(wrd, ZH)
	if err != nil {
		panic(err)
	}

	var gw gowords.Words
	if locale := os.Getenv("LANG"); strings.Contains(locale, "zh_CN") {
		gw = wordsZH
	} else {
		gw = wordsEN
	}
	W.Words = gw

	// println(mygw.Get("Hello world"))
	// println(mygw.Get("He\"s world"))
}

// func mainx() {

// 	const EN core.Suffix = "_EN"
// 	const FA core.Suffix = "_ZH"
// 	var err error
// 	var fileSource *os.File
// 	fileSource, err = os.Open("lang.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer fileSource.Close()

// 	wrd, err := gowords.NewWordsFile(fileSource, core.Separator, core.Comment)

// 	if err != nil {
// 		panic(err)
// 	}

// 	wordsEN, err := gowords.NewWithSuffix(wrd, EN)
// 	if err != nil {
// 		panic(err)
// 	}

// 	wordsZH, err := gowords.NewWithSuffix(wrd, FA)
// 	if err != nil {
// 		panic(err)
// 	}

// 	value1en := wordsEN.Get("key1")
// 	println(value1en)

// 	value2en, found2en := wordsEN.Find("key2")
// 	println(value2en, found2en)

// 	value1fa := wordsZH.Get("key1")
// 	println(value1fa)

// 	value2fa, found2fa := wordsZH.Find("key2")
// 	println(value2fa, found2fa)

// 	println(wordsEN.Get("Hello world"))
// 	println(wordsZH.Get("Hello world"))

// 	var gw gowords.Words
// 	if locale := os.Getenv("LANG"); strings.Contains(locale, "zh_CN") {
// 		gw = wordsZH
// 	} else {
// 		gw = wordsEN
// 	}

// 	var mygw MyWords
// 	mygw.Words = gw

// 	println(mygw.Get("Hello world"))
// 	println(mygw.Get("Hello 111 world"))

// }
