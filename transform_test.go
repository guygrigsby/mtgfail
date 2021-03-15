package mtgfail

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

func TestTappedOutTransform(t *testing.T) {
	log := log15.New()

	r, err := Normalize(ioutil.NopCloser(strings.NewReader(input)), log)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	fmt.Printf("BYTES %s", b)

}
func TestDeckboxTransform(t *testing.T) {
	log := log15.New()

	r, err := Normalize(ioutil.NopCloser(strings.NewReader(input)), log)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	fmt.Printf("BYTES %s", b)

}

var input = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <meta name="description" content="Card set print view" />
    <title>Deckbox.org card set print view</title>
</head>
<body>
        1 Bloodstained Mire<br/>1 Command Tower<br/>1 Fabled Passage<br/>1 Flooded Strand<br/>1 Great Furnace<br/>19 Island<br/>5 Mountain<br/>1 Mystic Sanctuary<br/>1 Polluted Delta<br/>1 Seat of the Synod<br/>1 Steam Vents<br/>1 Sulfur Falls<br/>1 Temple of the False God<br/>1 Wooded Foothills<br/>1 Anticipate<br/>1 Arcane Signet<br/>1 Argentum Armor<br/>1 Armory Automaton<br/>1 Basilisk Collar<br/>1 Batterskull<br/>1 Bloodforged Battle-Axe<br/>1 Brainstorm<br/>1 Brass Squire<br/>1 Chaos Warp<br/>1 Counterspell<br/>1 Cyclonic Rift<br/>1 Dalakos, Crafter of Wonders<br/>1 Deadeye Quartermaster<br/>1 Divination<br/>1 Dramatic Reversal<br/>1 Echo Storm<br/>1 Emry, Lurker of the Loch<br/>1 Etherium Sculptor<br/>1 Fabricate<br/>1 Foundry Inspector<br/>1 Goblin Engineer<br/>1 Hammer of Nazahn<br/>1 Herald of Kozilek<br/>1 Isochron Scepter<br/>1 Izzet Signet<br/>1 Jhoira, Weatherlight Captain<br/>1 Kazuul's Toll Collector<br/>1 Mask of Memory<br/>1 Mind Stone<br/>1 Mirran Spy<br/>1 Mirrormade<br/>1 Mystic Repeal<br/>1 Narset's Reversal<br/>1 Negate<br/>1 Nettle Drone<br/>1 Neurok Stealthsuit<br/>1 Rapid Hybridization<br/>1 Reality Shift<br/>1 Rhystic Study<br/>1 Saheeli, Sublime Artificer<br/>1 Saheeli, the Gifted<br/>1 Sai, Master Thopterist<br/>1 Scytheclaw<br/>1 Shadowspear<br/>1 Skullclamp<br/>1 Sol Ring<br/>1 Solemn Simulacrum<br/>1 Storm the Vault<br/>1 Swan Song<br/>1 Sword of Feast and Famine<br/>1 Sword of the Animist<br/>1 Talisman of Creativity<br/>1 Tezzeret the Seeker<br/>1 Tezzeret, Artifice Master<br/>1 Thopter Spy Network<br/>1 Thoughtcast<br/>1 Trailblazer's Boots<br/>1 Trinket Mage<br/>1 Trophy Mage<br/>1 Unwinding Clock<br/>1 Vandalblast<br/>1 Vedalken Engineer<br/>1 Zephyr Scribe<br/>

    <p><strong>Sideboard:</strong></p>
    1 Displace<br/>1 Ghostly Flicker<br/>

</body>
</html>`
