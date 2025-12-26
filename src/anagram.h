#ifndef ANAGRAM_H
#define ANAGRAM_H
#include <stdio.h>

int excluded(char character);
int read_line(FILE* filePointer, char** string);

/*
 * Checks to see if `string1` and `string2` are anagrams of one another. Will
 * return `true`, or `1`, if this is the case, else `false`, or `0`. If `NULL`
 * is passed in either of the string inputs, `false`, or `0` is returned.
 */
int isAnagrams(char *string1, char *string2);

#endif
