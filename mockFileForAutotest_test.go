package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	test()
	// без этого файла не выполняется автотест 2 инкремента
	//Messages:   	Невозможно получить результат выполнения команды: /usr/local/go/bin/go test -cover ./... Вывод:
	//go: warning: "./..." matched no packages
	// go test -cover ./... не работает с модулями, где есть файл go.mod
	// команда go test -cover работает только при вызове внутри модуля
}
