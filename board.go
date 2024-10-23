package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Starts at index 0
// [0] File A-H
// [1] Rank 8-1
type Square struct {
	X, Y uint8
}

type Vector2i struct {
	X, Y int
}

type Move struct {
	From, To Square
}

type MoveTable struct {
	Captures bool       // True if the move can capture
	Jumps    bool       // True if the move goes through other pieces/is a horse
	Pattern  []Vector2i // A list of squares to move from 0, 0
}

func ToAlgebraic(s Square) string {
	algebraic := []byte{'a', '8'}
	files := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	algebraic[0] = files[s.X]                    // Grabs the file
	algebraic[1] = strconv.Itoa(8 - int(s.Y))[0] // Converts the rank (uint8) to rank (ascii byte)

	return string(algebraic)
}

func FromAlgebraic(s string) Square {
	square := Square{0, 0}
	files := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	square.X = uint8(slices.Index(files, s[0])) // Grabs the file

	rank, err := strconv.Atoi(string(s[1]))
	ienil(err)
	square.Y = uint8(8 - rank) // Converts the rank (ascii byte) to rank (uint8)

	return square
}

type Board struct {
	// Array of Pieces, from A8 to H1
	Position [64]byte

	// True = white's turn, False = black's turn
	Active bool

	// [0] White kingside
	// [1] White queenside
	// [2] Black kingside
	// [3] Black queenside
	Castling [4]bool

	// En Passant target square (the square a pawn has just moved over)
	EnPassant Square

	// Number of halfmoves since the last capture or pawn advance
	FiftyMoveRule int

	// Number of fullmoves
	Fullmoves int
}

func (b Board) PrintState() {
	for i, c := range b.Position {
		if c == 0 {
			fmt.Print(".")
		} else {
			fmt.Printf("%c", c)
		}

		if i%8 == 7 {
			fmt.Print("\n")
		}
	}
}

func (b *Board) LoadFen(fen string) {
	fen_record := strings.Split(fen, " ")

	var new_position [64]byte
	var position_index int = 0

	digit_regex, err := regexp.Compile(`\d`)
	ienil(err)

	for _, c := range fen_record[0] {
		if c == '/' {
		} else {
			matched := digit_regex.MatchString(string(c))
			ienil(err)

			if matched {
				i, err := strconv.Atoi(string(c))
				ienil(err)
				position_index += i
			} else {
				new_position[position_index] = byte(c)
				position_index++
			}
		}
	}

	b.Position = new_position

	switch fen_record[1] {
	case "w":
		b.Active = true
	case "b":
		b.Active = false
	default:
		panic("FEN string corrupted: active color is neither 'w' or 'b'")
	}

	b.Castling[0] = strings.Contains(fen_record[2], "K")
	b.Castling[1] = strings.Contains(fen_record[2], "Q")
	b.Castling[2] = strings.Contains(fen_record[2], "k")
	b.Castling[3] = strings.Contains(fen_record[2], "q")

	if fen_record[3] == "-" {
		b.EnPassant = Square{255, 255}
	} else {
		b.EnPassant = FromAlgebraic(fen_record[3])
	}

	halfmoves, err := strconv.Atoi(fen_record[4])
	ienil(err)
	b.FiftyMoveRule = halfmoves

	fullmoves, err := strconv.Atoi(fen_record[5])
	ienil(err)
	b.Fullmoves = fullmoves
}

// Returns if a move is legal, a slice of extra moves to be done, and a slice of squares to be captured
func (b Board) IsMoveLegal(move Move) (bool, []Move, []Square) {
	//move_table := make(map[string][]byte)
	// switch b.Position[move.From.X+(move.From.Y*8)] {
	// case 'P':

	// }
	return false, []Move{}, []Square{}
}

// Makes move on the board and updates board state
// Fails silently if the move given is illegal
func (b *Board) Move(move Move) {
	legal, extra, captures := b.IsMoveLegal(move)
	if !legal {
		return
	}

	// Capture captured pieces
	for _, captured := range captures {
		b.Position[captured.X+(captured.Y*8)] = 0
	}

	// Move piece
	b.Position[move.To.X+(move.To.Y*8)] = b.Position[move.From.X+(move.From.Y*8)]
	b.Position[move.From.X+(move.From.Y*8)] = 0

	// Move extra returned moves (used for castling)
	for _, extra_move := range extra {
		b.Position[extra_move.To.X+(extra_move.To.Y*8)] = b.Position[extra_move.From.X+(extra_move.From.Y*8)]
		b.Position[extra_move.From.X+(extra_move.From.Y*8)] = 0
	}

	// If moved piece is white
	if strings.ToUpper(string(b.Position[move.To.X+(move.To.Y*8)])) == string(b.Position[move.To.X+(move.To.Y*8)]) {
		b.Active = false // Black's move
	} else {
		b.Active = true // White's move
		b.Fullmoves++
	}
}
