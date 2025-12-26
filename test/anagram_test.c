#include "../support/test/unit/munit/munit.h"
#include "../src/anagram.h"

MunitResult excluded_false(const MunitParameter params[], void* fixture) {
    const int result = excluded('c');
    munit_assert(result == 0);
    return MUNIT_OK;
}


MunitResult excluded_true(const MunitParameter params[], void* fixture) {
    const int result = excluded('-');
    munit_assert(result == 1);
    return MUNIT_OK;
}

MunitResult is_anagram_true(const MunitParameter params[], void* fixture) {
    const int result = isAnagrams("dosmokiesunplug", "muksludgepoison");
    munit_assert(result == 1);
    return MUNIT_OK;
}

MunitResult is_anagram_true2(const MunitParameter params[], void* fixture) {
    const int result = isAnagrams("softcoollizards", "dracozoltfossil");
    munit_assert(result == 1);
    return MUNIT_OK;
}

MunitResult is_anagram_false(const MunitParameter params[], void* fixture) {
    const int result = isAnagrams("dosmokiesunpluf", "muksludgepoison");
    munit_assert(result == 0);
    return MUNIT_OK;
}

MunitResult is_anagram_null1(const MunitParameter params[], void* fixture) {
    const int result = isAnagrams("dosmokiesunpluf", NULL);
    munit_assert(result == 0);
    return MUNIT_OK;
}

MunitResult is_anagram_null2(const MunitParameter params[], void* fixture) {
    const int result = isAnagrams(NULL, "muksludgepoison");
    munit_assert(result == 0);
    return MUNIT_OK;
}

MunitTest tests[] = {
  {
    "/excluded", /* name */
    excluded_true, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/excluded", /* name */
    excluded_false, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/anagram", /* name */
    is_anagram_true, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/anagram_draco", /* name */
    is_anagram_true2, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/anagram", /* name */
    is_anagram_false, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/anagram_null", /* name */
    is_anagram_null1, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  {
    "/anagram_null", /* name */
    is_anagram_null2, /* test */
    NULL, /* setup */
    NULL, /* tear_down */
    MUNIT_TEST_OPTION_NONE, /* options */
    NULL /* parameters */
  },
  /* Mark the end of the array with an entry where the test
   * function is NULL */
  { NULL, NULL, NULL, NULL, MUNIT_TEST_OPTION_NONE, NULL }
};

static const MunitSuite suite = {
  "/pokedex", /* name */
  tests, /* tests */
  NULL, /* suites */
  1, /* iterations */
  MUNIT_SUITE_OPTION_NONE /* options */
};

int main(int argc, char* argv[MUNIT_ARRAY_PARAM(argc + 1)]) {
    return munit_suite_main(&suite, NULL, argc, argv);
}
