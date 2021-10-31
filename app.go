/*
   @title      大六壬神课排盘系统（精确）
   @author     bingxio, 丙杺
   @email      bingxio@qq.com
   @date       2021-10-26 18:40:55
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	File_2080 = "2080.txt"
	File_1980 = "1980.txt"

	PORT = ":8080"
)

var (
	Year = 1980

	Data2080 = []Data_2080{}
	Data1980 = []Data_1980{}

	A = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	B = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	C = []string{"亥", "戌", "酉", "申", "未", "午", "巳", "辰", "卯", "寅", "丑", "子"}

	D = []string{"贵", "腾", "朱", "六", "勾", "青", "空", "白", "常", "玄", "阴", "后"}
	E = []string{"孙", "父", "兄", "鬼", "财"}
)

type Data_2080 struct {
	Year  string
	JieQi []Item_2080
}

type Item_2080 struct {
	Name string
	Date time.Time
}

type Data_1980 struct {
	Year string
	Item [][]string
}

type Class struct {
	A []string `json:"a"`
	B []string `json:"b"`
	C []string `json:"c"`
	D []string `json:"d"`
}

type Response struct {
	Date  string   `json:"date"`
	Gz    []string `json:"gz"`
	J     string   `json:"j"`
	Kw    []string `json:"kw"`
	Dp    []string `json:"dp"`
	Js    []string `json:"js"`
	Tp    []string `json:"tp"`
	Class Class    `json:"class"`
}

func (r Response) Stringer() {
	x := fmt.Sprintf("%s --> %+v %s", r.Date, r.Gz, r.J)
	fmt.Println(x)
}

func main() {
	a, err := rfile_2080()
	if err != nil {
		panic(err)
	}
	Data2080 = a

	b, err := rfile_1980()
	if err != nil {
		panic(err)
	}
	Data1980 = b

	/* http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		dt := r.URL.Query().Get("dt")

		if dt == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		e := evaluate(dt)
		b, err := json.Marshal(e)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = rw.Write(b)
	})

	fmt.Println("listen and serve http handler on port", PORT)

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		panic(err)
	} */

	t := time.Now().Format("2006-01-02 15:04")

	e := evaluate(t)
	// e.Stringer()
	Display(e)
}

func peekGz(s string, rev bool) string {
	x := indexOf(A, string([]rune(s)[0]))
	y := indexOf(B, string([]rune(s)[1]))

	if !rev {
		x++
		y++

		if x == 10 {
			x = 0
		}
		if y == 12 {
			y = 0
		}
	} else {
		x--
		y--

		if x == -1 {
			x = 9
		}
		if y == -1 {
			y = 11
		}
	}
	return fmt.Sprintf("%s%s", A[x], B[y])
}

func indexOf(element []string, value string) int {
	for i, v := range element {
		if v == value {
			return i
		}
	}
	panic("unexpected error")
}

func getKw(gz string) []string {
	p := []string{}
	for {
		gz = peekGz(gz, false)

		a := string([]rune(gz)[0])
		b := string([]rune(gz)[1])

		if a == "甲" {
			p = append(p, b)

			gz = peekGz(gz, false)
			b = string([]rune(gz)[1])

			p = append(p, b)
			break
		}
	}
	return p
}

func hourPeriods(t time.Time) string {
	h := t.Hour()
	switch h {
	case 23, 0:
		return B[0]
	case 1, 2:
		return B[1]
	case 3, 4:
		return B[2]
	case 5, 6:
		return B[3]
	case 7, 8:
		return B[4]
	case 9, 10:
		return B[5]
	case 11, 12:
		return B[6]
	case 13, 14:
		return B[7]
	case 15, 16:
		return B[8]
	case 17, 18:
		return B[9]
	case 19, 20:
		return B[10]
	case 21, 22:
		return B[11]
	}
	panic("unexpected error")
}

func evDp(hz string, j string) []string {
	dp := []string{}

	i := indexOf(B, j)
	p := indexOf(B, hz)

	x := i + (12 - p)
	if x > 12 {
		x = x - 12
	}
	for len(dp) != 12 {
		if x >= 12 {
			x = 0
		}
		dp = append(dp, B[x])
		x++
	}
	return dp
}

func evJs(hz string, dg string, dp []string) []string {
	js := make([]string, 12)

	fmt.Println("日干：", dg, "时支：", hz, dp)
	var a, b, p int

	switch dg {
	case "甲", "戊", "庚":
		a = 1
		b = 7
	case "乙", "己":
		a = 0
		b = 8
	case "丙", "丁":
		a = 11
		b = 9
	case "辛":
		a = 6
		b = 2
	case "壬", "癸":
		a = 5
		b = 3
	}

	switch indexOf(B, hz) {
	case 3, 4, 5, 6, 7, 8:
		p = a
	case 9, 10, 11, 0, 1, 2:
		p = b
	}

	fmt.Println("贵人：", B[p])
	var reverse bool

	switch indexOf(dp, B[p]) {
	case 11, 0, 1, 2, 3, 4:
		reverse = false
	case 5, 6, 7, 8, 9, 10:
		reverse = true
	}

	i := indexOf(dp, B[p])

	for j := 0; j < 12; j++ {
		js[i] = D[j]

		if reverse {
			i--
		} else {
			i++
		}
		if i == -1 {
			i = 11
		}
		if i == 12 {
			i = 0
		}
	}

	fmt.Println(js)
	return js
}

