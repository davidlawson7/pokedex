# `pokegram` - Pokemon Anagram Solver

A fun little program to solve anagrams of Pokemon. Specifically it'll see what
Pokemon are could be an anagram for the provided single word text input, and
list them to stdout.

- [Installation](#installation)
- [About](#about)
- [Examples](#examples)
- [Usage](#usage)
- [Testing](#testing)
- [TODO](#todo)
- [License](#license)

## Installation

## About

Designed to eventually be run a small tiny computer I could put inside a 3d
printed "PokeDex". Something you could look up any Pokemon, and get useful
details when playing Nuzlocks (i.e. ev, stats, move and evolution levels, etc).

This Anagram solver is just a misc function it could do for "fun". Built
initially to be a standalone program I can run via the comandline but later
apart of the "PokeDex".

## Examples

```console
$ anagram -st dosmokiesunplug
muk


$ anagram dosmokiesunplug --types --species
muk
```

## Usage

```console
$ anagram --help
  Usage: anagram [OPTION...] WORD
  Pokemon Anagram Solver -- Solve anagrams using pokemon name, species and type.

    -s, --species              Includes the Pokemons species
    -t, --type                 Includes the Pokemons type (both if exist)
    -?, --help                 Give this help list
        --usage                Give a short usage message
    -V, --version              Print program version

  Report bugs to <david.lawson.95@outlook.com>.
```

## Testing
There are two types of tests. Unit and integration. Both are very simple,
designed to help my sanity a little.

```console
$ make test
```

## TODO

There remains a little to do in this repo.
1. Refactor code to be cleaner.
2. Finish the unit test suite.
3. Write bash based integration test suit testing againsts valid cases.
4. Work out install instructions and verify.

## License

MIT License
