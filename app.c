/*
  小干时局，排盘系统
  丙杺著于，辛丑年秋

  bingxio@qq.com
*/
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

static const char *A[] = {
    "甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"};

static const char *B[] = {
    "子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"};

FILE *fp_2080 = NULL;
FILE *fp_1980 = NULL;

typedef enum { TA, TB } type;

typedef struct {
    int year;
    int month;
    int day;
    int hour;
} date;

typedef struct {
    int day;
    int hour;
} jq;

typedef struct {
    uint8_t g1[2];
    uint8_t g2[2];
    uint8_t g3[2];
    uint8_t g4[2];
} gz;

typedef struct {
    char *date;
    uint8_t i1, i2, i3, i4;
    gz *gz;
    uint8_t k1, k2;
} env;

typedef struct {
    uint8_t h1, h2;
} hp_map;

static const hp_map hplist[] = {
    {23, 24},
    {1,  2 },
    {3,  4 },
    {5,  6 },
    {7,  8 },
    {9,  10},
    {11, 12},
    {13, 14},
    {15, 16},
    {17, 18},
    {19, 20},
    {21, 22},
};

void panic(const char *m) {
    printf("%s\n", m);
    exit(EXIT_SUCCESS);
}

date *parse_input(char *d) {
    char *tp = malloc(5);
    date *pd = malloc(sizeof(date));

    uint8_t status = 0;

    for (int i = 0, j = 0; i < 128; i++) {
        char c = d[i];

        if (c == ' ' || c == '\0') {
            int trans = atoi(tp);
            if (trans == 0) {
                panic("请输入十进制整数");
            }

            switch (status) {
            case 0:
                pd->year = trans;
                break;
            case 1:
                pd->month = trans;
                break;
            case 2:
                pd->day = trans;
                break;
            case 3:
                pd->hour = trans;
                break;
            }
            status++;
            memset(tp, 0, 5);
            j = 0;

            if (c == '\0') {
                break;
            }
        } else {
            tp[j++] = c;
        }
    }

    free(tp);
    return pd;
}

jq *parse_2080(date *dp) {
    fp_2080 = fopen("2080.txt", "r");
    if (fp_2080 == NULL) {
        panic("打开文件失败");
    }

    char line[100];

    while (!feof(fp_2080)) {
        fgets(line, 100, fp_2080);

        char year[5];
        for (int i = 7, j = 0; i < 11; i++) {
            year[j++] = line[i];
        }
        year[4] = '\0';

        char month[3];
        month[0] = line[12];
        month[1] = line[13];
        month[2] = '\0';

        if (atoi(year) == dp->year && atoi(month) == dp->month) {
            break;
        }
    }
    if (feof(fp_2080)) {
        panic("检查日期是否合法");
    }
    jq *q = malloc(sizeof(jq));
    // printf("--> %s", line);

    char day[3];
    char hour[3];

    day[0] = line[15];
    day[1] = line[16];
    day[2] = '\0';

    hour[0] = line[18];
    hour[1] = line[19];
    hour[2] = '\0';

    q->day = atoi(day);
    q->hour = atoi(hour);

    fclose(fp_2080);
    return q;
}

uint8_t index_of(type t, char *w) {
    const char **arr = NULL;
    int len;

    if (t == TA) {
        arr = A;
        len = 10;
    }
    if (t == TB) {
        arr = B;
        len = 12;
    }

    for (int i = 0; i < len; i++) {
        if (strcmp(arr[i], w) == 0) {
            return i;
        }
    }
    panic("未知符号");
}

gz *parse_1980(date *dp) {
    fp_1980 = fopen("1980.txt", "r");
    if (fp_1980 == NULL) {
        panic("打开文件失败");
    }

    char line[100];

    while (!feof(fp_1980)) {
        fgets(line, 100, fp_1980);

        char year[5];
        for (int i = 0; i < 4; i++) {
            year[i] = line[i];
        }
        year[4] = '\0';

        if (atoi(year) == dp->year) {
            for (int i = 1; i <= 12 && i != dp->month; i++) {
                fgets(line, 100, fp_1980);
            }
            break;
        }
    }
    // printf("--> %s", line);

    gz *g = malloc(sizeof(gz));
    int status = 1, left = 1;

    for (int i = 0, j = 5; i < 6; i++) {
        char pp[4] = {};
        strncpy(pp, line + j, 3);

        uint8_t *arr = NULL;

        switch (status) {
        case 1:
            arr = g->g1;
            break;
        case 2:
            arr = g->g2;
            break;
        case 3:
            arr = g->g3;
            break;
        }

        if (left) {
            arr[0] = index_of(TA, pp);
            left = 0;
        } else {
            arr[1] = index_of(TB, pp);
            left = 1;
            status++;
        }
        memset(pp, 0, 4);

        if (i % 2 == 0) {
            j += 3;
        } else {
            j += 4;
        }
    }

    return g;
}

void peek_gz(uint8_t *kv) {
    kv[0]++;
    kv[1]++;

    uint8_t a = kv[0];
    uint8_t b = kv[1];

    if (a == 10) {
        kv[0] = 0;
    }
    if (b == 12) {
        kv[1] = 0;
    }
}

char *printer(uint8_t *kv) {
    char *con = malloc(20);
    sprintf(con, "%s%s", A[kv[0]], B[kv[1]]);
    return con;
}

