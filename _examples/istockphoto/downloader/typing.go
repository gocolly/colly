package downloader

const (
	nameMediaType      = "MediaType"
	nameOrientations   = "Orientations"
	nameNumberOfPeople = "NumberOfPeople"

	Photo         = "photography"
	Illustrations = "illustration"
	Vectors       = "illustration&assetfiletype=eps"

	Square              = "square"
	Vertical            = "vertical"
	Horizontal          = "horizontal"
	PanoramicVertical   = "panoramicvertical"
	PanoramicHorizontal = "panoramichorizontal"

	NoPeople      = "none"
	OnePerson     = "one"
	TwoPeople     = "two"
	GroupOfPeople = "group"

	UNDEFINED = "undefined"
)

type filterMediaType struct {
	Photo         string
	Illustrations string
	Vectors       string
	Undefined     string
	optional      []string
}

type filterOrientations struct {
	Square              string
	Vertical            string
	Horizontal          string
	PanoramicVertical   string
	PanoramicHorizontal string
	Undefined           string
	optional            []string
}

type filterNumberOfPeople struct {
	NoPeople      string
	OnePerson     string
	TwoPeople     string
	GroupOfPeople string
	Undefined     string
	optional      []string
}

var (
	MediaType = &filterMediaType{
		Photo:         Photo,
		Illustrations: Illustrations,
		Vectors:       Vectors,
		Undefined:     UNDEFINED,
		optional:      []string{Photo, Illustrations, Vectors, UNDEFINED},
	}
	Orientations = &filterOrientations{
		Square:              Square,
		Vertical:            Vertical,
		Horizontal:          Horizontal,
		PanoramicVertical:   PanoramicVertical,
		PanoramicHorizontal: PanoramicHorizontal,
		Undefined:           UNDEFINED,
		optional:            []string{Square, Vertical, Horizontal, PanoramicVertical, PanoramicHorizontal, UNDEFINED},
	}
	NumberOfPeople = &filterNumberOfPeople{
		NoPeople:      NoPeople,
		OnePerson:     OnePerson,
		TwoPeople:     TwoPeople,
		GroupOfPeople: GroupOfPeople,
		Undefined:     UNDEFINED,
		optional:      []string{NoPeople, OnePerson, TwoPeople, GroupOfPeople, UNDEFINED},
	}
)

var (
	queryMap = map[string][]string{
		nameMediaType:      MediaType.optional,
		nameOrientations:   Orientations.optional,
		nameNumberOfPeople: NumberOfPeople.optional,
	}

	queryDefault = map[string]string{
		nameMediaType:      MediaType.Photo,
		nameOrientations:   Orientations.Square,
		nameNumberOfPeople: NumberOfPeople.NoPeople,
	}
)

// RefactorInvalidQueryType Automatically correct deviated parameters back
// to default values during parameter checking
func RefactorInvalidQueryType(queryType, query string) string {
	for _, val := range queryMap[queryType] {
		if val == query {
			return query
		}
	}
	return queryDefault[queryType]
}
