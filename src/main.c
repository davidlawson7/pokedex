#include "anagram.h"
#include <argp.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const char *argp_program_version = "poke-anagram 1.0";
const char *argp_program_bug_address = "<david.lawson.95@outlook.com>";

static char doc[] = "Pokemon Anagram Solver -- Solve anagrams using pokemon "
                    "name, species and type.";

static char args_doc[] = "WORD";

static struct argp_option options[] = {
    {"species", 's', 0, 0, "Includes the Pokemons species"},
    {"type", 't', 0, 0, "Includes the Pokemons type (both if exist)"},
    {0}};

static error_t parseOpt(int key, char *arg, struct argp_state *state) {
  struct arguments *arguments = state->input;

  switch (key) {
  case 's':
    arguments->species = 1;
    break;
  case 't':
    arguments->types = 1;
    break;

  case ARGP_KEY_ARG:
    if (state->arg_num >= 1) {
      /* Too many arguments. */
      argp_usage(state);
    }
    arguments->args[state->arg_num] = arg;
    break;
  case ARGP_KEY_END:
    if (state->arg_num < 1) {
      /* Not enough arguments. */
      argp_usage(state);
      break;
    }
  default:
    return ARGP_ERR_UNKNOWN;
  }
  return 0;
}

static struct argp argp = {options, parseOpt, args_doc, doc};

int main(int argc, char *argv[]) {
  struct arguments arguments;
  arguments.species = 0;
  arguments.types = 0;

  argp_parse(&argp, argc, argv, 0, 0, &arguments);

  FILE *filePointer;
  filePointer = fopen("data/pokemon", "r");
  int code = 0;
  char *line = NULL;

  while ((code = read_line(filePointer, &line, &arguments)) != 0) {
    if (isAnagrams(arguments.args[0], line)) {
      for (int i= 0; i < code; ++i) {
        printf("%c", line[i]);
      }
      printf("%c\n", '\0');
    }
    free(line);
  }

  fclose(filePointer);
  exit(0);
}
