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
	path = "../"
	port = ":7880"
)

var (
	Year = 1980

	JieQi  = []Jq{}
	GanZhi = []Gz{}

	A = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	B = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	C = []string{"青", "明", "刑", "朱", "匱", "德", "白", "玉", "牢", "玄", "司", "勾"}

	HMap = map[int]int{
		23: 0,
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  2,
		5:  3,
		6:  3,
		7:  4,
		8:  4,
		9:  5,
		10: 5,
		11: 6,
		12: 6,
		13: 7,
		14: 7,
		15: 8,
		16: 8,
		17: 9,
		18: 9,
		19: 10,
		20: 10,
		21: 11,
		22: 11,
	}
)

type Jq struct {
	Year  string
	JieQi []Item
}

type Item struct {
	Name string
	Date time.Time
}

type Gz struct {
	Year string
	Item [][]string
}

func rfile_2080() error {
	f, err := os.Open(path + "2080.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	item := []Item{}

	Year = 1980
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		if line == "\n" {
			JieQi = append(JieQi, Jq{
				Year:  fmt.Sprintf("%d", Year),
				JieQi: item,
			})
			Year += 1
			item = []Item{}
		} else {
			i := strings.Split(line, " ")

			s := fmt.Sprintf("%s %s", i[1], i[2])
			r := []rune(s)
			r = append(r[:19], r[20:]...)

			t, _ := time.Parse("2006-01-02 15:04:05", string(r))

			item = append(item, Item{
				Name: i[0],
				Date: t,
			})
		}
	}
	return nil
}

func rfile_1980() error {
	f, err := os.Open(path + "1980.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	item := [][]string{}

	Year = 1980
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		if line == "\n" {
			GanZhi = append(GanZhi, Gz{
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
	return nil
}

func initMemory() error {
	err := rfile_2080()
	if err != nil {
		return err
	}
	err = rfile_1980()
	if err != nil {
		return err
	}
	return nil
}

func appointDate(t time.Time) []string {
	item := []Item{}
	j := Item{}

	for _, v := range JieQi {
		if v.Year == fmt.Sprint(t.Year()) {
			item = v.JieQi
			break
		}
	}
	for i := 0; i < 24; i += 2 {
		if item[i].Date.Month() == t.Month() {
			j = item[i]
			break
		}
	}

	var g []string

	for _, v := range GanZhi {
		if v.Year == fmt.Sprint(t.Year()) {
			g = v.Item[t.Month()-1]
			break
		}
	}
	for i := 0; i < t.Day()-1; i++ {
		g[2] = peek(g[2])
	}
	if t.After(j.Date) {
		g[1] = peek(g[1])
	}

	var h string

	switch string([]rune(g[2])[0]) {
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
		x, ok := HMap[t.Hour()]
		if !ok {
			panic("what's problem?")
		}
		if B[x] == string([]rune(h)[1]) {
			break
		}
		h = peek(h)
	}

	g = append(g, h)

	i := t.Minute()
	d := 0

	if i < 10 {
		d = 0
	}
	if i >= 10 && i < 20 {
		d = 1
	}
	if i >= 20 && i < 30 {
		d = 2
	}
	if i >= 30 && i < 40 {
		d = 3
	}
	if i >= 40 && i < 50 {
		d = 4
	}
	if i >= 50 {
		d = 5
	}
	if t.Hour()%2 == 0 {
		d += 6
	}
	g = append(g, B[d])

	return g
}

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {})

	// fmt.Println("listen port:", Port)

	// err := http.ListenAndServe(Port, mux)
	// if err != nil {
	// 	panic(err)
	// }
	err := initMemory()
	if err != nil {
		panic(err)
	}

	x, y, z := paipan(time.Now())
	fmt.Printf("%v \n%v \n%v \n%v \n", y, z, x, B)
}

func paipan(t time.Time) ([]string, []string, []string) {
	// t, _ := time.Parse("2006-01-02 15:04", "2021-11-09 22:23")
	date := appointDate(t)
	fmt.Println(t.Format("2006-01-02 15:04:05"), date)

	var x, y, z []string

	for i := 0; i < 12; i++ {
		x = append(x, " ")
		y = append(y, " ")
		z = append(z, " ")
	}

	a := index(B, string([]rune(date[3])[1]))
	b := index(B, string([]rune(date[4])))
	j := a

	for i := 0; i < 12; i++ {
		x[j] = B[b]

		j += 1
		b += 1

		if j == 12 {
			j = 0
		}
		if b == 12 {
			b = 0
		}
	}

	p := 0

	switch string([]rune(date[3])[0]) {
	case "甲":
		p = 5
	case "乙":
		p = 6
	case "丙":
		p = 7
	case "丁":
		p = 8
	case "戊":
		p = 9
	case "己":
		p = 0
	case "庚":
		p = 1
	case "辛":
		p = 2
	case "壬":
		p = 3
	case "癸":
		p = 4
	}

	j = a
	for i := 0; i < 12; i++ {
		y[j] = A[p]

		j += 1
		p += 1

		if j == 12 {
			j = 0
		}
		if p == 10 {
			p = 0
		}
	}
	j = a
	for i := 0; i < 12; i++ {
		z[j] = C[i]

		j += 1
		if j == 12 {
			j = 0
		}
	}

	return x, y, z
}

func index(p []string, t string) int {
	for i, v := range p {
		if v == t {
			return i
		}
	}
	panic("unexpected error")
}

func peek(s string) string {
	i := index(A, string([]rune(s)[0]))
	j := index(B, string([]rune(s)[1]))

	i += 1
	j += 1

	if i == 10 {
		i = 0
	}
	if j == 12 {
		j = 0
	}
	return fmt.Sprintf("%s%s", A[i], B[j])
}

/*
func stringer(x, y []string) {
	fmt.Printf(`
		    %s%s%s%s
		    %s%s%s%s
		    巳午未申
		%s%s辰    酉%s%s
		%s%s卯    戌%s%s
		    寅丑子亥
		    %s%s%s%s
		    %s%s%s%s

`, y[5], y[6], y[7], y[8], x[5], x[6], x[7], x[8],

		y[4], x[4], x[9], y[9],
		y[3], x[3], x[10], y[10],

		x[2], x[1], x[0], x[11], y[2], y[1], y[0], y[11])
}
*/