func evTp(dgz string, dp []string) []string {
	tp := make([]string, 12)

	var xz string

	for {
		dgz = peekGz(dgz, true)

		a := string([]rune(dgz)[0])
		b := string([]rune(dgz)[1])

		if a == "甲" {
			xz = b
			break
		}
	}

	fmt.Println("旬：", "甲"+xz)
	i := indexOf(dp, xz)
	j := 0

	for a := 0; a < 12; a++ {
		tp[i] = A[j]

		j++
		i++
		if j == 10 {
			j = 0
		}
		if i == 12 {
			i = 0
		}
	}

	fmt.Println(tp)
	return tp
}

func evClass(dgz string, dp []string, js []string, tp []string) Class {
	a := make([]string, 4)
	b := make([]string, 4)
	c := make([]string, 4)
	d := make([]string, 4)

	d[1] = string([]rune(dgz)[1])
	d[3] = string([]rune(dgz)[0])

	var i int

	switch d[3] {
	case "甲":
		i = 2
	case "乙":
		i = 4
	case "丙", "戊":
		i = 5
	case "丁", "己":
		i = 7
	case "庚":
		i = 8
	case "辛":
		i = 10
	case "壬":
		i = 11
	case "癸":
		i = 1
	}

	d[2] = dp[i]
	c[3] = dp[i]

	i = indexOf(B, dp[i])
	c[2] = dp[i]

	i = indexOf(B, string([]rune(dgz)[1]))
	d[0] = dp[i]
	c[1] = dp[i]

	i = indexOf(B, dp[i])
	c[0] = dp[i]

	for j := 0; j < 4; j++ {
		i = indexOf(dp, c[j])
		b[j] = js[i]
		a[j] = tp[i]
	}

	fmt.Println(a, b, c, d)
	return Class{a, b, c, d}
}

func EvalGZ(date string) ([]string, Item_2080, Item_2080) {
	t, _ := time.Parse("2006-01-02 15:04", date)

	gz := []string{}
	jq := []Item_2080{}
	jqItemA, jqItemB := Item_2080{}, Item_2080{}

	for _, v := range Data1980 {
		if v.Year == fmt.Sprint(t.Year()) {
			gz = v.Item[t.Month()-1]
		}
	}
	for _, v := range Data2080 {
		if v.Year == fmt.Sprint(t.Year()) {
			jq = v.JieQi
		}
	}
	for i := 0; i < 24; i += 2 {
		if jq[i].Date.Month() == t.Month() {
			jqItemA = jq[i]
			jqItemB = jq[i+1]
			break
		}
	}
	pre := gz[2]

	for i := 0; i < int(t.Day())-1; i++ {
		pre = peekGz(pre, false)
	}
	if t.After(jqItemA.Date) {
		gz[1] = peekGz(gz[1], false)
	}
	gz[2] = pre

	var h string
	switch string([]rune(pre)[0]) {
	case "甲", "己":
		h = "甲子"
	case "乙", "庚":
		h = "丙子"
	case "丙", "辛":
		h = "戊子"
	case "丁", "壬":
		h = "庚子"
	case "戊", "癸":
		h = "壬子"
	}
	for i := 0; i < 12; i++ {
		if hourPeriods(t) == string([]rune(h)[1]) {
			if t.Hour() == 23 {
				gz[2] = peekGz(gz[2], false)
			}
			break
		} else {
			h = peekGz(h, false)
		}
	}
	gz = append(gz, h)

	return gz, jqItemA, jqItemB
}

func evaluate(date string) Response {
	t, _ := time.Parse("2006-01-02 15:04", date)
	gz, _, jqItemB := EvalGZ(date)

	var j string
	p := 2
	for i := 0; i < 12; i++ {
		if B[p] == string([]rune(gz[1])[1]) {
			j = C[i]
		}
		p++
		if p == 12 {
			p = 0
		}
	}

	if t.Before(jqItemB.Date) {
		tp := indexOf(B, j)
		tp++
		if tp == 12 {
			tp = 0
		}
		j = B[tp]
	}

	hz := string([]rune(gz[3])[1])
	dg := string([]rune(gz[2])[0])

	kw := getKw(gz[2])
	dp := evDp(hz, j)
	js := evJs(hz, dg, dp)
	tp := evTp(gz[2], dp)
	cl := evClass(gz[2], dp, js, tp)

	return Response{date, gz, j, kw, dp, js, tp, cl}
}

func rfile_2080() ([]Data_2080, error) {
	f, err := os.Open(File_2080)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	item := []Item_2080{}
	data := []Data_2080{}

	Year = 1980
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		if line == "\n" {
			data = append(data, Data_2080{
				Year:  fmt.Sprintf("%d", Year),
				JieQi: item,
			})
			Year += 1
			item = []Item_2080{}
		} else {
			i := strings.Split(line, " ")

			s := fmt.Sprintf("%s %s", i[1], i[2])
			r := []rune(s)
			r = append(r[:19], r[20:]...)

			t, _ := time.Parse("2006-01-02 15:04:05", string(r))

			item = append(item, Item_2080{
				Name: i[0],
				Date: t,
			})
		}
	}
	return data, nil
}

func rfile_1980() ([]Data_1980, error) {
	f, err := os.Open(File_1980)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	item := [][]string{}
	data := []Data_1980{}

	Year = 1980
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		if line == "\n" {
			data = append(data, Data_1980{
				Year: fmt.Sprintf("%d", Year),
				Item: item,
			})
			Year += 1
			item = [][]string{}
		} else {
			y := line[5:11]
			m := line[12:18]
			d := line[19:25]

			item = append(item, []string{y, m, d})
		}
	}
	return data, nil
}
