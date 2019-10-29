/* jsmin.c
   2013-03-29

Copyright (c) 2002 Douglas Crockford  (www.crockford.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

The Software shall be used for Good, not Evil.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

// @File: jsmin.go
// @Date: 2018/09/10 17:34:30
// https://github.com/douglascrockford/JSMin
*/

package main

import (
	"fmt"
	"errors"
)

const EOF = -1

var (
	theA         int
	theB         int
	theLookahead  = EOF
	theX          = EOF
	theY          = EOF
)

var code []byte
var minjs []byte

func jserror(s string) {
	panic("JSMIN Error: " + s)
}

/* isAlphanum -- return true if the character is a letter, digit, underscore,
   dollar sign, or non-ASCII character.
*/
func isAlphanum(c int) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || c == '_' || c == '$' || c == '\\' || c > 126
}

func getc() int {
	if len(code) <= 0 {
		return EOF
	}
	ret := code[0]
	code = code[1:]

	return int(ret)
}

func putc(c int)  {
	minjs = append(minjs, byte(c))
}

/* get -- return the next character from stdin. Watch out for lookahead. If
   the character is a control character, translate it to a space or
   linefeed.
*/
func get() int {
	c := theLookahead
	theLookahead = EOF
	if c == EOF {
		c = getc()
	}
	if c >= ' ' || c == '\n' || c == EOF {
		return c
	}
	if c == '\r' {
		return '\n'
	}
	return ' '
}

/* peek -- get the next character without getting it.
 */
func peek() int {
	theLookahead = get()
	return theLookahead
}

/* next -- get the next character, excluding comments. peek() is used to see
   if a '/' is followed by a '/' or '*'.
*/
func next() int {
	c := get()
	if c == '/' {
		switch peek() {
		case '/':
			for {
				c = get()
				if c <= '\n' {
					break
				}
			}
			break
		case '*':
			get()
			for c != ' ' {
				switch get() {
				case '*':
					if peek() == '/' {
						get()
						c = ' '
					}
					break
				case EOF:
					jserror("Unterminated comment.")
				}
			}
			break
		}
	}
	theY = theX
	theX = c
	return c
}

/* action -- do something! What you do is determined by the argument:
        1   Output A. Copy B to A. Get the next B.
        2   Copy B to A. Get the next B. (Delete A).
        3   Get the next B. (Delete B).
   action treats a string as a single character. Wow!
   action recognizes a regular expression if it is preceded by ( or , or =.
*/
func action(d int) {
	switch d {
	case 1:
		putc(theA)
		if (theY == '\n' || theY == ' ') && (theA == '+' || theA == '-' || theA == '*' || theA == '/') && (theB == '+' || theB == '-' || theB == '*' || theB == '/') {
			putc(theY)
		}
		fallthrough
	case 2:
		theA = theB
		if theA == '\'' || theA == '"' || theA == '`' {
			for {
				putc(theA)
				theA = get()
				if theA == theB {
					break
				}
				if theA == '\\' {
					putc(theA)
					theA = get()
				}
				if theA == EOF {
					jserror("Unterminated string literal.")
				}
			}
		}
		fallthrough
	case 3:
		theB = next()
		if theB == '/' && (theA == '(' || theA == ',' || theA == '=' || theA == ':' || theA == '[' || theA == '!' || theA == '&' || theA == '|' || theA == '?' || theA == '+' || theA == '-' || theA == '~' || theA == '*' || theA == '/' || theA == '{' || theA == '\n') {
			putc(theA)
			if theA == '/' || theA == '*' {
				putc(' ')
			}
			putc(theB)
			for {
				theA = get()
				if theA == '[' {
					for {
						putc(theA)
						theA = get()
						if theA == ']' {
							break
						}
						if theA == '\\' {
							putc(theA)
							theA = get()
						}
						if theA == EOF {
							jserror("Unterminated set in Regular Expression literal.")
						}
					}
				} else if theA == '/' {
					switch peek() {
					case '/':
						fallthrough
					case '*':
						jserror("Unterminated set in Regular Expression literal.")
					}
					break
				} else if theA == '\\' {
					putc(theA)
					theA = get()
				}
				if theA == EOF {
					jserror("Unterminated Regular Expression literal.")
				}
				putc(theA)
			}
			theB = next()
		}
	}
}

/* jsmin -- Copy the input to the output, deleting the characters which are
   insignificant to JavaScript. Comments will be removed. Tabs will be
   replaced with spaces. Carriage returns will be replaced with linefeeds.
   Most spaces and linefeeds will be removed.
*/
func jsmin() {
	if peek() == 0xEF {
		get()
		get()
		get()
	}
	theA = '\n'
	action(3)
	for theA != EOF {
		switch theA {
		case ' ':
			if isAlphanum(theB) {
				action(1)
			} else {
				action(2)
			}
		case '\n':
			switch theB {
			case '{':
				fallthrough
			case '[':
				fallthrough
			case '(':
				fallthrough
			case '+':
				fallthrough
			case '-':
				fallthrough
			case '!':
				fallthrough
			case '~':
				action(1)
			case ' ':
				action(3)
			default:
				if isAlphanum(theB) {
					action(1)
				} else {
					action(2)
				}
			}
		default:
			switch theB {
			case ' ':
				if isAlphanum(theA) {
					action(1)
				} else {
					action(3)
				}
			case '\n':
				switch theA {
				case '}':
					fallthrough
				case ']':
					fallthrough
				case ')':
					fallthrough
				case '+':
					fallthrough
				case '-':
					fallthrough
				case '"':
					fallthrough
				case '\'':
					fallthrough
				case '`':
					action(1)
				default:
					if isAlphanum(theA) {
						action(1)
					} else {
						action(3)
					}
				}
			default:
				action(1)
			}
		}
	}
}

/* main -- Output any command line arguments as comments
   and then minify the input.
*/
func MinJS(originalJSCode []byte) (minJSCode []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			err = errors.New(e.(string))
		}
	}()

	code = originalJSCode[:]
	minjs = []byte{}
	jsmin()
	return minjs, nil
}
