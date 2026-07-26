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
	"os"
	"os/exec"
)

func RunLoggedInDir(dir string, exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func RunLoggedInDirMust(dir string, exe string, args ...string) {
	Must(RunLoggedInDir(dir, exe, args...))
}

func RunMust(exe string, args ...string) string {
	cmd := exec.Command(exe, args...)
	d, err := cmd.CombinedOutput()
	out := string(d)
	PanicIf(err != nil, "'%s' failed with '%s', out:\n'%s'\n", cmd.String(), err, out)
	return out
}

func RunLoggedMust(exe string, args ...string) string {
	cmd := exec.Command(exe, args...)
	d, err := cmd.CombinedOutput()
	out := string(d)
	PanicIf(err != nil, "'%s' failed with '%s', out:\n'%s'\n", cmd.String(), err, out)
	fmt.Printf("%s:\n%s\n", cmd.String(), out)
	return out
}
