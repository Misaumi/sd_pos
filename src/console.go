package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Console struct {
	f   *Facade
	in  *bufio.Scanner
	out io.Writer
	eof bool
}

func NewConsole(f *Facade, in io.Reader, out io.Writer) *Console {
	return &Console{f: f, in: bufio.NewScanner(in), out: out}
}

func (c *Console) ask(prompt string) string {
	fmt.Fprint(c.out, prompt)
	if !c.in.Scan() {
		c.eof = true
		return ""
	}
	return strings.TrimSpace(c.in.Text())
}

func (c *Console) askFloat(prompt string) (float64, bool) {
	s := strings.Replace(c.ask(prompt), ",", ".", 1)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintln(c.out, "Ошибка: нужно число")
		return 0, false
	}
	return v, true
}

func (c *Console) printErr(err error) {
	if err != nil {
		fmt.Fprintln(c.out, "Ошибка:", err)
	}
}

func (c *Console) menu() {
	fmt.Fprintln(c.out, `
=== ПОС: Преподавательская Оценка Студентов ===
1  Добавить студента
2  Список студентов
3  Добавить критерий
4  Список критериев
5  Выставить оценку
6  Показать оценки студента
7  Исправить оценку
8  Отменить оценку
9  Итог студента
10 Экспорт итогов в CSV
11 Сменить правило итога (all / rated)
0  Выход`)
}

func (c *Console) Run() {
	for !c.eof {
		c.menu()
		switch c.ask("> ") {
		case "1":
			st, err := c.f.AddStudent(c.ask("Имя: "))
			c.printErr(err)
			if err == nil {
				fmt.Fprintf(c.out, "Добавлен %s (%s)\n", st.Name(), st.ID())
			}
		case "2":
			c.printStudents()
		case "3":
			name := c.ask("Название: ")
			max, ok1 := c.askFloat("max_score: ")
			w, ok2 := c.askFloat("Вес: ")
			if ok1 && ok2 {
				cr, err := c.f.AddCriterion(name, max, w)
				c.printErr(err)
				if err == nil {
					fmt.Fprintf(c.out, "Добавлен критерий %s (%s)\n", cr.Name(), cr.ID())
				}
			}
		case "4":
			c.printCriteria()
		case "5":
			sid := c.ask("id студента: ")
			cid := c.ask("id критерия: ")
			v, ok := c.askFloat("Балл: ")
			if ok {
				comment := c.ask("Комментарий (можно пусто): ")
				c.printErr(c.f.SetGrade(sid, cid, v, comment))
			}
		case "6":
			c.printGrades(c.ask("id студента: "))
		case "7":
			sid := c.ask("id студента: ")
			cid := c.ask("id критерия: ")
			if v, ok := c.askFloat("Новый балл: "); ok {
				c.printErr(c.f.ChangeGrade(sid, cid, v))
			}
		case "8":
			sid := c.ask("id студента: ")
			cid := c.ask("id критерия: ")
			c.printErr(c.f.CancelGrade(sid, cid))
		case "9":
			t, err := c.f.Total(c.ask("id студента: "))
			c.printErr(err)
			if err == nil {
				fmt.Fprintf(c.out, "Итог: %.2f%%\n", t)
			}
		case "10":
			path := c.ask("Файл (Enter = totals.csv): ")
			if path == "" {
				path = "totals.csv"
			}
			err := c.f.ExportTotals(path)
			c.printErr(err)
			if err == nil {
				fmt.Fprintln(c.out, "Сохранено в", path)
			}
		case "11":
			c.printErr(c.f.SetStrategy(c.ask("all или rated: ")))
		case "0":
			return
		}
	}
}

func (c *Console) printStudents() {
	for _, s := range c.f.ListStudents() {
		fmt.Fprintf(c.out, "%s  %s\n", s.ID(), s.Name())
	}
}

func (c *Console) printCriteria() {
	for _, cr := range c.f.ListCriteria() {
		fmt.Fprintf(c.out, "%s  %s  max=%.1f  вес=%.1f\n", cr.ID(), cr.Name(), cr.MaxScore(), cr.Weight())
	}
}

func (c *Console) printGrades(studentID string) {
	list, err := c.f.ListGrades(studentID)
	if err != nil {
		c.printErr(err)
		return
	}
	for _, e := range list {
		fmt.Fprintf(c.out, "%s: %.1f  (%s) %s\n", e.CriterionID(), e.Score().Value(),
			e.Date().Format("02.01.2006"), e.Comment())
	}
}

// RunDemo - сквозной сценарий из ТЗ
func RunDemo(f *Facade, out io.Writer) {
	fmt.Fprintln(out, "--- Демо ---")
	a, _ := f.AddStudent("Иван Петров")
	b, _ := f.AddStudent("Мария Сидорова")
	lec, _ := f.AddCriterion("Лекции", 10, 1)
	sem, _ := f.AddCriterion("Семинары", 10, 3)
	hw, _ := f.AddCriterion("Домашки", 20, 4)

	f.SetGrade(a.ID(), lec.ID(), 8, "")
	f.SetGrade(a.ID(), sem.ID(), 9, "активен")
	f.SetGrade(a.ID(), hw.ID(), 15, "")
	f.SetGrade(b.ID(), lec.ID(), 10, "")
	f.SetGrade(b.ID(), hw.ID(), 18, "")

	fmt.Fprintln(out, "Попытка поставить 12 из 10:")
	fmt.Fprintln(out, "  ", f.SetGrade(b.ID(), sem.ID(), 12, ""))
	fmt.Fprintln(out, "Попытка оценить критерий c99:")
	fmt.Fprintln(out, "  ", f.SetGrade(b.ID(), "c99", 5, ""))

	for _, s := range f.ListStudents() {
		t, _ := f.Total(s.ID())
		fmt.Fprintf(out, "Итог %s: %.2f%%\n", s.Name(), t)
	}
	if err := f.ExportTotals("totals.csv"); err != nil {
		fmt.Fprintln(out, "Ошибка экспорта:", err)
		return
	}
	fmt.Fprintln(out, "Итоги выгружены в totals.csv")
}
