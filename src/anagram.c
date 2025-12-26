#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

#define INITIAL_SIZE 128

int excluded(char character) {
    int ignoreList[] = {',', ' ', '-'};
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

        char c = (char)tolower(character);
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

int isAnagrams(char *string1, char *string2) {
    if (strlen(string1) != strlen(string2)) {
        return 0;
    }

    int freq[26] = {0};
    
    for (int i = 0; string1[i] != '\0'; i++) {
        freq[string1[i] - 'a']++;
    }

    for (int i = 0; string2[i] != '\0'; i++) {
        freq[string2[i] - 'a']--;
    }

    for (int i = 0; i < 26; i++) {
        if (freq[i] != 0) {
            return 0;
        }
    }

    return 1;
}
