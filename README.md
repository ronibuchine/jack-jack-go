# Jack Compiler
A compiler for Jack programming language written in Go... That's it...

## Adding Features
Features will be added by the following workflow:
    1. Checkout a new branch for the feature with the name of the feature as the branch name using `git checkout -b <branchname>`
    2. Continue to work on all feature related items on that branch.
    3. Once completed, submit a PR and select a code reviewer to merge the feature.    
    4. After merged the branch can be deleted.


## Test Script

### Compiling
Compile with `make` on linux. On windows?? -- calm your farm... you can use mingw (command is `mingw32-make`)

### Usage 

`jjg [-tokens] [-ast] <files>`

`-tokens` will output the xml of the tokenized jack file

`-ast` will output the xml of the ast

`files` can be any number of Jack files that are in your current working directory or can be a directoy with Jack files.

### Testing

`./tester.py` or `./tester.py <chapters youd like to test>`.

Can also be run explicitly as a python program like `python tester.py 7`

Currently only handles chapters 7 and 8

Alternatively one can use the Makefile to run tests or compile an executable with `make` or `make test` for Linux systems
and `mingw32-make` and `mingw32-make test` for Windows. To run tests for specific files use `make test<test-#>` or `mingw32-make test<test-#>`

## Flabby Bird

To play you can open up the VMEmulator program, press load program and select the flabby-bird directory. 

Enjoy!
