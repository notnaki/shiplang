
// Define the board as a 3x3 array
let board = [][]string{
    []string{" "," "," "},
    []string{" "," "," "},
    []string{" "," "," "}
};

// Function to display the current state of the board
fn displayBoard() {
    foreach (row in board) {
        show(row);
    }
}

// Function to check if a player has won
fn checkWin(player) {
    // Check rows
    foreach (i in range(3)) {
        if (board[i][0] == player && board[i][1] == player && board[i][2] == player) {
            return true;
        }
    }
    // Check columns
    foreach (i in range(3)) {
        if (board[0][i] == player && board[1][i] == player && board[2][i] == player) {
            return true;
        }
    }
    // Check diagonals
    if (board[0][0] == player && board[1][1] == player && board[2][2] == player) {
        return true;
    }
    if (board[0][2] == player && board[1][1] == player && board[2][0] == player) {
        return true;
    }
    return false;
}

// Function to check if the board is full
fn checkDraw() {
    foreach (row in board) {
        foreach (cell in row) {
            if (cell == null) {
                return false;
            }
        }
    }
    return true;
}

// Function to make a move
fn makeMove(row, col, player) {
    if (row < 0 || row >= 3 || col < 0 || col >= 3 || board[row][col] != null) {
        return false; // Invalid move
    }
    board[row][col] = player;
    return true;
}

// Main game loop
while (true) {
    displayBoard();
    let row = ask("Enter row (0-2): ").toNum();
    let col = ask("Enter column (0-2): ").toNum();
    if (makeMove(row, col, "X")) {
        show("Invalid move! Try again.");
    }
    if (checkWin("X")) {
        show("Player X wins!");
        break;
    }
    if (checkDraw()) {
        show("It's a draw!");
        break;
    }
    

}
