# `pokegram` - Pokemon Anagram Solver

A fun little program to solve anagrams of Pokemon. Specifically it'll see what
Pokemon are could be an anagram for the provided single word text input, and
list them to stdout.

- [Installation](#installation)
- [About](#about)
- [Examples](#examples)

## Installation

## About

Designed to eventually be run a small tiny computer I could put inside a 3d
printed "PokeDex". Something you could look up any Pokemon, and get useful
details when playing Nuzlocks (i.e. ev, stats, move and evolution levels, etc).

This Anagram solver is just a misc function it could do for "fun". Built
initially to be a standalone program I can run via the comandline but later
apart of the "PokeDex".

## Examples

```shell
$ anagram -st dosmokiesunplug
muk


$ anagram dosmokiesunplug --types --species
muk
```
