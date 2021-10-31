/*
  小干时局，排盘系统
  丙杺著于，辛丑年秋
  
  bingxio@qq.com
*/
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

const char *A[] = {
  "甲", "乙", "丙", "丁", "戊",
  "己", "庚", "辛", "壬", "癸"
};

const char *B[] = {
  "子", "丑", "寅", "卯", "辰", "巳",
  "午", "未", "申", "酉", "戌", "亥"
};

FILE *fp_2080 = NULL;
FILE *fp_1980 = NULL;

typedef enum { TA, TB } type;

typedef struct {
  char *date;
  uint8_t i1, i2, i3, i4;
} env;

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

void panic(const char *m) {
  printf("%s\n", m);
  exit(EXIT_SUCCESS);
}

date *parse_input(char *d) {
  char *tp = malloc(sizeof(char) * 4);
  date *pd = malloc(sizeof(date));

  uint8_t status = 0;

  for (int i = 0, j = 0; i < 128; i++) {
    char c = d[i];

    if (c == ' ' || c == '\0') {
      // printf("--> %s\n", tp);

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
      memset(tp, 0, 4);
      j = 0;
      if (c == '\0') {
        break;
      }
    } else {
      tp[j++] = c;
    }
  }
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

    // printf("%s", line);
    // printf("--> %s %s | \n\n", year, month);

    if (atoi(year) == dp->year &&
        atoi(month) == dp->month) {
          break;
    }
  }
  if (feof(fp_2080)) {
    panic("检查日期是否合法");
  }
  jq *q = malloc(sizeof(jq));
  printf("%s", line);

  char day[3];
  char hour[3];

  day[0] = line[15];
  day[1] = line[16];
  day[2] = '\0';

  hour[0] = line[18];
  hour[1] = line[19];
  hour[2] = '\0';

  printf("%s %s\n", day, hour);
  q->day = atoi(day);
  q->hour = atoi(hour);

  fclose(fp_2080);
  return q;
}

uint8_t index_of(type t, char *w) {
  switch (t) {
    case TA:
      break;
    case TB:
      break;
  }
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
      for (int i = 0; i < 12; i++) {
        fgets(line, 100, fp_1980);

        if (i + 1 == dp->month) {
          break;
        }
      }
      break;
    }
  }

  printf("%s", line);

  gz *g = malloc(sizeof(gz));
  int status;
  char item[5];

  item[0] = line[5];
  item[1] = line[6];

  item[2] = line[7];
  item[3] = line[8];

  item[4] = '\0';

  printf("%s\n", item);

  // for (int i = 0; i < 4; i++) {

  // }
  return g;
}

int main(int argc, char **argv) {
  printf("年月日时 yyyy MM dd hh\t:\t");

  char *input_date = malloc(128);
  scanf("%[^\n]", input_date);
  getchar();

  date *d = parse_input(input_date);
  printf("%d 年 %d 月 %d 日 %d 时\n",
    d->year, d->month, d->day, d->hour);

  jq *q = parse_2080(d);
  gz *g = parse_1980(d);

  env *e = malloc(sizeof(env));
  printf("%d\n", e == NULL);
}
