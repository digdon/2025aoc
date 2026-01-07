package main

import (
	"fmt"
	"strings"
)

func generateMatrix(machine Machine) [][]int {
	// Convert the buttons and joltages into a matrix. Each row consists of:
	//   - one column per button, indicating whether that button affects this joltage (cell > 0)
	//   - one final column indicating the joltage value
	// ie, the row contains a joltage value and which buttons affect it.
	matrix := make([][]int, len(machine.joltages))

	for row, joltage := range machine.joltages {
		matrix[row] = make([]int, len(machine.buttons)+1)
		matrix[row][len(machine.buttons)] = joltage

		// for col, button := range machine.buttons {
		// 	for _, lightIndex := range button {
		// 		if lightIndex == row {
		// 			matrix[row][col]++
		// 			break
		// 		}
		// 	}
		// }
	}

	// Apply the buttons to the matrix. Each button is a list of which light indices it toggles.
	// Each button corresponds to a column in the matrix, and each light index for each button
	// corresponds to a row in the matrix.
	for col, button := range machine.buttons {
		for _, lightIndex := range button {
			matrix[lightIndex][col]++
		}
	}

	return matrix
}

// Compute the Integer-Only Row Echelon Form (IREF) of the given matrix
func computeIREF(origMatrix [][]int) [][]int {
	// Make a copy of the original matrix to work on
	iere := make([][]int, len(origMatrix))

	for i := range origMatrix {
		iere[i] = make([]int, len(origMatrix[i]))
		copy(iere[i], origMatrix[i])
	}

	maxRows, maxCols := len(iere), len(iere[0])
	numVars := maxCols - 1
	pivotRow := 0

	// Now we start applying integer-only Gaussian elimination to the matrix.
	// Processing each column in turn, we look for the first non-zero value in the current column,
	// swap that row to the "top" (the active/pivot row), and eliminate that column from all other rows.
	for col := 0; col < numVars && pivotRow < maxRows; col++ {
		// Find the first non-zero pivot in this column
		swapRow := pivotRow

		for row := pivotRow; row < maxRows; row++ {
			if iere[row][col] != 0 {
				swapRow = row
				break
			}
		}

		if iere[swapRow][col] == 0 {
			// No pivot in this column, move to next column
			continue
		}

		// Swap the rows
		iere[pivotRow], iere[swapRow] = iere[swapRow], iere[pivotRow]

		// Make pivot positive if it's negative, and apply to the entire row
		if iere[pivotRow][col] < 0 {
			for c := range maxCols {
				iere[pivotRow][c] *= -1
			}
		}

		pivotVal := iere[pivotRow][col]

		// Eliminate column in all other rows using integer operations
		for row := range maxRows {
			if row == pivotRow {
				continue
			}

			if iere[row][col] != 0 {
				factor := iere[row][col]

				for k := range maxCols {
					// Muliply the current cell by the pivot value, and subtract
					// the pivot row cell multiplied by the current cell value (factor)
					iere[row][k] = (iere[row][k] * pivotVal) - (iere[pivotRow][k] * factor)
				}

				// Reduce row by GCD (if applicable)
				g := calculateRowGCD(iere[row])

				if g > 1 {
					for k := range maxCols {
						iere[row][k] /= g
					}
				}
			}
		}

		pivotRow++
	}

	// Drop any all-zero rows at the bottom
	for i := maxRows - 1; i >= 0; i-- {
		allZero := true

		for j := maxCols - 1; j >= 0; j-- {
			if iere[i][j] != 0 {
				allZero = false
				break
			}
		}

		if allZero {
			iere = iere[:i]
		} else {
			break
		}
	}

	return iere
}

func calculateRowGCD(row []int) int {
	g := 0

	for i := range len(row) {
		g = gcd(g, abs(row[i]))
	}

	return g
}

func gcd(a, b int) int {
	a = abs(a)
	b = abs(b)

	for b > 0 {
		a, b = b, a%b
	}

	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generatePivotVars(matrix [][]int) []int {
	maxRows, maxCols := len(matrix), len(matrix[0])
	numVars := maxCols - 1
	pivotVars := make([]int, numVars)

	for i := range pivotVars {
		pivotVars[i] = -1
	}

	// Identify pivot variables
	for row, col := 0, 0; col < numVars && row < maxRows; col++ {
		if matrix[row][col] == 1 {
			pivotVars[col] = row
			row++
		}
	}

	return pivotVars
}

func lookForQuickSolution(origMatrix [][]int, rref [][]int, pivotVars []int) int {
	// Does the RREF have any free variables? If NOT, we can directly calculate the solution.
	hasFreeVars := false
	for _, pv := range pivotVars {
		if pv == -1 {
			hasFreeVars = true
			break
		}
	}

	if !hasFreeVars {
		maxRows, valuePos := len(rref), len(rref[0])-1
		presses := 0

		for row := range maxRows {
			presses += rref[row][valuePos]
		}

		return presses
	}

	// Are there any rows in the original matrix that have all buttons pressed?
	maxRows, maxCols := len(origMatrix), len(origMatrix[0])

	for row := range maxRows {
		allPressed := true

		for col := 0; col < maxCols-1; col++ {
			if origMatrix[row][col] == 0 {
				allPressed = false
				break
			}
		}

		if allPressed {
			// This row has all buttons pressed, so the solution is just the constant term
			return origMatrix[row][maxCols-1]
		}
	}

	return -1
}

func printEquations1(matrix [][]int, pivotVars []int) {
	maxCols := len(matrix[0])
	numVars := maxCols - 1

	// Print equations
	for i, pv := range pivotVars {
		if pv == -1 {
			fmt.Printf("x%d is free\n", i+1)
		} else {
			constant := matrix[pv][maxCols-1]
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("x%d", i+1))

			// Check for dependencies on other (free?) variables
			for j := i + 1; j < numVars; j++ {
				val := matrix[pv][j]

				if val != 0 {
					if val < 0 {
						sb.WriteString(" - ")
						val = -val
					} else if val > 0 {
						sb.WriteString(" + ")
					}

					if val > 1 {
						sb.WriteString(fmt.Sprintf("%d", val))
					}

					sb.WriteString(fmt.Sprintf("x%d", j+1))
				}
			}

			sb.WriteString(fmt.Sprintf(" = %d", constant))
			fmt.Println(sb.String())
		}
	}
}
