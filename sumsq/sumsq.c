#include <stdio.h>

int square(int num) { return num * num; }

int sum_squares(int x, int y) { return square(x) + square(y); }

int main() {
  printf("%d\n", sum_squares(2, 10));
  return 0;
}
