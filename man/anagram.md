Pokemon Anagram Solver "2025" "General Commands Manual"
=======================================================

NAME
----

`pokegram` - Pokémon Anagram Solver.

SYNOPSIS
--------

`pokegram [OPTIONS] <anagram>`

`pokegram [-s | -t] <anagram>`

DESCRIPTION
-----------

Pokémon Anagram Solver.

`pokegram` is able to solve anagrams of different Pokémon names in combination
with their type, species, and other traits. By default, `pokegram` will only
use a Pokémons name to check against. If additional traits are in the anagram,
you can use the flags defined in this manual to include them in the check.

OPTIONS
-------

`-s`, `--species`
  Tells `pokegram` to include the Pokémons species.

`-t`, `--types`
  Tells `pokegram` to include the Pokémons types.
