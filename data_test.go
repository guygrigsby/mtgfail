package mtgfail

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const entryJSON = `
{
      "object": "card",
      "id": "c7e6915b-2077-45aa-93e7-29b5f14beb51",
      "oracle_id": "7b3b5be0-8bec-43c9-bd61-39fd92e0d705",
      "multiverse_ids": [
        439473
      ],
      "tcgplayer_id": 153180,
      "cardmarket_id": 313994,
      "name": "Half-Orc, Half-",
      "lang": "en",
      "released_at": "2017-12-08",
      "uri": "https://api.scryfall.com/cards/c7e6915b-2077-45aa-93e7-29b5f14beb51",
      "scryfall_uri": "https://scryfall.com/card/ust/84/half-orc-half-?utm_source=api",
      "layout": "augment",
      "highres_image": true,
      "image_uris": {
        "small": "https://c1.scryfall.com/file/scryfall-cards/small/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.jpg?1572374034",
        "normal": "https://c1.scryfall.com/file/scryfall-cards/normal/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.jpg?1572374034",
        "large": "https://c1.scryfall.com/file/scryfall-cards/large/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.jpg?1572374034",
        "png": "https://c1.scryfall.com/file/scryfall-cards/png/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.png?1572374034",
        "art_crop": "https://c1.scryfall.com/file/scryfall-cards/art_crop/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.jpg?1572374034",
        "border_crop": "https://c1.scryfall.com/file/scryfall-cards/border_crop/front/c/7/c7e6915b-2077-45aa-93e7-29b5f14beb51.jpg?1572374034"
      },
      "mana_cost": "",
      "cmc": 0,
      "type_line": "Creature — Orc Warrior",
      "oracle_text": "Trample\nAt the beginning of each end step, if an opponent was dealt damage this turn,\nAugment {1}{R}{R} ({1}{R}{R}, Reveal this card from your hand: Combine it with target host. Augment only as a sorcery.)",
      "power": "+3",
      "toughness": "+1",
      "colors": [
        "R"
      ],
      "color_indicator": [
        "R"
      ],
      "color_identity": [
        "R"
      ],
      "keywords": [
        "Trample",
        "Augment"
      ],
      "legalities": {
        "standard": "not_legal",
        "future": "not_legal",
        "historic": "not_legal",
        "pioneer": "not_legal",
        "modern": "not_legal",
        "legacy": "not_legal",
        "pauper": "not_legal",
        "vintage": "not_legal",
        "penny": "not_legal",
        "commander": "not_legal",
        "brawl": "not_legal",
        "duel": "not_legal",
        "oldschool": "not_legal"
      },
      "games": [
        "paper"
      ],
      "reserved": false,
      "foil": true,
      "nonfoil": true,
      "oversized": false,
      "promo": false,
      "reprint": false,
      "variation": false,
      "set": "ust",
      "set_name": "Unstable",
      "set_type": "funny",
      "set_uri": "https://api.scryfall.com/sets/83491685-880d-41dd-a4af-47d2b3b17c10",
      "set_search_uri": "https://api.scryfall.com/cards/search?order=set&q=e%3Aust&unique=prints",
      "scryfall_set_uri": "https://scryfall.com/sets/ust?utm_source=api",
      "rulings_uri": "https://api.scryfall.com/cards/c7e6915b-2077-45aa-93e7-29b5f14beb51/rulings",
      "prints_search_uri": "https://api.scryfall.com/cards/search?order=released&q=oracleid%3A7b3b5be0-8bec-43c9-bd61-39fd92e0d705&unique=prints",
      "collector_number": "84",
      "digital": false,
      "rarity": "uncommon",
      "card_back_id": "0aeebaf5-8c7d-4636-9e82-8c27447861f7",
      "artist": "Kev Walker",
      "artist_ids": [
        "f366a0ee-a0cd-466d-ba6a-90058c7a31a6"
      ],
      "illustration_id": "ed8bb73e-fba9-4e7f-9888-217c8b635fd6",
      "border_color": "silver",
      "frame": "2015",
      "full_art": false,
      "textless": false,
      "booster": true,
      "story_spotlight": false,
      "prices": {
        "usd": "0.16",
        "usd_foil": "0.47",
        "eur": "0.02",
        "eur_foil": "0.29",
        "tix": null
      },
      "related_uris": {
        "gatherer": "https://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=439473",
        "tcgplayer_decks": "https://decks.tcgplayer.com/magic/deck/search?contains=Half-Orc%2C+Half-&page=1&utm_campaign=affiliate&utm_medium=api&utm_source=scryfall",
        "edhrec": "https://edhrec.com/route/?cc=Half-Orc%2C+Half-",
        "mtgtop8": "https://mtgtop8.com/search?MD_check=1&SB_check=1&cards=Half-Orc%2C+Half-"
      },
      "purchase_uris": {
        "tcgplayer": "https://shop.tcgplayer.com/product/productsearch?id=153180&utm_campaign=affiliate&utm_medium=api&utm_source=scryfall",
        "cardmarket": "https://www.cardmarket.com/en/Magic/Products/Search?referrer=scryfall&searchString=Half-Orc%2C+Half-&utm_campaign=card_prices&utm_medium=text&utm_source=scryfall",
        "cardhoarder": "https://www.cardhoarder.com/cards?affiliate_id=scryfall&data%5Bsearch%5D=Half-Orc%2C+Half-&ref=card-profile&utm_campaign=affiliate&utm_medium=card&utm_source=scryfall"
      }
    }
`

func TestEntry(t *testing.T) {
	b := []byte(entryJSON)
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	require.NoError(t, err)
}
