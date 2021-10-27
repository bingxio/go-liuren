/*
   @title      1980 - 2080 年干支、节气（中气获取月将）
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

	TypeA = iota
	TypeB
)

var (
	Year = 1980

	Data2080 = []Data_2080{}
	Data1980 = []Data_1980{}

	A = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	B = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	C = []string{"亥", "戌", "酉", "申", "未", "午", "巳", "辰", "卯", "寅", "丑", "子"}
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

type Response struct {
	Date   string   `json:"date"`
	GanZhi []string `json:"ganzhi"`
	Jiang  string   `json:"jiang"`
}

func (r Response) Stringer() {
	x := fmt.Sprintf("%s --> %+v %s", r.Date, r.GanZhi, r.Jiang)
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
	e := evaluate("2021-10-27 16:37")
	e.Stringer()
}

func indexOf(t int, value string) int {
	p := A
	if t == TypeB {
		p = B
	}
	for i, v := range p {
		if v == value {
			return i
		}
	}
	panic("unexpected error")
}

func getJiang(yue string) string {
	p := 2
	for i := 0; i < 12; i++ {
		if B[p] == yue {
			return C[i]
		}
		p++
		if p == 12 {
			p = 0
		}
	}
	panic("unexpected error")
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

func evaluate(et string) Response {
	t, _ := time.Parse("2006-01-02 15:04", et)

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

	next := func(s string) string {
		x := indexOf(TypeA, string([]rune(s)[0])) + 1
		y := indexOf(TypeB, string([]rune(s)[1])) + 1

		if x == 10 {
			x = 0
		}
		if y == 12 {
			y = 0
		}
		return fmt.Sprintf("%s%s", A[x], B[y])
	}
	pre := gz[2]

	for i := 0; i < int(t.Day())-1; i++ {
		pre = next(pre)
	}
	if t.After(jqItemA.Date) {
		gz[1] = next(gz[1])
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
		h = next(h)

		if hourPeriods(t) == string([]rune(h)[1]) {
			break
		}
	}
	gz = append(gz, h)

	j := getJiang(string([]rune(gz[1])[1]))

	if t.Before(jqItemB.Date) {
		tp := indexOf(TypeB, j)
		tp++
		if tp == 12 {
			tp = 0
		}
		j = B[tp]
	}
	return Response{
		Date:   et,
		GanZhi: gz,
		Jiang:  j,
	}
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
