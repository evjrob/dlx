# dlx

dlx is a go package that is based off of Donald Knuth's dancing links implementation of algorithm X for solving the exact cover problem.
It allows you to define a dlx matrix of a given number of columns and then add rows in the form of an int slice of the column indeces that contain a one. A solution can then be found using Solve() or SolveComplete() depending if only a single solution is sought or if all possible solutions are wanted.

The key methods exposed are:
- NewMatrix(columnCount int)
- (matrix) AddRow(rowIndeces []int{})
- (matrix) Solve()
- (matrix) SolveComplete()

Solutions are in the form of a map[int][]int where the map keys are the row numbers based on when they were added, and the int slice at a key are the column indeces of the 1 values in that row.

## example

To solve the example matrix provided in [Donald Knuth's paper](https://arxiv.org/abs/cs/0011047) on the dancing links algorithm:

```
m := NewMatrix(7)

rows := make([][]int, 6)
rows[0] = []int{2, 4, 5}
rows[1] = []int{0, 3, 6}
rows[2]	= []int{1, 2, 5}
rows[3] = []int{0, 3}
rows[4] = []int{1, 6}
rows[5] = []int{3, 4, 6}

for _, row := range rows {
  m.AddRow(row)
}

solution, success := m.Solve()
```

Which will return:

```
map[3:[0 3] 0:[4 5 2] 4:[1 6]]
```
