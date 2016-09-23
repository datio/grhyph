### A Syllabification and Hyphenation Library for Modern Greek
The syllabification of Greek text, especially when written informally, can be problematic.
Written in the Greek alphabet, words may miss accents or diaeresis diacritics, reducing the result correctness in algorithms that follow a simplified Modern Greek grammar ruleset.
Greek content written using Latin characters, also known as Greeklish, is a common occurrence online. Hyphenating Greeklish words has its own set of challenges, such as how some character sequences map from one alphabet to the other.
Some words may include syllabic vowels that should not be separated on hyphenation, due to a phenomenon known as synizesis, but instead should be combined into a single syllable.
This repository contains the implementation a Modern Greek hyphenation library that provides support for exceptions using regular expressions and a test CLI program (*[grhyph_cli](https://github.com/datio/grhyph/tree/master/grhyph_cli)*).

[![grhyph_cli](https://i.imgur.com/8klAJt5.png)](https://asciinema.org/a/epf5dnx24w7uwm09aonol2kdl)