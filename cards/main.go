package main

func main() {
	cards := newDeck()

	// cards.print()
	hand, remainingCards := deal(cards, 5)

	hand.print()
	remainingCards.print()
}

func newCard() string {
	return "Five of Diamonds"
}