env *eval(date *dp, jq *q, gz *g, uint8_t p) {
    if (dp->day >= q->day && dp->hour >= q->hour) {
        peek_gz(g->g2);
    }
    for (int i = 1; i != dp->day; i++) {
        peek_gz(g->g3);
    }
    if (dp->hour == 23) {
        peek_gz(g->g3);
    }

    int hp = 0;
    for (; hp < 12; hp++) {
        hp_map m = hplist[hp];

        if (dp->hour == m.h1 || dp->hour == m.h2) {
            break;
        }
    }

    uint8_t gp;
    switch (g->g3[0]) {
    case 0:
    case 5:
        gp = 0;
        break;
    case 1:
    case 6:
        gp = 2;
        break;
    case 2:
    case 7:
        gp = 4;
        break;
    case 3:
    case 8:
        gp = 6;
        break;
    case 4:
    case 9:
        gp = 8;
        break;
    }

    uint8_t arr[2];
    arr[0] = gp;
    arr[1] = 0;

    for (int i = 0; i < 12 && arr[1] != hp; i++) {
        peek_gz(arr);
    }

    g->g4[0] = arr[0];
    g->g4[1] = arr[1];

    arr[0] = gp;
    arr[1] = 0;

    for (int i = 0; i < 12 && arr[1] != p; i++) {
        peek_gz(arr);
    }

    uint8_t i3 = arr[0];

    env *e = malloc(sizeof(env));

    char *date = malloc(128);
    sprintf(date, "%s %s %s %s", printer(g->g1), printer(g->g2), printer(g->g3),
        printer(g->g4));
    e->date = date;

    uint8_t hdex = 5;
    for (int i = 0; i < 10 && g->g4[0] != i; i++) {
        hdex++;
        if (hdex == 10) {
            hdex = 0;
        }
    }

    e->i1 = hdex;
    e->i2 = i3;
    e->i3 = p;
    e->i4 = g->g4[1];

    if (g->g3[0] == 9) {
        if (g->g3[1] == 11) {
            e->k1 = 0;
            e->k2 = 1;
        } else {
            e->k1 = g->g3[1] + 1;
            e->k2 = g->g3[1] + 2;
        }
    } else {
        uint8_t a = g->g3[0];
        uint8_t b = g->g3[1];

        while (a != 9) {
            a++;
            b++;

            if (b == 12) {
                b = 0;
            }
        }
        if (b == 11) {
            e->k1 = 0;
            e->k2 = 1;
        }
    }

    return e;
}

void save(env *e, char *name) {
    char file_name[20] = {};

    struct stat filestat;

    if (stat("./tmp", &filestat) != 0) {
        system("mkdir tmp");
    }

    sprintf(file_name, "./tmp/%s.json", name);

    FILE *jp = fopen(file_name, "w");

    char y[10];
    sprintf(y, "%s%s", A[e->gz->g1[0]], B[e->gz->g1[1]]);

    char m[10];
    sprintf(m, "%s%s", A[e->gz->g2[0]], B[e->gz->g2[1]]);

    char d[10];
    sprintf(d, "%s%s", A[e->gz->g3[0]], B[e->gz->g3[1]]);

    char h[10];
    sprintf(h, "%s%s", A[e->gz->g4[0]], B[e->gz->g4[1]]);

    char i1[10];
    sprintf(i1, "%s", A[e->i1]);

    char i2[10];
    sprintf(i2, "%s%s", A[e->i2], B[e->i3]);

    char i3[10];
    sprintf(i3, "%s", B[e->i4]);

    char k1[10];
    char k2[10];

    sprintf(k1, "%s", B[e->k1]);
    sprintf(k2, "%s", B[e->k2]);

    fprintf(jp,
        "{ \n\t\"gz\": [\"%s\", \"%s\", \"%s\", \"%s\"], \n\t\"st\": [\"%s\", "
        "\"%s\", "
        "\"%s\"], \n\t\"kw\": [\"%s\", \"%s\"] \n}\n",
        y, m, d, h, i1, i2, i3, k1, k2);

    fclose(jp);
}

// ./app 2021 11 4 11 3 -json 123
int main(int argc, char **argv) {
    char *input_date = malloc(128);
    char *df = malloc(5);
    char *fname;

    bool save_json = false;

    if (argc > 1) {
        memset(input_date, 0, 128);

        for (int i = 1; i <= 4; i++) {
            strcat(input_date, argv[i]);

            if (i + 1 != 5) {
                strcat(input_date, " ");
            }
        }
        memcpy(df, argv[5], 5);

        if (argc == 8) {
            save_json = strcmp(argv[6], "-json") == 0;
            fname = argv[7];
        }
    } else {
        printf("年月日时 yyyy MM dd hh\t\t:\t");

        scanf("%[^\n]", input_date);
        getchar();

        printf("请输入 1 - 12 之间的一个数字\t:\t");

        scanf("%[^\n]", df);
        getchar();
    }

    uint8_t p = (uint8_t)atoi(df);
    if (p == 0) {
        panic("请输入整数");
    }

    date *d = parse_input(input_date);

    jq *q = parse_2080(d);
    gz *g = parse_1980(d);

    env *e = eval(d, q, g, p - 1);
    e->gz = g;

    if (save_json) {
        save(e, fname);
    } else {
        printf("\
干支：%s\n\
局时：\n\
\
               %s\n\
             %s%s                       小干时（排盘工具）\n\
               %s\n\
",
            e->date, A[e->i1], A[e->i2], B[e->i3], B[e->i4]);
    }

    free(e->date);

    free(input_date);
    free(df);
}