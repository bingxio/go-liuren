"""
    @title      获取 1980 年到 2080 年的 24 节气交节时间
    @author     bingxio, 丙杺
    @email      bingxio@qq.com
    @date       2021-10-25 23:16:12
"""

from typing import List, Tuple
import requests as rq

from bs4 import BeautifulSoup

URL = "http://jieqi.ttcha.net/"
YEAR = 1980

FILE = open("2080.txt", "w")


def wfile(result):
    for i in result:
        FILE.write("%s %s\n" % (i[0], i[1]))
    FILE.write("\n")


def parse(span) -> List[Tuple[str, str]]:
    pos = 0

    i, p = 0, 2
    r = []

    while pos < 24:
        name = span[i].text
        date = span[p].text

        r.append((name, date))

        i = p + 1
        p += 3
        pos += 1
    return r


for i in range(0, 100):
    rp = rq.get("%s%s.html" % (URL, YEAR))
    print("(%d) : %d --> %d => OK" % (i + 1, YEAR, rp.status_code))

    soup = BeautifulSoup(rp.text, "html.parser")
    span = soup.select("ul span")

    result = parse(span)
    wfile(result)

    YEAR += 1

FILE.close()
print("Done!")
