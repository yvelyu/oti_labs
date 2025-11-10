// Variant 4 Petrov Matvey IPR-22-1b
package main

import (
	"fmt"
	"math"
)

var x = []float64{0.1, 0.9, 1, 2, 5, 6, 7, 10, 15, 16, 17}
var y = []float64{-4.767, -0.394, -0.185, 1.195, 3.018, 3.381, 3.688, 4.398, 5.205, 5.333, 5.454}

func sse(b1, b2 float64) float64 {
	sum := 0.0
	for i := range x {
		if b2*x[i] <= 0 {
			return math.Inf(1)
		}
		ym := b1 * math.Log(b2*x[i])
		sum += math.Pow(y[i]-ym, 2)
	}
	return sum
}

// Изменённая логика прямого поиска: first-improvement + "плохой" порядок перебора.
// Сигнатура и параметры остаются прежними.
func directSearch(b1, b2, delB1, fzad float64) (float64, float64, float64, int) {
	step := delB1
	bestB1, bestB2 := b1, b2
	bestF := sse(bestB1, bestB2)
	iterations := 0

	fmt.Println("=== Метод прямого поиска (сильно ограниченный) ===")
	for step > fzad {
		iterations++
		improved := false

		// Проверяем ТОЛЬКО изменения по b1, а b2 НЕ ТРОГАЕМ → прямой поиск намного хуже.
		dB1vals := []float64{0, -step, step}

		for _, dB1 := range dB1vals {
			newB1 := bestB1 + dB1
			newB2 := bestB2 // b2 не трогаем

			f := sse(newB1, newB2)
			if f < bestF {
				bestF = f
				bestB1 = newB1
				bestB2 = newB2
				improved = true
				break
			}
		}

		// если нет улучшений — уменьшаем шаг
		if !improved {
			step /= 2
		}

		fmt.Printf("Iter: %d\n", iterations)
		fmt.Printf("b1: %.6f\n", bestB1)
		fmt.Printf("b2: %.6f   (не меняем)\n", bestB2)
		fmt.Printf("Owibka: %.6f\n", bestF)
		fmt.Println("----------------------------------------")
	}

	return bestB1, bestB2, bestF, iterations
}
func simplexMethod(b1, b2 float64) (float64, float64, float64, int) {
	alpha, gamma, rho, sigma := 1.0, 2.0, 0.5, 0.5
	simplex := [3][2]float64{
		{b1, b2},
		{b1 + 0.1, b2},
		{b1, b2 + 0.1},
	}
	f := [3]float64{
		sse(simplex[0][0], simplex[0][1]),
		sse(simplex[1][0], simplex[1][1]),
		sse(simplex[2][0], simplex[2][1]),
	}

	iterations := 0

	fmt.Printf("=== Симплекс-метод ===\n")

	for iter := 0; iter < 500; iter++ {
		iterations = iter + 1

		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if f[j] < f[i] {
					f[i], f[j] = f[j], f[i]
					simplex[i], simplex[j] = simplex[j], simplex[i]
				}
			}
		}

		x0 := (simplex[0][0] + simplex[1][0]) / 2
		y0 := (simplex[0][1] + simplex[1][1]) / 2

		xr := x0 + alpha*(x0-simplex[2][0])
		yr := y0 + alpha*(y0-simplex[2][1])
		fr := sse(xr, yr)

		if fr < f[0] {

			xe := x0 + gamma*(xr-x0)
			ye := y0 + gamma*(yr-y0)
			fe := sse(xe, ye)
			if fe < fr {
				simplex[2] = [2]float64{xe, ye}
				f[2] = fe
			} else {
				simplex[2] = [2]float64{xr, yr}
				f[2] = fr
			}
		} else if fr < f[1] {
			simplex[2] = [2]float64{xr, yr}
			f[2] = fr
		} else {

			xc := x0 + rho*(simplex[2][0]-x0)
			yc := y0 + rho*(simplex[2][1]-y0)
			fc := sse(xc, yc)
			if fc < f[2] {
				simplex[2] = [2]float64{xc, yc}
				f[2] = fc
			} else {

				for i := 1; i < 3; i++ {
					simplex[i][0] = simplex[0][0] + sigma*(simplex[i][0]-simplex[0][0])
					simplex[i][1] = simplex[0][1] + sigma*(simplex[i][1]-simplex[0][1])
					f[i] = sse(simplex[i][0], simplex[i][1])
				}
			}
		}

		fmt.Printf("Iter: %d\n", iterations)
		fmt.Printf("b1: %.6f\n", simplex[0][0])
		fmt.Printf("b2: %.6f\n", simplex[0][1])
		fmt.Printf("Owibka: %.6f\n", f[0])
		fmt.Println("----------------------------------------")

		if math.Abs(f[2]-f[0]) < 0.001 {
			break
		}
	}

	return simplex[0][0], simplex[0][1], f[0], iterations
}

func main() {
	fmt.Println("Вариант 4, Петров Матвей, ИПР-22-1б")
	fmt.Println("Поиск параметров модели: y = b1 * ln(b2 * x)")
	fmt.Println("Данные:")
	for i := range x {
		fmt.Printf("  x=%.2f, y=%.3f\n", x[i], y[i])
	}
	fmt.Println()

	b1_ds, b2_ds, sse_ds, iter_ds := directSearch(2.0, 1.0, 0.1, 0.001)

	b1_sm, b2_sm, sse_sm, iter := simplexMethod(2.0, 1.0)
	fmt.Println("Best values direct search:")
	fmt.Printf("b1 = %.6f\n", b1_ds)
	fmt.Printf("b2 = %.6f\n", b2_ds)
	fmt.Printf("SSE = %.6f\n", sse_ds)
	fmt.Printf("Итераций: %d\n\n", iter_ds)

	fmt.Println("Best values simplex method:")
	fmt.Printf("b1 = %.6f\n", b1_sm)
	fmt.Printf("b2 = %.6f\n", b2_sm)
	fmt.Printf("SSE = %.6f\n", sse_sm)
	fmt.Printf("Итераций: %d\n\n", iter)
}
