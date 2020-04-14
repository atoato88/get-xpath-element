package main

// This program refers to https://qiita.com/h-hiroki/items/04d8c6636968c07a438e
// Thanks to @h-hiroki

// This program use ChromeWebDriver
// See https://chromedriver.chromium.org/downloads
// And move program in download file to $PATH.

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sclevine/agouti"
	"github.com/spf13/pflag"
)

func main() {
	var driver *agouti.WebDriver
	var (
		debug bool
		file  string
		url   string
		xpath string
	)

	pflag.BoolVarP(&debug, "debug", "d", false, "If true, chrome window will open")
	pflag.StringVarP(&file, "file", "f", "", "Input file has xpath expression")
	pflag.Parse()
	args := pflag.Args()
	fmt.Println(args)

	// From Tokyo Sta. to Inadadudumi Sta. at starts on 2020/04/14 08:45
	url = "https://transit.yahoo.co.jp/search/result?flatlon=&fromgid=&from=%E6%9D%B1%E4%BA%AC&tlatlon=&togid=&to=%E7%A8%B2%E7%94%B0%E5%A0%A4&viacode=&via=&viacode=&via=&viacode=&via=&y=2020&m=04&d=14&hh=08&m2=5&m1=4&type=1&ticket=ic&expkind=1&ws=3&s=0&al=1&shin=1&ex=1&hb=1&lb=1&sr=1&kw=%E7%A8%B2%E7%94%B0%E5%A0%A4"
	if file == "" {
		xpath = "//div[@id='contents-body']/div[@id='main']//div[@id='srline']//dl/dd/ul/li//span[@class='mark' and contains(text(), '円')]"
	} else {
		fp, err := os.Open(file)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			xpath = scanner.Text()
			break
		}
	}

	if debug == true {
		driver = agouti.ChromeDriver()
	} else {
		driver = agouti.ChromeDriver(
			agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}),
		)
	}

	if err := driver.Start(); err != nil {
		log.Fatalf("driverの起動に失敗しました : %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("セッション作成に失敗しました : %v", err)
	}

	// ウェブページに遷移する
	if err := page.Navigate(url); err != nil {
		log.Fatalf("error occured when get page", err)
	}

	// 対象のXpathを取得する
	elem := page.FindByXPath(xpath)
	str, err := elem.Text()
	if err != nil {
		log.Fatalf("error occured when get xpath element", err)
	} else {
		fmt.Println(str)
	}

	os.Exit(0)

	//	// framesetの中の要素を検索するには一旦該当のフレームにフォーカスしなければならない
	//	// 任意のフレームにフォーカスする
	//	if err := page.FindByXPath("/html/frameset/frame[1]").SwitchToFrame(); err != nil {
	//		log.Fatalf("阿部寛の左側frameにフォーカスできませんでした : %v", err)
	//	}
	//	// 任意のフォーカス上の要素をクリック
	//	if err := page.FindByXPath("html/body/table/tbody/tr[10]/td[3]/p/a").Click(); err != nil {
	//		log.Fatalf("阿部寛の写真集が見つかりませんでした : %v", err)
	//	}
	//
	//	// フレームのフォーカス外すためrootにもどる
	//	if err := page.SwitchToRootFrame(); err != nil {
	//		log.Fatalf("しょうも無いエラーが発生しました : %v", err)
	//	}
	//
	//	// スクショとる
	//	if err := page.Screenshot("./abe_hiroshi.jpg"); err != nil {
	//		log.Fatalf("スクショ取れまへん : %v", err)
	//	}
}
