# cmdbox xlua

an extended lua runtime based on golang lua 5.5

## usage

```
$ cmdbox help xlua 
Synopsis:
        extended lua environment
Usage:
        cmdbox xlua [flags] [lua-script]

Available flags:
  -D string
        data-file to load into '_data'
  -e string
        statement to execute
  -extlib
        register extlib (default true)
  -f string
        script-file to execute
  -i    interactive (repl) mode
  -v    print version and exit
```

## base variables

* `arg` – command line arguments

## extensions (extlib)

### base functions

* `apairs` – complementary function to pairs and ipairs
* `stringify(a[,l]) -> string` – stringifies lua structures up to level depth

### utf8 functions

char/rune (not byte) wise

* `utf8.charat (s, n) -> string` 
* `utf8.length (s) -> int`



