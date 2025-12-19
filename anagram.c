#include <stdio.h>
#include <stdlib.h>

#define INITIAL_SIZE 128

int excluded(char character) {
    int ignoreList[] = {',', ' '};
    int length = sizeof(ignoreList) / sizeof(ignoreList[0]);
    int i;

    for (i = 0; i < length; i++) {
        if (ignoreList[i] == character) {
            return 1;
        }
    }
    return 0;
}

int read_line(FILE* filePointer, char** string) {
    char* buffer = malloc(INITIAL_SIZE);
    if (!buffer) {
        return 1;
    }

    size_t size = INITIAL_SIZE;
    size_t length = 0;
    int character;

    while ((character = fgetc(filePointer)) != EOF && character != '\n') {
        if (length + 1 >= size) {
            size *= 2;
            char* tempBuffer = realloc(buffer, size);

            if (!tempBuffer) {
                free(buffer);
                return 2;
            }

            buffer = tempBuffer;
        }
        
        char c = (char)character;
        if (excluded(c)) {
            continue;
        }

        buffer[length++] = c;
    }

    if (character == EOF) {
        free(buffer);
        return 3;
    }

    buffer[length] = '\0';

    char* tempBuffer = realloc(buffer, length + 1);
    if (!tempBuffer) {
        free(buffer);
        return 2;
    }
    *string = tempBuffer;
    return 0;
}

int main() {
    FILE *filePointer;
    filePointer = fopen("pokemon", "r");
    int code = 0;
    char* line = NULL;

    while ((code = read_line(filePointer, &line)) == 0) {
        printf("%s", line);
        free(line);
    }

    if (code == 3) {
        printf("The END");
    }

    fclose(filePointer);
}
