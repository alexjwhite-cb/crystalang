# Jet Code Style Guide

* Variables should be `camelCased`
* Non-member functions should also be `camelCased`
* Objects, Descriptors, their Properties, and member functions should be `PascalCased`
* Opening braces `{` should be on the same line as the method declaration
* Closing braces `}` should be the first line item at the end of a code block
  * The only exception to the above is cases in which function code can reasonably and readably take place on a single line:
```
meth isEmpty: x { ( len(x) < 1 ? true :: false )-> }
```
* Jet is not whitespace sensitive. Indenting within loops, function calls that continue onto new lines, etc is encouraged, but not currently enforced.
* Semi-colons `;` and new line characters `\n` and `\r` are considered equivalent in terms of line endings. Ergo, `;` should only be used to denote multiple statements on the same line, else they can be omitted entirely.