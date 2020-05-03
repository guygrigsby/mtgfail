
up:
	go run cmd/mtgfail/main.go -file ./deck.txt -bulk ./scryfall-default-cards.json
	cp out.json ~/Library/Tabletop\ Simulator/Saves/Saved\ Objects/testing
