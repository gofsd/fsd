// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    deck, err := UnmarshalDeck(bytes)
//    bytes, err = deck.Marshal()
//
//    en, err := UnmarshalEn(bytes)
//    bytes, err = en.Marshal()
//
//    expression, err := UnmarshalExpression(bytes)
//    bytes, err = expression.Marshal()
//
//    filter, err := UnmarshalFilter(bytes)
//    bytes, err = filter.Marshal()
//
//    lang, err := UnmarshalLang(bytes)
//    bytes, err = lang.Marshal()
//
//    tag, err := UnmarshalTag(bytes)
//    bytes, err = tag.Marshal()
//
//    topWordsByGoogle, err := UnmarshalTopWordsByGoogle(bytes)
//    bytes, err = topWordsByGoogle.Marshal()
//
//    tr, err := UnmarshalTr(bytes)
//    bytes, err = tr.Marshal()
//
//    translation, err := UnmarshalTranslation(bytes)
//    bytes, err = translation.Marshal()
//
//    user, err := UnmarshalUser(bytes)
//    bytes, err = user.Marshal()
//
//    userAccount, err := UnmarshalUserAccount(bytes)
//    bytes, err = userAccount.Marshal()

package types

import "encoding/json"

func (r *Deck) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Deck) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Deck) GetID() int {
	return r.ID
}

func (r *En) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *En) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *En) GetID() int {
	return r.ID
}

func (r *Expression) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Expression) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Expression) GetID() int {
	return r.ID
}

func (r *Filter) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Filter) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Filter) GetID() int {
	return r.ID
}

func (r *Lang) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Lang) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Lang) GetID() int {
	return r.ID
}

func (r *Tag) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Tag) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Tag) GetID() int {
	return r.ID
}

func (r *TopWordsByGoogle) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TopWordsByGoogle) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *TopWordsByGoogle) GetID() int {
	return r.ID
}

func (r *Tr) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Tr) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Tr) GetID() int {
	return r.ID
}

func (r *Translation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Translation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Translation) GetID() int {
	return r.ID
}

func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *User) GetID() int {
	return r.ID
}

func (r *UserAccount) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UserAccount) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *UserAccount) GetID() int {
	return r.ID
}

type En struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	PartsOfSpeech []string `json:"partsOfSpeech"`
	Level         string   `json:"level"`
	Frequency     int      `json:"frequency"`
	Examples      []string `json:"examples"`
	Transcription string   `json:"transcription"`
	Meanings      []string `json:"meanings"`
}

type Expression struct {
	ID  int    `json:"id"`
	Exp string `json:"exp"`
}

type Filter struct {
	ID     int   `json:"id"`
	DeckID int   `json:"deck_id"`
	Type   []int `json:"type"`
}

type Lang struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Tag struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

type TopWordsByGoogle struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

type Tr struct {
	ID               int                  `json:"id"`
	NormalizedSource string               `json:"normalizedSource"`
	DisplaySource    string               `json:"displaySource"`
	Translations     []TranslationElement `json:"translations"`
}

type TranslationElement struct {
	NormalizedTarget string            `json:"normalizedTarget"`
	DisplayTarget    string            `json:"displayTarget"`
	PosTag           string            `json:"posTag"`
	Confidence       float64           `json:"confidence"`
	PrefixWord       string            `json:"prefixWord"`
	BackTranslations []BackTranslation `json:"backTranslations"`
}

type BackTranslation struct {
	NormalizedText string `json:"normalizedText"`
	DisplayText    string `json:"displayText"`
	NumExamples    int    `json:"numExamples"`
	FrequencyCount int    `json:"frequencyCount"`
}

type Translation struct {
	ID     int    `json:"id"`
	TagID  int    `json:"tag_id"`
	LangID int    `json:"lang_id"`
	Tr     string `json:"tr"`
}

type User struct {
	ID      int    `json:"id"`
	UUID    string `json:"uuid"`
	Created int    `json:"created"`
	Updated int    `json:"updated"`
	Deleted int    `json:"deleted"`
}

type UserAccount struct {
	ID        int        `json:"id"`
	Decks     []Deck     `json:"decks"`
	StatItems []StatItem `json:"stat_items"`
}

type Deck struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StatItem struct {
	DeckID int   `json:"deck_id"`
	CardID int   `json:"card_id"`
	Right  bool  `json:"right"`
	Date   int   `json:"date"`
	Type   []int `json:"type"`
}
