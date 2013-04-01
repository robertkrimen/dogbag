package kilt

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Kilt struct {
}

func New() *Kilt {
	return &Kilt{}
}

func (self Kilt) Sha1Path(path ...string) string {
	return Sha1Path(path...)
}

func Sha1Path(path ...string) string {
	file, err := os.Open(filepath.Join(path...))
	if err != nil {
		return ""
	}
	return Sha1Of(file)
}

func (self Kilt) Sha1Of(src io.Reader) string {
	return Sha1Of(src)
}

func Sha1Of(src io.Reader) string {
	hash := sha1.New()
	_, err := io.Copy(hash, src)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (self Kilt) Sha1(data []byte) string {
	return Sha1(data)
}

func Sha1(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func (self Kilt) GraveTrim(target string) string {
	return GraveTrim(target)
}

func GraveTrim(target string) string {
	// Discard \r? Go already does this for raw string literals.
	end := len(target)

	last := 0
	index := 0
	for index = 0; index < end; index++ {
		chr := rune(target[index])
		if chr == '\n' || !unicode.IsSpace(chr) {
			last = index
			break
		}
	}
	if index >= end {
		return ""
	}
	start := last
	if rune(target[start]) == '\n' {
		// Skip the leading newline
		start++
	}

	last = end - 1
	for index = last; index > start; index-- {
		chr := rune(target[index])
		if chr == '\n' || !unicode.IsSpace(chr) {
			last = index
			break
		}
	}
	stop := last
	result := target[start : stop+1]
	return result
}

func Symlink(oldname, newname string, overwrite bool) error {
	err := os.Symlink(oldname, newname)
	if err == nil {
		return nil // Success
	}
	if !os.IsExist(err) {
		return err // Failure
	}
	// Failure, file exists
	symbolic := false
	{
		stat, err := os.Lstat(newname)
		if err != nil {
			return err
		}
		symbolic = stat.Mode()&os.ModeSymlink != 0
	}
	if !symbolic {
		return err
	}
	if !overwrite {
		return nil
	}
	err = os.Remove(newname)
	if err != nil {
		return err
	}
	return os.Symlink(oldname, newname)
}

func (self Kilt) Symlink(oldname, newname string, overwrite bool) error {
	return Symlink(oldname, newname, overwrite)
}

type QuoteWord struct {
	Value string
	Quote string
	Space [2]string
	Index [2]int
}

type _quoteScan struct {
	input  string
	length int
	index  int
	last   int
	done   bool
}

func (self *_quoteScan) next() rune {
	if self.index >= self.length {
		self.done = true
		self.last = -1
		return -1
	}
	next, width := utf8.DecodeRuneInString(self.input[self.index:])
	self.last = width
	self.index += width
	return next
}

func (self *_quoteScan) backup() {
	if self.last != -1 {
		self.index -= self.last
		self.last = -1
	}
}

func (self Kilt) QuoteParse(input string) []QuoteWord {
	return QuoteParse(input)
}

func QuoteParse(input string) []QuoteWord {
	result := []QuoteWord{}
	scan := _quoteScan{
		input:  input,
		length: len(input),
		index:  0,
	}

	for !scan.done {
		item := QuoteWord{}
		item.Index[0] = scan.index

		text := bytes.Buffer{}
	Space:
		for !scan.done {
			next := scan.next()
			switch {
			case unicode.IsSpace(next):
				text.WriteRune(next)
			default:
				scan.backup()
				break Space
			}
		}

		item.Space[0] = text.String()
		if length := len(result); length > 0 {
			result[length-1].Space[1] = item.Space[0]
		}

		if !scan.done {
			text.Reset()

			delimit := 0
			delimeter := rune(-1)
			next := scan.next()
			if strings.ContainsRune("\"'", next) {
				delimit += 1
				delimeter = next
			} else {
				scan.backup()
			}

			if !scan.done {
			Value:
				for {
					next := scan.next()
					switch next {
					case -1:
						break Value
					case '\\':
						next = scan.next()
						text.WriteRune(next)
					default:
						if delimit == 0 {
							if unicode.IsSpace(next) {
								scan.backup()
								break Value
							} else if strings.ContainsRune("\"'", next) {
								scan.backup()
								break Value
							}
						} else {
							if next == delimeter {
								delimit += 1
								break Value
							}
						}
						text.WriteRune(next)
					}
				}

				item.Value = text.String()
				if delimit != 0 {
					text.Reset()
					text.WriteRune(delimeter)
					text.WriteString(item.Value)
					if delimit == 2 {
						text.WriteRune(delimeter)
					}
					item.Quote = text.String()
				}
			}
		}

		//    case '"', '\'':
		//    case '\\':
		//        value = self.next()
		//        if isLineTerminator(value) {
		//            if quote == '/' {
		//                return errorIllegal()
		//            }
		//            self.scanEndOfLine(value, false)
		//            continue
		//        }
		//        if quote == '/' { // RegularExpression
		//            // TODO Handle the case of [\]?
		//            text.WriteRune('\\')
		//            text.WriteRune(value)
		//            continue
		//        }
		//        switch value {
		//        case 'n':
		//            text.WriteRune('\n')
		//        case 'r':
		//            text.WriteRune('\r')
		//        case 't':
		//            text.WriteRune('\t')
		//        case 'b':
		//            text.WriteRune('\b')
		//        case 'f':
		//            text.WriteRune('\f')
		//        case 'v':
		//            text.WriteRune('\v')
		//        case '0':
		//            text.WriteRune(0)
		//        case 'u':
		//            result := self.scanHexadecimalRune(4)
		//            if result != utf8.RuneError {
		//                text.WriteRune(result)
		//            } else {
		//                return errorIllegal()
		//            }

		//        case 'x':
		//            result := self.scanHexadecimalRune(2)
		//            if result != utf8.RuneError {
		//                text.WriteRune(result)
		//            } else {
		//                return errorIllegal()
		//            }
		//        default:
		//            text.WriteRune(value)
		//        }
		//        // TODO Octal escaping
		//}

		if scan.index == item.Index[0] {
			break
		}
		item.Index[1] = scan.index
		result = append(result, item)
	}
	return result
}

func (self Kilt) WriteAtomicFile(filename string, data io.Reader, mode os.FileMode) error {
	return WriteAtomicFile(filename, data, mode)
}

func WriteAtomicFile(filename string, data io.Reader, mode os.FileMode) error {
	parent := filepath.Dir(filename)
	tmp, err := ioutil.TempDir(parent, ".tmp.")
	if err != nil {
		return err
	}
	if len(parent) >= len(tmp) { // Should never, ever happen
		panic(fmt.Sprintf("%s < %s", tmp, parent))
	}
	defer os.RemoveAll(tmp)

	tmpname := filepath.Join(tmp, filepath.Base(filename))
	file, err := os.OpenFile(tmpname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	return os.Rename(tmpname, filename)
}
