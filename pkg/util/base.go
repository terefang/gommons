// originally taken from https://github.com/kjk/common/ under:
//MIT License
//
//Copyright (c) 2021 Krzysztof Kowalczyk
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

func PanicIf(cond bool, args ...any) {
	if !cond {
		return
	}
	s := "condition failed"
	if len(args) > 0 {
		s = fmt.Sprintf("%s", args[0])
		if len(args) > 1 {
			s = fmt.Sprintf(s, args[1:]...)
		}
	}
	panic(s)
}

func PanicIfErr(err error, args ...any) {
	if err == nil {
		return
	}
	s := err.Error()
	if len(args) > 0 {
		s = fmt.Sprintf("%s", args[0])
		if len(args) > 1 {
			s = fmt.Sprintf(s, args[1:]...)
		}
	}
	panic(s)
}

func GetErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func IsWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func IsMac() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}

func IsWinOrMac() bool {
	return IsWindows() || IsMac()
}

func IsLinux() bool {
	return strings.Contains(runtime.GOOS, "linux")
}

func GetCallstackFrames(skipFrames int) []string {
	var callers [32]uintptr
	n := runtime.Callers(skipFrames+1, callers[:])
	frames := runtime.CallersFrames(callers[:n])
	var cs []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		s := frame.File + ":" + strconv.Itoa(frame.Line)
		cs = append(cs, s)
	}
	return cs
}

func GetCallstack(skipFrames int) string {
	frames := GetCallstackFrames(skipFrames + 1)
	return strings.Join(frames, "\n")
}

func Push[S ~[]E, E any](s *S, els ...E) {
	*s = append(*s, els...)
}

func SliceLimit[S ~[]E, E any](s S, max int) S {
	if len(s) > max {
		return s[:max]
	}
	return s
}
