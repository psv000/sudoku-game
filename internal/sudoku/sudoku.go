package sudoku

import (
	"math/rand"
)

const gridSize = 9

type Board [gridSize][gridSize]int

type PuzzleResponse struct {
    Puzzle   Board `json:"puzzle"`
    Solution Board `json:"solution"`
}

// GenerateSudoku создает новую головоломку
func GenerateSudoku(difficulty string) (Board, Board) {
    var board Board
    var solution Board

    fillDiagonal(&board)
    solveSudoku(&board)
    solution = board

    var holes int
    switch difficulty {
    case "beginner":
        holes = 30 + rand.Intn(5) // 30-34 empty cells
    case "easy":
        holes = 40 + rand.Intn(5) // 40-44
    case "medium":
        holes = 48 + rand.Intn(5) // 48-52
    case "hard":
        holes = 53 + rand.Intn(4) // 53-56
    case "expert":
        holes = 57 + rand.Intn(3) // 57-59
    default:
        holes = 48 // mid by default
    }
    removeDigits(&board, holes)

    return board, solution
}

// isSafe checks if a number can be placed in a given cell
func isSafe(board *Board, row, col, num int) bool {
    // Rows check
    for x := 0; x < gridSize; x++ {
        if board[row][x] == num {
            return false
        }
    }
    // Columns check
    for x := 0; x < gridSize; x++ {
        if board[x][col] == num {
            return false
        }
    }
    // 3x3 check
    startRow := row - row%3
    startCol := col - col%3
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if board[i+startRow][j+startCol] == num {
                return false
            }
        }
    }
    return true
}

// solveSudoku solves Sudoku with backtracking
func solveSudoku(board *Board) bool {
    for i := 0; i < gridSize; i++ {
        for j := 0; j < gridSize; j++ {
            if board[i][j] == 0 {
                for num := 1; num <= gridSize; num++ {
                    if isSafe(board, i, j, num) {
                        board[i][j] = num
                        if solveSudoku(board) {
                            return true
                        }
                        board[i][j] = 0 // backtrack
                    }
                }
                return false
            }
        }
    }
    return true
}

// fillDiagonal fills diagonal 3x3 blocks
func fillDiagonal(board *Board) {
    for i := 0; i < gridSize; i = i + 3 {
        fillBox(board, i, i)
    }
}

// fillBox fills block 3x3 random digits
func fillBox(board *Board, row, col int) {
    var num int
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            for {
                num = rand.Intn(gridSize) + 1
                if isSafe(board, row, col, num) {
                    break
                }
            }
            board[row+i][col+j] = num
        }
    }
}

// removeDigits removes K digits from grid
func removeDigits(board *Board, k int) {
    count := k
    for count != 0 {
        i := rand.Intn(gridSize)
        j := rand.Intn(gridSize)
        if board[i][j] != 0 {
            count--
            board[i][j] = 0
        }
    }
}