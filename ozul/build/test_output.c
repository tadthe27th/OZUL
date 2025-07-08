#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    int health = 100;
    double speed = 3.140000;
    char* name = "Pikachu";
    health = (health + 50);
    speed = (speed / 2.000000);
    char str_buffer_3[256];
    strcpy(str_buffer_3, "Hello, ");
    strcat(str_buffer_3, name);
    char* greeting = str_buffer_3;
    printf("%d\n", health);
    printf("%f\n", speed);
    printf("%s\n", greeting);
    int userInput;
    scanf("%d", &userInput);
    printf("%d\n", userInput);
    return 0;
}