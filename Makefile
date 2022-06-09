
MAIN=compiler.go

default:
	go build -o jjg $(MAIN)

test:
	python tester.py

test7:
	python tester.py 7

test8:
	python tester.py 8

test9:
	python tester.py 9

test10:
	python tester.py 10

test11:
	python tester.py 11


