package console

import (
	"fmt"
	"strings"
)

func (a *App) Demo() error {
	fmt.Fprintln(a.output, "ДЕМОНСТРАЦИЯ: пустой парк, все виды объектов, отказ, дубликат, ошибки ввода, отчёты.")
	demo := New(a.fleet, strings.NewReader(demoInput), a.output)
	demo.echo = true
	return demo.Run()
}

const demoInput = `3
4
5
1
1
Городской самокат
S-001
6
да
2.5
1
2
Электровелосипед
E-001
5
да
4
1
3
Обычный велосипед
B-001
10
да
2
1
Защитный шлем
H-001
2
2
Зарядная станция
C-001
1.5
1
1
Неисправный самокат
S-002
8
нет
3
2
1
Дубликат номера
S-001
1
3
Неверная простота
B-002
11
2
2
Неверная энергия
C-002
-1
1
9
не команда
3
4
5
0
`
