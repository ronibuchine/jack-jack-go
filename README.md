# Jack Compiler
A compiler for Jack programming language written in Go... That's it...

## Adding Features
Features will be added by the following workflow:
    1. Checkout a new branch for the feature with the name of the feature as the branch name using `git checkout -b <branchname>`
    2. Continue to work on all feature related items on that branch.
    3. Once completed, submit a PR and select a code reviewer to merge the feature.    
    4. After merged the branch can be deleted.


## Test Script

### usage: 
`./tester.py` or `./tester.py <chapters youd like to test>`.

Can also be run explicitly as a python program like `python tester.py 7`

Currently only handles chapters 7 and 8

Alternitavely one can use the Makefile to run tests or compile an executable with `make` or `make test` for Linux systems
and `mingw32-make` and `mingw32-make test` for Windows. To run tests for specific files use `make test<test-#>` or `mingw32-make test<test-#>`
