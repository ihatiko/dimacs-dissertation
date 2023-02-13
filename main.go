package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Clause struct {
	MChVector  []int   // характеристический вектор
	MVariables []int64 // входящие переменные
	// длина дизъюнкта есть m_variables.length()
	MNeg int   // число отрицаний
	MSum int64 // сумма номеров переменных
}

func (c *Clause) GetMChVectorData(index int) int {
	return c.MChVector[index]
}

func (c *Clause) AddDimacsLiteral(v int) bool {
	signum := B2i(0 < v) - B2i(v < 0)
	// индекс вставки можно вычислить как signum(v)*v
	index := signum*v - 1 // нумерация в DIMACS c единицы, вызывающая ф-ция не может перенумеровать, т.к. есть отриц. литеры!.
	pr := signum * c.MChVector[index]
	// контрарную пару можно обнаружить по произведениию знака на значение литеры
	// pr = signum * ch_vector[index] < 0
	if pr == -1 { // контрарная пара
		return false
	}
	if pr == 0 {
		c.MChVector[index] = signum

		c.MVariables = append(c.MVariables, int64(index))
		c.MSum += int64(index) // добавляем сумму номеров пер.
	} // случай pr == 1 - дубликат литеры
	return true
}
func pushBackOff(data []int64, value int64) {

}

type ClauseList struct {
	MList []Clause // массив списков по длинам
}

type Colla struct {
	/*	std::vector<ClauseList> m_lists;  // вектор списков дизъюнктов
		size_t m_vars; // число переменных
		int m_clauses; // число дизъюнктов
		std::vector<int> m_ocs; // счетчики вхождений переменных
		std::vector<int> m_pos;  // счетчики положительных вхождений
		std::vector<int> m_neg;  // счетчики отрицательных вхождений
		bool m_modified = true; // признак модификации (удалена клауза или литера)*/

	MLists    []ClauseList // вектор списков дизъюнктов
	MVars     int          // число переменных
	MClauses  int          // число дизъюнктов
	MOcs      []int        //счетчики вхождений переменных
	MPos      []int        // счетчики положительных вхождений
	MNeg      []int        // счетчики отрицательных вхождений
	MModified bool         // признак модификации (удалена клауза или литера)
}

func (c *Colla) AddClause(clause *Clause) {
	fmt.Println(*clause)
	/*	// обновляем счетчики литер
		size_t index; int value;
		for (size_t i = 0; i < clause.m_variables.size(); ++i) {
			index = clause.m_variables[i]; value = clause[index];
			this->m_ocs[index] += 1; // считаем вхождение
			if (value > 0) // в зависимости от знака очередной литеры
				this->m_pos[index] += 1; // увеличиваем счетчик положительных литер,
			else
			this->m_neg[index] += 1; // либо отрицательных
		}

		size_t length = clause.length();
		if (this->m_lists.size() < length) // при необъодимости досоздаем новые списки
			this->m_lists.resize(length);
		this->m_lists[length - 1].insert(clause); // вставляем клаузу на место
		this->m_clauses += 1; // увеличиваем счетчик дизъюнктов*/
}

func DecompSat(file *os.File, assigment *[]int) {
	// Функция определяет выполнимость формулы
	// Вход: std::fstream &dimacs - файл формата DIMACS
	// в случае ошибки чтения файла выбрасывается соотв. исключение
	// Возвр. значение: целочисленное, обозн. результат решения
	// 0 - КНФ тождественно ложна
	// 1 - КНФ тожд. истинна
	// 2 - КНФ выполнима, вып. набор в &assignment
	// выходной пар.: &assignment - выполняющий набор
	//

	cnf := Colla{MModified: true}
	// 1. построение стр-ры COLLA
	CollaBuilder(file, &cnf)
}
func CollaBuilder(file *os.File, colla *Colla) {
	// построение структры COLLA из файла
	// Вход: файловый поток input, выход: экземпляр класса COLLA, colla
	// возвращает код ошибки: 0 -успешно, 1- ошибка чтения файла
	rd := bufio.NewReader(file)
	var parsedFile []string
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		parsedFile = append(parsedFile, strings.Replace(line, "\n", "", 1))
	}

	if len(parsedFile) < 2 {
		return
	}

	headers := strings.Split(parsedFile[0], " ")
	vars, _ := strconv.Atoi(headers[2])    // считываем число переменных
	clauses, _ := strconv.Atoi(headers[3]) // Считываем число дизъюнктов
	fmt.Println(vars, clauses)

	colla.MVars = vars
	colla.MClauses = clauses
	colla.MOcs = make([]int, vars)
	colla.MNeg = make([]int, vars)
	colla.MPos = make([]int, vars)
	for i := 1; i < len(parsedFile); i++ {
		clause := &Clause{MChVector: make([]int, vars), MNeg: 0, MSum: 0, MVariables: []int64{}}
		if ParseDimacsLine(parsedFile[i], clause) {
			colla.AddClause(clause)
		}
	}
}
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
func ParseDimacsLine(data string, clause *Clause) bool {
	// функция строит хар. вектор по строковому представлению дизъюнкта
	// возвращает true, если дизъюнкт построен: не было контрарных пар.
	// s - входная строка, vars - число литер в формуле
	// c - выходной дизъюнкт
	// Предусловие - дизъюкнт с "пуст"

	splittedData := strings.Split(data, " ")
	offset := 0
	v, _ := strconv.Atoi(splittedData[offset])
	offset += 1
	for {
		if v == 0 {
			break
		}
		sigNum := B2i(0 < v) - B2i(v < 0)
		// индекс вставки можно вычислить как signum(v)*v
		insertIndex := sigNum*v - 1
		pr := sigNum * clause.GetMChVectorData(insertIndex)
		// контрарную пару можно обнаружить по произведениию
		// pr = signum * clause[index]
		if pr == -1 { // контрарная пара
			return false
		}
		if pr == 0 {
			clause.AddDimacsLiteral(v)
		}
		v, _ = strconv.Atoi(splittedData[offset])
		offset += 1
	}
	return true
}

const (
	assetPath = "assets/test-data-set.cnf"
)

func main() {
	file, err := os.Open(assetPath)
	var assigment []int
	if err != nil {
		panic(err)
	}
	DecompSat(file, &assigment)
}
