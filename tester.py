#!/usr/bin/env python3

import subprocess as sp
import os
import sys

MAIN_FILE_PATH = os.path.join('.', 'cmd', 'compiler', 'compiler.go')

EMULATOR = os.path.join('emulators', 'CPUEmulator' + '.sh' if os.name == 'posix' else '.bat')

def test_outputs(chapters):
    for chapter in [os.path.join('tests', '0' + chap) for chap in chapters]:
        for category in os.listdir(chapter):
            for test in os.listdir(os.path.join(chapter, category)):
                test_dir = os.path.join(chapter, category, test, '')
                sp.run(['go', 'run', MAIN_FILE_PATH,  test_dir])
                os.replace(test + '.asm', test_dir + test + '.asm')
                out = sp.run(['./' + EMULATOR, test_dir + test + '.tst'], capture_output=True)
                print("")
                if out.stdout == b'End of script - Comparison ended successfully\n':
                    print(f"Test {test} ran correctly on the emulator")
                else:
                    print(f"Test {test} had issues")
                    print(out.stderr.decode("utf-8"))


def test_compiler():
    print("testing compiler")
    try:
        sp.run(['go', 'test', '-v', './...'], check=True)
    except:
        print("tests failed")
        exit()
    print("tests succesful")


def main():
    chapters = sys.argv[1:] if len(sys.argv) > 1 else ['7']
    test_compiler()
    test_outputs(chapters)

if __name__ == '__main__':
    main()
