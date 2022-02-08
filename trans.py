from datetime import datetime

YEAR = 2020

P = -1
B = []

for line in open("data.txt"):
    if line == "\n":
        continue

    line = line[:-1]
    P += 1

    if P == 24:
        print(B, end=",\n")

        P = 0
        YEAR += 1
        B = []

    if len(B) == 0:
        B.append(YEAR)

    INNER = []
    INNER.append(line[:2])  # 名称

    t = datetime.strptime(line[3:], "%Y-%m-%d %H:%M:%S")

    INNER.append(t.month)
    INNER.append(t.day)
    INNER.append(t.hour)

    B.append(INNER)

    # if YEAR == 2023:
    #     break

    if YEAR == 2079 and P + 1 == 24:
        print(B)
