package main

import "log"

func main() {
	var board Board
	board.LoadFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	board.Move(Move{FromAlgebraic("e2"), FromAlgebraic("e4")})

	board.PrintState()

	log.Println(board.Active)
	log.Println(board.Castling)
	log.Println(board.EnPassant)
	log.Println(board.FiftyMoveRule)
	log.Println(board.Fullmoves)
}
