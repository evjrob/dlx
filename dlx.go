package dlx

// Matrix is a 2D representation of the exact cover problem consisting of 1s and 0s
type Matrix interface {
	// AddRow allows the puzzle to be built up by adding one row at a time.
	// A row is an int slice of column indeces for the "1" elements in the matrix
	AddRow([]int) int

	// Solve attempts to find a subset of the rows that contain a single 1 in
	// each of the columns of the puzzle. It returns when the first solution is found.
	Solve() (map[int][]int, bool)

	// SolveComplete attempts to find all subsets of the rows that satisfy the exact
	// cover problem. It provides a solutions channel and a success channel
	SolveComplete() (<-chan map[int][]int, <-chan bool)
}

// matrix is the 2D circular doubly linked list where a node has pointers in the
// up, down, left, and right directions to its neighbouring nodes.
type matrix struct {
	root         *columnHeader //The left most column header that never has nodes assigned to it.
	columns      []columnHeader
	solutionRows []matrixElement
	rowCount     int
}

// A matrix element is either a node or a columnHeader in the dlx matrix
type matrixElement interface {
	rowIndex() int
	setRowIndex(int)
	parentColumn() *columnHeader
	setParentColumn(*columnHeader)
	up() matrixElement
	down() matrixElement
	left() matrixElement
	right() matrixElement
	setUp(matrixElement)
	setDown(matrixElement)
	setLeft(matrixElement)
	setRight(matrixElement)
}

// node is an element of the matrix signifying a one
type node struct {
	rowInd                            int
	parentColumnPtr                   *columnHeader
	upPtr, downPtr, leftPtr, rightPtr matrixElement
}

func (n *node) rowIndex() int {
	return n.rowInd
}
func (n *node) setRowIndex(i int) {
	n.rowInd = i
}
func (n *node) parentColumn() *columnHeader {
	return n.parentColumnPtr
}
func (n *node) setParentColumn(c2 *columnHeader) {
	n.parentColumnPtr = c2
}
func (n *node) up() matrixElement {
	return n.upPtr
}
func (n *node) down() matrixElement {
	return n.downPtr
}
func (n *node) left() matrixElement {
	return n.leftPtr
}
func (n *node) right() matrixElement {
	return n.rightPtr
}
func (n *node) setUp(e matrixElement) {
	n.upPtr = e
}
func (n *node) setDown(e matrixElement) {
	n.downPtr = e
}
func (n *node) setLeft(e matrixElement) {
	n.leftPtr = e
}
func (n *node) setRight(e matrixElement) {
	n.rightPtr = e
}

// column is a single column in the matrix struct
type columnHeader struct {
	node      // column headers are just special nodes.
	columnLen int
	columnInd int
}

func (c *columnHeader) columnIndex() int {
	return c.columnInd
}
func (c *columnHeader) setColumnIndex(i int) {
	c.columnInd = i
}
func (c *columnHeader) len() int {
	return c.columnLen
}
func (c *columnHeader) incrementLen() {
	c.columnLen++
}
func (c *columnHeader) decrementLen() {
	c.columnLen--
}

// NewMatrix acts as a constructor to return a functioning but empty instantiation
// of a matrix that contains only they number of columnHeaders specified
func NewMatrix(columnCount int) Matrix {
	var solutionRows []matrixElement
	columnHeaders := make([]columnHeader, columnCount+1)
	for i := 0; i < columnCount+1; i++ {
		head := &columnHeaders[i]
		head.setRowIndex(-1)
		head.setParentColumn(head)
		head.setUp(head)
		head.setDown(head)
		if i == 0 {
			head.setLeft(&columnHeaders[columnCount]) // left most heads's left pointer is to the right most head
		} else {
			head.setLeft(&columnHeaders[i-1])
		}
		if i == columnCount {
			head.setRight(&columnHeaders[0]) // right most head's right point is to the left most head
		} else {
			head.setRight(&columnHeaders[i+1])
		}
		head.setColumnIndex(i - 1)

	}
	return &matrix{&columnHeaders[0], columnHeaders, solutionRows, 0}
}

// AddRow allows a constraint row in the form of an int slice of node indeces to be added.
func (m *matrix) AddRow(rowIndeces []int) int {
	nodeCount := len(rowIndeces)
	row := make([]node, nodeCount)
	for i, c := range rowIndeces {
		currentNode := &row[i]
		column := &m.columns[c+1]
		downNode := column.up()
		currentNode.setRowIndex(m.rowCount)
		currentNode.setParentColumn(column)
		column.incrementLen() // Increase the length of this column by 1
		// Set the currentNode up and down pointers
		currentNode.setUp(downNode)
		currentNode.setDown(column)
		// Set the neighbouring nodes above and below to point to currentNode
		downNode.setDown(currentNode)
		column.setUp(currentNode)
		if i == 0 {
			currentNode.setLeft(&row[nodeCount-1]) // left most node's left pointer is to the right most node
		} else {
			currentNode.setLeft(&row[i-1])
		}
		if i == nodeCount-1 {
			currentNode.setRight(&row[0]) // right most node's right point is to the left most node
		} else {
			currentNode.setRight(&row[i+1])
		}
	}
	m.rowCount++

	return m.rowCount - 1
}

