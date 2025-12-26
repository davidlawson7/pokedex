CC = gcc

CFLAGS = -pedantic -Wall -Werror -std=gnu99

# The name of the source files
SOURCES = munit.h munit.c anagram.h anagram.c

anagram:
	$(CC) $(CFLAGS) ./src/anagram.c ./src/main.c -o ./build/anagram

dir:
	mkdir -p build/opts

test: build_test
	./build/test

build_test: test/anagram_test.c 
	$(CC) $(CFLAGS) ./support/test/unit/munit/munit.c ./src/anagram.c ./test/anagram_test.c  -o ./build/test

default: all

all:
	$(CC) $(CFLAGS) $(SOURCES) -o $(EXE)
