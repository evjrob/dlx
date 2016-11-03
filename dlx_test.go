package dlx

import (
  "fmt"
  "testing"
)

func TestSolve(t *testing.T) {
  columnCount := 7
  rows := make([][]int, 6)
  rows[0] = []int{2, 4, 5}
	rows[1] = []int{0, 3, 6}
	rows[2]	= []int{1, 2, 5}
	rows[3] = []int{0, 3}
	rows[4] = []int{1, 6}
  rows[5] = []int{3, 4, 6}

  m := NewMatrix(columnCount)

  for _, row := range rows {
    m.AddRow(row)
  }

  solution, success := m.Solve()

  if success {
    fmt.Printf("Passed the Knuth paper example: %v\n", solution)
  } else {
    t.Errorf("Failed the knuth paper example.")
  }
}

func TestSolveNoSolution(t *testing.T) {
  columnCount := 7
  rows := make([][]int, 2)
  rows[0] = []int{2, 4, 5}
	rows[1] = []int{0, 3, 6}

  m := NewMatrix(columnCount)

  for _, row := range rows {
    m.AddRow(row)
  }

  solution, success := m.Solve()

  if !success {
    fmt.Printf("Returned false for success on impossible example: %v\n", solution)
  } else {
    t.Errorf("Failed the test on an impossible example.")
  }
}

func TestSolveComplete(t *testing.T) {
  columnCount := 4
  rows := make([][]int, 4)
  rows[0] = []int{0, 1}
	rows[1] = []int{2, 3}
  rows[2] = []int{0, 1, 2}
	rows[3] = []int{3}

  solutionCount := 0
  finished := false
  success := false

  m := NewMatrix(columnCount)

  for _, row := range rows {
    m.AddRow(row)
  }

  solutionChannel, successChannel := m.SolveComplete()

  for !finished {
    select {
    case solution := <-solutionChannel:
        solutionCount++
        fmt.Printf("Solution %v: %v.\n", solutionCount, solution)
      case success = <-successChannel:
        finished = true
    }
  }

  if success {
    fmt.Printf("Passed the multiple solutions example with %v solutions.\n", solutionCount)
  } else {
    t.Errorf("Failed the test on the SolveComplete example.")
  }
}

func TestSolveCompleteNoSolutions(t *testing.T) {
  columnCount := 4
  rows := make([][]int, 4)
  rows[0] = []int{0, 1}

  solutionCount := 0
  finished := false
  success := false

  m := NewMatrix(columnCount)

  for _, row := range rows {
    m.AddRow(row)
  }

  solutionChannel, successChannel := m.SolveComplete()

  for !finished {
    select {
    case solution := <-solutionChannel:
        solutionCount++
        fmt.Printf("Solution %v: %v.\n", solutionCount, solution)
      case success = <-successChannel:
        finished = true
    }
  }

  if !success {
    fmt.Printf("Returned false for the impossible SolveComplete example with %v solutions.\n", solutionCount)
  } else {
    t.Errorf("Failed the test on an impossible SolveComplete example.")
  }
}
