package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"go/ast"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/ast/inspector"
	"strings"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

// Это пакет со статическими анализаторами
// Для запуска используйте команду main [-flag] [package]
// или main help для вызова справки.
// Включает в себя следующие анализаторы:
// 1. errcheck.Analyzer - анализирует необработанные ошибки
// 2. inspect.Analyzer - анализирует абстрактное синтаксическое дерево (AST) кода
// 3. printf.Analyzer - анализирует использование функций форматирования printf
// 4. shadow.Analyzer - анализирует переопределение переменных внутри блоков
// 5. shift.Analyzer - анализирует использование сдвига влево/вправо на недопустимое количество битов
// 6. structtag.Analyzer - анализирует использование некорректных тегов структур
// Кроме того, в коде также определен пользовательский анализатор exitcall,
// который анализирует использование функции os.Exit в функции main пакета main.
func main() {
	var mychecks []*analysis.Analyzer
	// Добавление анализаторов класса SA пакета staticcheck.io в mychecks
	for _, v := range staticcheck.Analyzers {
		//всех анализаторов класса SA пакета staticcheck.io;
		if strings.Contains(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
		// Добавление анализаторов других классов пакета staticcheck.io в mychecks
		if v.Analyzer.Name == "ST1001" || v.Analyzer.Name == "QF1007" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	var exitAnalizer = &analysis.Analyzer{
		Name: "exitcall",
		Doc:  "Анализирует использование os.Exit в функции main пакета main",
		Run:  run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	mychecks = append(
		mychecks,
		exitAnalizer,
		// стандартных статических анализаторов пакета
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		errcheck.Analyzer,
	)

	multichecker.Main(
		mychecks...,
	)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspct := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspct.Preorder(nodeFilter, func(node ast.Node) {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}
		fun, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		ident, ok := fun.X.(*ast.Ident)
		if !ok || ident.Name != "os" || fun.Sel.Name != "Exit" {
			return
		}
		pass.Reportf(callExpr.Pos(), "прямой вызов os.Exit в функции main пакета main запрещен")
	})

	return nil, nil
}
