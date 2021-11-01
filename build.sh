CC=gcc-11
$CC -c -g app.c
$CC -fsanitize=address app.o -o app
rm -f *.o
