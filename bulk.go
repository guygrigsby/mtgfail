package mtgfail

type Bulk map[string]*Entry

type Entry struct {
	Object          string        `json:"object"`
	ID              string        `json:"id"`
	OracleID        string        `json:"oracle_id"`
	MultiverseIds   []interface{} `json:"multiverse_ids"`
	Name            string        `json:"name"`
	Lang            string        `json:"lang"`
	ReleasedAt      string        `json:"released_at"`
	URI             string        `json:"uri"`
	ScryfallURI     string        `json:"scryfall_uri"`
	Layout          string        `json:"layout"`
	HighresImage    bool          `json:"highres_image"`
	ImageUris       ImageUris     `json:"image_uris"`
	ManaCost        string        `json:"mana_cost"`
	Cmc             float64       `json:"cmc"`
	TypeLine        string        `json:"type_line"`
	OracleText      string        `json:"oracle_text"`
	Colors          []string      `json:"colors"`
	ColorIdentity   []string      `json:"color_identity"`
	Legalities      Legalities    `json:"legalities"`
	Games           []string      `json:"games"`
	Reserved        bool          `json:"reserved"`
	Foil            bool          `json:"foil"`
	Nonfoil         bool          `json:"nonfoil"`
	Oversized       bool          `json:"oversized"`
	Promo           bool          `json:"promo"`
	Reprint         bool          `json:"reprint"`
	Variation       bool          `json:"variation"`
	Set             string        `json:"set"`
	SetName         string        `json:"set_name"`
	SetType         string        `json:"set_type"`
	SetURI          string        `json:"set_uri"`
	SetSearchURI    string        `json:"set_search_uri"`
	ScryfallSetURI  string        `json:"scryfall_set_uri"`
	RulingsURI      string        `json:"rulings_uri"`
	PrintsSearchURI string        `json:"prints_search_uri"`
	CollectorNumber string        `json:"collector_number"`
	Digital         bool          `json:"digital"`
	Rarity          string        `json:"rarity"`
	CardBackID      string        `json:"card_back_id"`
	Artist          string        `json:"artist"`
	ArtistIds       []string      `json:"artist_ids"`
	IllustrationID  string        `json:"illustration_id"`
	BorderColor     string        `json:"border_color"`
	Frame           string        `json:"frame"`
	FullArt         bool          `json:"full_art"`
	Textless        bool          `json:"textless"`
	Booster         bool          `json:"booster"`
	StorySpotlight  bool          `json:"story_spotlight"`
	EdhrecRank      int           `json:"edhrec_rank"`
	RelatedUris     RelatedUris   `json:"related_uris"`
}
type ImageUris struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	Png        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}
type Legalities struct {
	Standard  string `json:"standard"`
	Future    string `json:"future"`
	Historic  string `json:"historic"`
	Pioneer   string `json:"pioneer"`
	Modern    string `json:"modern"`
	Legacy    string `json:"legacy"`
	Pauper    string `json:"pauper"`
	Vintage   string `json:"vintage"`
	Penny     string `json:"penny"`
	Commander string `json:"commander"`
	Brawl     string `json:"brawl"`
	Duel      string `json:"duel"`
	Oldschool string `json:"oldschool"`
}
type RelatedUris struct {
	TcgplayerDecks string `json:"tcgplayer_decks"`
	Edhrec         string `json:"edhrec"`
	Mtgtop8        string `json:"mtgtop8"`
}
