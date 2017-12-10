package hanabi

import (
	"fmt"
	"math/rand"
)

const (
	maxHintCount = 8
	bombCount    = 3
)

var (
	ErrOutOfTurn  = fmt.Errorf("player played card out of turn")
	ErrOutOfBombs = fmt.Errorf("Out of Bombs")
	ErrGameOver   = fmt.Errorf("Game ended")
)

// TODO: represent hints (struct?)

// Board holds game state
type Board struct {
	Piles              map[Color]Number
	Deck               []*Card
	Discard            []*Card
	Players            []*Player
	CurrentPlayerIndex int
	Hints              int
	Bombs              int
}

func NewBoard(n int) *Board {
	b := &Board{
		Piles: make(map[Color]Number),
		Deck:  newDeck(),
		Hints: maxHintCount,
		Bombs: bombCount,
	}
	for i := 0; i < n; i++ {
		b.Players = append(b.Players, &Player{
			ID:   i,
			Hand: b.drawCards(5), // TODO: card counts per player count
		})
	}
	return b
}

func (b *Board) card(pi, ci int) *Card {
	// TODO: check if pi and ci are valid
	player := b.Players[pi]
	card := player.Hand[ci]
	copy(player.Hand[ci:], player.Hand[ci+1:])
	newCard, ok := b.drawCard()
	player.Hand[len(player.Hand)-1] = newCard
	return card
}

func (b *Board) GiveHint(pi int, ci int) error {
	if !b.decrementHint() {
		return ErrNoHintsLeft
	}
	return nil
}

func (b *Board) Play(pi int, ci int) error {
	if b.CurrentPlayerIndex != pi {
		return ErrOutOfTurn
	}
	card := b.card(pi, ci)
	color := card.C
	pileNumber := b.Piles[color]
	if pileNumber+1 == card.N {
		b.Piles[color]++
		if card.N == 5 {
			b.incrementHint()
		}
	} else {
		b.Bombs--
		b.Discard = append(b.Discard, card)
		if b.Bombs <= 0 {
			return ErrOutOfBombs
		}
	}
	if len(b.Deck) == 0 {
		return ErrGameOver
	}
	return nil
}

func (b *Board) Discard(pi int, ci int) error {
	if b.CurrentPlayerIndex != pi {
		return ErrOutOfTurn
	}
	card := b.card(pi, ci)
	b.incrementHint()
	b.Discard = append(b.Discard, card)
	if len(b.Deck) == 0 {
		return ErrGameOver
	}
	return nil
}

func (b *Board) decrementHint() bool {
	if b.Hints > 0 {
		b.Hints--
		return true
	}
	return false
}

func (b *Board) incrementHint() {
	if b.Hints < maxHintCount {
		b.Hints++
	}
}

func (b *Board) drawCards(n int) []*Card {
	var cards []*Card
	for i := 0; i < n; i++ {
		card := b.drawCard()
		cards = append(cards, b.drawCard())
	}
	return cards
}

func (b *Board) drawCard() *Card {
	card := b.Deck[0]
	b.Deck = b.Deck[1:]
	return card
}

func newDeck() []*Card {
	var deck []*Card
	counts := []int{0, 3, 2, 2, 2, 1}
	for color := Blue; color <= Yellow; color++ {
		// 3x1, 2x2, 2x3, 2x4, 1x5
		for num, cnt := range counts {
			for i := 0; i < cnt; i++ {
				deck = append(deck, &Card{
					C: Color(color),
					N: Number(num),
				})
			}
		}
	}

	shuffledIndices := rand.Perm(len(deck))
	shuffledDeck := make([]*Card, len(deck))
	for i, v := range shuffledIndices {
		shuffledDeck[v] = deck[i]
	}

	return shuffledDeck
}

// Color is color of a card
type Color int

// The Hanabi Colors
const (
	NoColor = iota
	Blue
	Green
	Red
	White
	Yellow
)

// Number is the number on a card
type Number int

// Card is a card.
type Card struct {
	C Color  `json:"c"`
	N Number `json:"n"`
}

// Player represents a player
type Player struct {
	ID   int
	Hand []*Card
}
