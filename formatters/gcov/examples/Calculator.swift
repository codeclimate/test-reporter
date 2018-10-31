#include <stdio.h>
#include <string.h>

int compute(char* alpha, char* beta) {
	int distance = 0;
	int len = strlen(alpha)<strlen(beta) ? strlen(alpha) : strlen(beta);

	for(int i=0; i<len; i++) {
		if(alpha[i] != beta[i]) {
			distance++;
		}
	}

	return distance;
}

int main() {
	if(compute("", "AAGC") != 0) {
		printf("Test failed");
	}

	if(compute("AGA", "AGC") != 1) {
		printf("Test failed");
	}
}
