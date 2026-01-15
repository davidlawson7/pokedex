#ifndef ANAGRAM_H
#define ANAGRAM_H
#include <stdio.h>

struct arguments;

typedef struct arguments {
  char *args[1];
  int species, types;
} arguments;

int excluded(char character);
int read_line(FILE* filePointer, char** string, arguments *arguments);

/*
 * Checks to see if `string1` and `string2` are anagrams of one another. Will
 * return `true`, or `1`, if this is the case, else `false`, or `0`. If `NULL`
 * is passed in either of the string inputs, `false`, or `0` is returned.
 */
int isAnagrams(char *string1, char *string2);

/*
 * Takes a Pokemon and its associated line details. Based on the provided flags
 * it will then convert it to 1 long string to check anagrams against.
 * 
 * An example string would be:
 *
 * `Bulbasaur, Seed, Grass, Poison`
 *
 * Below are some example with this input and what they would produce:
 *
 * ```
 * 
 * char* line = "Bulbasaur, Seed, Grass, Poison";
 *
 * const a = parse_line(&line, 0, 0);
 * const b = parse_line(&line, 1, 0);
 * const c = parse_line(&line, 0, 1);
 * const d = parse_line(&line, 1, 1);
 * 
 * ```
 *
 * yields:
 * 
 * printf("%s\n", line);
 *
 * `bulbasaur`
 * `bulbasaurseed`
 * `bulbasaurgrasspoison`
 * `bulbasaurseedgrasspoison`
 *
 */
int parse_line(char** line, arguments* arguments);

#endif