// cover a column in the matrix
func cover(column *columnHeader) {
	// Bypass the column
	leftColumn := column.left()
	rightColumn := column.right()
	leftColumn.setRight(rightColumn)
	rightColumn.setLeft(leftColumn)

	// For each node down the column
	downNode := column.down()
	for downNode != column {
		// For each node right
		rightNode := downNode.right()
		for rightNode != downNode {
			// Bypass the rightNode vertically and reduce the column len by one
			aboveRightNode := rightNode.up()
			belowRightNode := rightNode.down()
			aboveRightNode.setDown(belowRightNode)
			belowRightNode.setUp(aboveRightNode)
			rightNode.parentColumn().decrementLen()
			rightNode = rightNode.right()
		}
		downNode = downNode.down()
	}
}

// uncover a column in the matrix
func uncover(column *columnHeader) {
	// For each node up the column (from the bottom)
	upNode := column.up()
	for upNode != column {
		leftNode := upNode.left()
		// For each node left
		for leftNode != upNode {
			// Put this node back into its column and increase the column len by one
			aboveLeftNode := leftNode.up()
			belowLeftNode := leftNode.down()
			aboveLeftNode.setDown(leftNode)
			belowLeftNode.setUp(leftNode)
			leftNode.parentColumn().incrementLen()
			leftNode = leftNode.left()
		}
		upNode = upNode.up()
	}
	// Put the column back into the columnHeaders linked list
	leftColumn := column.left()
	rightColumn := column.right()
	leftColumn.setRight(column)
	rightColumn.setLeft(column)
}

func (m *matrix) getSolution() (solution map[int][]int) {
	solution = make(map[int][]int)
	for _, rowStart := range m.solutionRows {
		rowIndex := rowStart.rowIndex()
		solution[rowIndex] = []int{rowStart.parentColumn().columnIndex()}
		rowNode := rowStart.right()
		for rowNode != rowStart {
			solution[rowIndex] = append(solution[rowIndex], rowNode.parentColumn().columnIndex())
			rowNode = rowNode.right()
		}
	}
	return solution
}

// Search attempts to find an or all exact covers. It returns solutions over the
// results channel and has boolean channel to indicate if a solution has been
// found when the search is complete.
func (m *matrix) search(solutions chan map[int][]int, solutionFound chan bool) {
	matrixRoot := m.root

	// If only the root is left, a solution has been found
	if matrixRoot.right() == matrixRoot {
		solutions <- m.getSolution()
		solutionFound <- true
		return
	}

	// Pick the left most column of minimum length
	currentColumn := matrixRoot.right()
	minColumn := matrixRoot.right()
	for currentColumn != matrixRoot {
		if currentColumn.parentColumn().len() < minColumn.parentColumn().len() {
			minColumn = currentColumn
		}
		currentColumn = currentColumn.right()
	}

	cover(minColumn.parentColumn())

	// For each node down the covered minColumn
	childSolutionFound := false
	downNode := minColumn.down()
	for downNode != minColumn {
		m.solutionRows = append(m.solutionRows, downNode) // Add this row to the solution
		rightNode := downNode.right()
		// Cover the column of each node to the right of the current downNode
		for rightNode != downNode {
			cover(rightNode.parentColumn())
			rightNode = rightNode.right()
		}
		// Call search again on the remaining uncovered portion of the matrix.
		childSuccessChannel := make(chan bool)
		go m.search(solutions, childSuccessChannel)
		childSuccess := <-childSuccessChannel
		if childSuccess {
			childSolutionFound = true
		}
		// Uncover the columns
		leftNode := downNode.left()
		for leftNode != downNode {
			uncover(leftNode.parentColumn())
			leftNode = leftNode.left()
		}
		downNode = downNode.down()
		m.solutionRows = m.solutionRows[:len(m.solutionRows)-1] // remove this row from the solution
	}
	uncover(minColumn.parentColumn())
	solutionFound <- childSolutionFound // send the finished signal over the channel when done.
	return
}

// Solve returns a solution and a boolean to indicate if a solution was found.
// The function will return true when the first viable solution is found or false
// after all possibilities have been searched.
func (m *matrix) Solve() (map[int][]int, bool) {
	var solution map[int][]int
	solutions := make(chan map[int][]int)
	finished := make(chan bool)
	success := false
	go m.search(solutions, finished)
	select {
	case solution = <-solutions:
		success = true
	case success = <-finished:
		if !success {
			solution = make(map[int][]int)
		} else {
			solution = <-solutions
		}
	}
	return solution, success
}

// SolveComplete returns two channels. One contains all possible solutions and the
// other returns either true or false at the end of the search to indicate if any
// solutions were found.
func (m *matrix) SolveComplete() (<-chan map[int][]int, <-chan bool) {
	solutions := make(chan map[int][]int)
	finished := make(chan bool)
	go m.search(solutions, finished)

	return solutions, finished
}
