
serve:
	go run cmd/server/main.go -bulk ./scryfall-default-cards-local.json
up:
	go run cmd/mtgfail/main.go -file .examples/deck.txt -bulk ./scryfall-default-cards-local.json
	cp out.json ~/Library/Tabletop\ Simulator/Saves/Saved\ Objects/testing
test:
	docker build -t testing-mtg . -f Dockerfile.test

.PHONY: test
