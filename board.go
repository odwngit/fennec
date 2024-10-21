package main

import (
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Starts at index 0
// [0] File A-H
// [1] Rank 8-1
type Square [2]uint8

func (s Square) ToAlgebraic() string {
	algebraic := []byte{'a', '8'}
	files := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	algebraic[0] = files[s[0]]                    // Grabs the file
	algebraic[1] = strconv.Itoa(8 - int(s[1]))[0] // Converts the rank (uint8) to rank (ascii byte)

	return string(algebraic)
}

func FromAlgebraic(s string) Square {
	square := [2]uint8{0, 0}
	files := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}

	square[0] = files[slices.Index(files, s[0])] // Grabs the file
	rank, err := strconv.Atoi(string(s[1]))
	ienil(err)
	square[1] = uint8(8 - rank) // Converts the rank (ascii byte) to rank (uint8)

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

func (b *Board) LoadFen(fen string) {
	fen_record := strings.Split(fen, " ")

	var new_position [64]byte
	var position_index int = 0

	digit_regex, err := regexp.Compile(`\d`)
	ienil(err)

	for _, c := range fen_record[0] {
		log.Printf("Processing %c...\n", c)
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
