package main

func main() {
	cards := deck{"Ace of Dimonds", newCard()}
	cards = append(cards, "Six of Spades")

	// fmt.Println(cards)
	// for i, card := range cards {
	// 	fmt.Println(i, card)
	// }
	cards.print()
}

func newCard() string {
	return "Five of Diamonds"
}
