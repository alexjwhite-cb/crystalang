# Jet Code Style Guide

* Variables should be `camelCased`
* Non-member functions should also be `camelCased`
* Objects, Descriptors, their Properties, and member functions should be `PascalCased`
* Opening braces `{` should be on the same line as the method declaration
* Closing braces `}` should on a new line at the end of a code block
```
meth isEmpty: x { ( len(x) < 1 ? true :: false )-> }
```
* Jet is not whitespace sensitive. Indenting within loops, function calls that continue onto new lines, etc. is encouraged, but not currently enforced.
* Semicolons `;` and new line characters `\n` and `\r` are considered equivalent in terms of line endings. Ergo, `;` should only be used to denote multiple statements on the same line, else they can be omitted entirely.