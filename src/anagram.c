#include "anagram.h"
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

int read_line(FILE* filePointer, char** string, arguments *arguments) {
    char* buffer = malloc(INITIAL_SIZE);
    if (!buffer) {
        return 0;
    }

    // Control variables
    int column = 0; // 0 = name, 1 = species, 2 = type1, 3 = type2
    int shouldSkip = 0; // 0 = dont skip, 1 = skip
    int positionOfName = 0; // Position in buffer where name ends

    size_t size = INITIAL_SIZE;
    size_t length = 0;
    int uChar;

    while ((uChar = fgetc(filePointer)) != EOF && uChar != '\n') {
        if (length + 1 >= size) {
            size *= 2;
            char* tempBuffer = realloc(buffer, size);

            if (!tempBuffer) {
                free(buffer);
                return 0;
            }

            buffer = tempBuffer;
        }

        char character = (char)tolower(uChar);
        if (excluded(character)) {
            if (character == ',') {
                // Time to move the reader along
                ++column;
                
                switch (column) {
                    case 1:
                        positionOfName = length; 
                        shouldSkip = arguments->species == 0; 
                        break;
                    case 2:
                    case 3:
                        shouldSkip = arguments->types == 0;
                        break;
                    default:
                        shouldSkip = 0;
                }
            }

            continue;
        }

        if (shouldSkip) {
            continue;
        }    

        buffer[length++] = character;
    }

    if (uChar == EOF) {
        free(buffer);
        return 0;
    }

    buffer[length] = '\0';

    char* tempBuffer = realloc(buffer, length + 1);
    if (!tempBuffer) {
        free(buffer);
        return 0;
    }
    *string = tempBuffer;
    return positionOfName;
}

int should_skip(int column, arguments *arguments) {
    return 0;
}

int isAnagrams(char *string1, char *string2) {
    if (string1 == NULL || string2 == NULL) {
        return 0;
    }

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
