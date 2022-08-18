package main

const (
	MediaType     = "MediaType"
	PHOTO         = "photography"
	ILLUSTRATIONS = "illustration"
	VECTORS       = "illustration&assetfiletype=eps"

	Orientations        = "Orientations"
	SQUARE              = "square"
	VERTICAL            = "vertical"
	HORIZONTAL          = "horizontal"
	PanoramicVertical   = "panoramicvertical"
	PanoramicHorizontal = "panoramichorizontal"

	NumberOfPeople = "NumberOfPeople"
	NoPeople       = "none"
	OnePerson      = "one"
	TwoPeople      = "two"
	GroupOfPeople  = "group"

	UNDEFINED = "undefined"
)

var (
	OptionalMediaType    = []string{PHOTO, ILLUSTRATIONS, VECTORS}
	OptionalOrientations = []string{SQUARE, VERTICAL, HORIZONTAL, PanoramicVertical, PanoramicHorizontal}
	OptionalNoPeople     = []string{NoPeople, OnePerson, TwoPeople, GroupOfPeople}

	queryMap = map[string][]string{
		MediaType:      OptionalMediaType,
		Orientations:   OptionalOrientations,
		NumberOfPeople: OptionalNoPeople,
	}

	queryDefault = map[string]string{
		MediaType:      PHOTO,
		Orientations:   SQUARE,
		NumberOfPeople: NoPeople,
	}
)

// RefactorInvalidQueryType automatically correct deviated parameters back to default values during parameter checking
func RefactorInvalidQueryType(queryType, query string) string {
	for _, val := range queryMap[queryType] {
		if val == query {
			return query
		}
	}
	return queryDefault[queryType]
}
