# Conversions to Table Top Simulator

For now We can only convert to TableTop Simulator format from deckbox.org, tappedout.net and basic text. The images are from scryfall and it just pulls the first one returned. I think that defaults to first printing.

## tappedout.net

curl https://api.mtg.fail/\?deck\=https://tappedout.net/mtg-decks/22-01-20-kess-storm

## deckbox.org

curl https://api.mtg.fail/\?deck\=https://deckbox.org/sets/2548127 

## Basic text

`curl -X POST https://api.mtg.fail -H 'Content-Type: text/plain' --data-binary @deck.txt` where `deck.txt` contents is in [this format](sample.txt)

### Pro tip
I curl and output right to the TTS library directory.

```
curl https://api.mtg.fail/\?deck\=https://tappedout.net/mtg-decks/22-01-20-kess-storm/ > ~/Library/Tabletop\ Simulator/Saves/Saved\ Objects/testing/tappedout.json

```

I am thinking I would like to build a frontend to go with this. This summer I am learning React.js so I can build the pretties.
