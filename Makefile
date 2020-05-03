
serve:
	go run cmd/server/main.go -bulk ./scryfall-default-cards.json
up:
	go run cmd/mtgfail/main.go -file ./deck.txt -bulk ./scryfall-default-cards.json
	cp out.json ~/Library/Tabletop\ Simulator/Saves/Saved\ Objects/testing

