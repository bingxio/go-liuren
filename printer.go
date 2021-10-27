package main

import "fmt"

const (
	Text = `
	阳历：%s
	干支：%s
	空亡：%s、%s
	月将：%s

	   丙 丁 戊 己			丙 壬 乙 己

	   贵 腾 朱 六			腾 白 阴 勾

	   申 酉 戌 亥			申 寅 巳 亥

     乙 后 未       子 勾 庚		寅 申 亥 戊

     甲 阴 午       丑 青 辛

	   巳 辰 卯 寅			官    寅 白

	   玄 常 白 空			子 戊 申 蛇

	   乙 甲 癸 壬			官    寅 白

`
)

func Display(r Response) {
	ft := func(g []string) string {
		var s string
		for i, v := range g {
			s += v
			if i+1 != len(g) {
				s += " "
			}
		}
		return s
	}

	fmt.Printf(
		Text, r.Date, ft(r.Gz), r.Kw[0], r.Kw[1], r.J,
	)
}
