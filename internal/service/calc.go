package service

import (
	"strconv"
	"unicode"
	"yaly-1/pkg/errs"
)

// CalcService - интерфейс сервисного уровня калькулятора
type CalcService struct{}

// Конструктор для структуры сервисного уровня калькулятора
func NewCalcService() *CalcService {
	return &CalcService{}
}

// Структура для хранения токенов
type tokens [](*token)

// Структура токена
type token struct {
	ttype string
	value string
	pos   int
}

// Вычисление выражения
func (c *CalcService) Calc(ex string) (float64, error) {
	tokensSlice, err := getTokens(ex)
	if err != nil {
		return 0, err
	}
	opStack := tokens{}
	numStack := []float64{}
	for _, tokenChar := range tokensSlice {
		if tokenChar.ttype == "int" || tokenChar.ttype == "float" {
			val, err := strconv.ParseFloat(tokenChar.value, 64)
			if err != nil {
				return 0, err
			}
			numStack = append(numStack, val)
			continue
		}
		if tokenChar.ttype == "op" {
			if tokenChar.value == ")" {
				closed := false
				for len(opStack) != 0 {
					op := opStack[len(opStack)-1]
					opStack = opStack[:len(opStack)-1]
					if op.value == "(" {
						closed = true
						break
					}

					if len(numStack) < 2 {
						return 0, errs.ErrExpressionNotValid
					}
					lnum := numStack[len(numStack)-1]
					numStack = numStack[:len(numStack)-1]
					rnum := numStack[len(numStack)-1]
					numStack = numStack[:len(numStack)-1]
					res, err := execute(op, lnum, rnum)
					if err != nil {
						return 0, errs.ErrExpressionNotValid
					}
					numStack = append(numStack, res)
				}
				if !closed {
					return 0, errs.ErrExpressionNotValid
				}
				continue
			}
			for {
				priority := getPriority(tokenChar)
				backPriority := []int{0, 0}
				if len(opStack) > 0 {
					lastOp := opStack[len(opStack)-1]
					backPriority = getPriority(lastOp)
				}
				if len(opStack) == 0 || tokenChar.value == "(" ||
					priority[1] > backPriority[1] ||
					(priority[1] == backPriority[1] && priority[0] == 0) {
					opStack = append(opStack, tokenChar)
					break
				}

				if len(numStack) < 2 {
					return 0, errs.ErrExpressionNotValid
				}
				op := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				lnum := numStack[len(numStack)-1]
				numStack = numStack[:len(numStack)-1]
				rnum := numStack[len(numStack)-1]
				numStack = numStack[:len(numStack)-1]
				res, err := execute(op, lnum, rnum)
				if err != nil {
					return 0, errs.ErrExpressionNotValid
				}
				numStack = append(numStack, res)
			}
		}
	}

	for len(opStack) != 0 {
		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		if op.value == "(" {
			return 0, errs.ErrExpressionNotValid
		}
		if len(numStack) < 2 {
			return 0, errs.ErrExpressionNotValid
		}

		lnum := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]
		rnum := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]
		res, err := execute(op, lnum, rnum)
		if err != nil {
			return 0, errs.ErrExpressionNotValid
		}
		numStack = append(numStack, res)
	}

	return numStack[0], nil
}

// Получение слайса токенов
func getTokens(ex string) (tokens, error) {
	tokensSlice := tokens{}
	offset := 0
	for len(ex) != 0 {
		chr := getNextToken(&ex)
		offset++
		if chr == "/" || chr == "*" || chr == "+" || chr == "-" || chr == "(" || chr == ")" {
			tokensSlice = append(tokensSlice, &token{"op", chr, offset})
			continue
		}
		if chr == " " || chr == "\t" || chr == "\r" {
			continue
		}
		if unicode.IsDigit(rune(chr[0])) {
			isFloat := false
			val := chr

			for len(ex) != 0 {
				chr := string(ex[0])
				if unicode.IsDigit(rune(chr[0])) {
					ex = ex[1:]
					val += chr
					continue
				}
				if chr == "." && !isFloat {
					ex = ex[1:]
					val += chr
					isFloat = true
					continue
				}
				break
			}

			tokensSlice = append(tokensSlice, &token{"float", val, offset})

			continue
		}
		return nil, errs.ErrExpressionNotValid
	}
	if len(tokensSlice) < 3 {
		return nil, errs.ErrExpressionNotValid
	}
	return tokensSlice, nil
}

// Получение следующего токена
func getNextToken(ex *string) string {
	result := (*ex)[0]
	*ex = (*ex)[1:]
	return string(result)
}

// Вычисление выражения
func execute(op *token, lnum, rnum float64) (float64, error) {
	switch op.value {
	case "+":
		return rnum + lnum, nil

	case "-":
		return rnum - lnum, nil

	case "*":
		return rnum * lnum, nil

	case "/":
		if lnum == 0 {
			return 0, errs.ErrDivisionByZero
		}
		return rnum / lnum, nil
	}
	return 0, nil
}

// Получение приоритета операции
func getPriority(op *token) []int {
	value := op.value
	if value == "/" || value == "*" {
		return []int{1, 2}
	}
	if value == "+" || value == "-" {
		return []int{1, 1}
	}
	return []int{0, 0}
}
