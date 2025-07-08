#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    int health = 100;
    double damage = 25.500000;
    double newHealth = (health - damage);
    printf("%f\n", newHealth);
    return 0;
}