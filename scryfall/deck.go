package scryfall

type Deck struct {
	Object       string      `json:"object"`
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	URI          string      `json:"uri"`
	ScryfallURI  string      `json:"scryfall_uri"`
	Description  interface{} `json:"description"`
	Trashed      bool        `json:"trashed"`
	InCompliance bool        `json:"in_compliance"`
	Sections     struct {
		Primary   []string `json:"primary"`
		Secondary []string `json:"secondary"`
	} `json:"sections"`
	Entries struct {
		Commanders []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Foil              bool    `json:"foil"`
			CardDigest        struct {
				Object          string `json:"object"`
				ID              string `json:"id"`
				OracleID        string `json:"oracle_id"`
				Name            string `json:"name"`
				ScryfallURI     string `json:"scryfall_uri"`
				ManaCost        string `json:"mana_cost"`
				TypeLine        string `json:"type_line"`
				CollectorNumber string `json:"collector_number"`
				Set             string `json:"set"`
				ImageUris       struct {
					Front string `json:"front"`
				} `json:"image_uris"`
			} `json:"card_digest"`
		} `json:"commanders"`
		Lands []struct {
			Object            string      `json:"object"`
			ID                string      `json:"id"`
			DeckID            string      `json:"deck_id"`
			Section           string      `json:"section"`
			Cardinality       float64     `json:"cardinality"`
			Count             int         `json:"count"`
			RawText           string      `json:"raw_text"`
			Found             bool        `json:"found"`
			PrintingSpecified bool        `json:"printing_specified"`
			Foil              bool        `json:"foil"`
			CardDigest        interface{} `json:"card_digest"`
		} `json:"lands"`
		Outside []struct {
			Object            string      `json:"object"`
			ID                string      `json:"id"`
			DeckID            string      `json:"deck_id"`
			Section           string      `json:"section"`
			Cardinality       float64     `json:"cardinality"`
			Count             int         `json:"count"`
			RawText           string      `json:"raw_text"`
			Found             bool        `json:"found"`
			PrintingSpecified bool        `json:"printing_specified"`
			Foil              bool        `json:"foil"`
			CardDigest        interface{} `json:"card_digest"`
		} `json:"outside"`
		Maybeboard []struct {
			Object            string      `json:"object"`
			ID                string      `json:"id"`
			DeckID            string      `json:"deck_id"`
			Section           string      `json:"section"`
			Cardinality       float64     `json:"cardinality"`
			Count             int         `json:"count"`
			RawText           string      `json:"raw_text"`
			Found             bool        `json:"found"`
			PrintingSpecified bool        `json:"printing_specified"`
			Foil              bool        `json:"foil"`
			CardDigest        interface{} `json:"card_digest"`
		} `json:"maybeboard"`
		Nonlands []struct {
			Object            string  `json:"object"`
			ID                string  `json:"id"`
			DeckID            string  `json:"deck_id"`
			Section           string  `json:"section"`
			Cardinality       float64 `json:"cardinality"`
			Count             int     `json:"count"`
			RawText           string  `json:"raw_text"`
			Found             bool    `json:"found"`
			PrintingSpecified bool    `json:"printing_specified"`
			Foil              bool    `json:"foil"`
			CardDigest        struct {
				Object          string `json:"object"`
				ID              string `json:"id"`
				OracleID        string `json:"oracle_id"`
				Name            string `json:"name"`
				ScryfallURI     string `json:"scryfall_uri"`
				ManaCost        string `json:"mana_cost"`
				TypeLine        string `json:"type_line"`
				CollectorNumber string `json:"collector_number"`
				Set             string `json:"set"`
				ImageUris       struct {
					Front string `json:"front"`
				} `json:"image_uris"`
			} `json:"card_digest"`
		} `json:"nonlands"`
	} `json:"entries"`
}
