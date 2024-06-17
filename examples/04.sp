let board = [][]string{
    []string{" "," "," "},
    []string{" "," "," "},
    []string{" "," "," "}
};



let x = 0;
let y = 0;

fn makeMove(sym:string,x:number,y:number, board:[][]string) {
    let row = board[y];
    row[x] = sym;
    board[y] = row;
    return board;
}

fn showBoard(board: [][]string) {
    foreach (r in board){
        show(r);
    }
}

while(true){
    showBoard(board);
     x = ask("row? ").toNum();
     y = ask("col? ").toNum();

     makeMove("x",x,y,board);

     showBoard(board);
     x = ask("row? ").toNum();
     y = ask("col? ").toNum();

     makeMove("o",x,y,board);
}