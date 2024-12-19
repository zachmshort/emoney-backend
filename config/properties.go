package config

type PropertyConfig struct {
	Name          string
	Group         string
	Position      int
	Price         int
	RentPrices    []int
	HouseCost     int
	MortgageValue int
	Images        []string
}

var DefaultProperties = []PropertyConfig{
	// Brown properties
	{
		Name:          "Mediterranean Avenue",
		Group:         "brown",
		Position:      1,
		Price:         60,
		RentPrices:    []int{2, 10, 30, 90, 160, 250},
		HouseCost:     50,
		MortgageValue: 30,
		Images:        []string{"mediterranean-avenue-front"},
	},
	{
		Name:          "Baltic Avenue",
		Group:         "brown",
		Position:      3,
		Price:         60,
		RentPrices:    []int{4, 20, 60, 180, 320, 450},
		HouseCost:     50,
		MortgageValue: 30,
		Images:        []string{"baltic-avenue-front"},
	},

	// Light Blue properties
	{
		Name:          "Oriental Avenue",
		Group:         "light-blue",
		Position:      6,
		Price:         100,
		RentPrices:    []int{6, 30, 90, 270, 400, 550},
		HouseCost:     50,
		MortgageValue: 50,
		Images:        []string{"oriental-avenue-front"},
	},
	{
		Name:          "Vermont Avenue",
		Group:         "light-blue",
		Position:      8,
		Price:         100,
		RentPrices:    []int{6, 30, 90, 270, 400, 550},
		HouseCost:     50,
		MortgageValue: 50,
		Images:        []string{"vermont-avenue-front"},
	},
	{
		Name:          "Connecticut Avenue",
		Group:         "light-blue",
		Position:      9,
		Price:         120,
		RentPrices:    []int{8, 40, 100, 300, 450, 600},
		HouseCost:     50,
		MortgageValue: 60,
		Images:        []string{"connecticut-avenue-front"},
	},

	// Pink properties
	{
		Name:          "St. Charles Place",
		Group:         "pink",
		Position:      11,
		Price:         140,
		RentPrices:    []int{10, 50, 150, 450, 625, 750},
		HouseCost:     100,
		MortgageValue: 70,
		Images:        []string{"st-charles-place-front"},
	},
	{
		Name:          "States Avenue",
		Group:         "pink",
		Position:      13,
		Price:         140,
		RentPrices:    []int{10, 50, 150, 450, 625, 750},
		HouseCost:     100,
		MortgageValue: 70,
		Images:        []string{"states-avenue-front"},
	},
	{
		Name:          "Virginia Avenue",
		Group:         "pink",
		Position:      14,
		Price:         160,
		RentPrices:    []int{12, 60, 180, 500, 700, 900},
		HouseCost:     100,
		MortgageValue: 80,
		Images:        []string{"virginia-avenue-front"},
	},

	// Orange properties
	{
		Name:          "St. James Place",
		Group:         "orange",
		Position:      16,
		Price:         180,
		RentPrices:    []int{14, 70, 200, 550, 750, 950},
		HouseCost:     100,
		MortgageValue: 90,
		Images:        []string{"st-james-place-front"},
	},
	{
		Name:          "Tennessee Avenue",
		Group:         "orange",
		Position:      18,
		Price:         180,
		RentPrices:    []int{14, 70, 200, 550, 750, 950},
		HouseCost:     100,
		MortgageValue: 90,
		Images:        []string{"tennessee-avenue-front"},
	},
	{
		Name:          "New York Avenue",
		Group:         "orange",
		Position:      19,
		Price:         200,
		RentPrices:    []int{16, 80, 220, 600, 800, 1000},
		HouseCost:     100,
		MortgageValue: 100,
		Images:        []string{"new-york-avenue-front"},
	},

	// Red properties
	{
		Name:          "Kentucky Avenue",
		Group:         "red",
		Position:      21,
		Price:         220,
		RentPrices:    []int{18, 90, 250, 700, 875, 1050},
		HouseCost:     150,
		MortgageValue: 110,
		Images:        []string{"kentucky-avenue-front"},
	},
	{
		Name:          "Indiana Avenue",
		Group:         "red",
		Position:      23,
		Price:         220,
		RentPrices:    []int{18, 90, 250, 700, 875, 1050},
		HouseCost:     150,
		MortgageValue: 110,
		Images:        []string{"indiana-avenue-front"},
	},
	{
		Name:          "Illinois Avenue",
		Group:         "red",
		Position:      24,
		Price:         240,
		RentPrices:    []int{20, 100, 300, 750, 925, 1100},
		HouseCost:     150,
		MortgageValue: 120,
		Images:        []string{"illinois-avenue-front"},
	},

	// Yellow properties
	{
		Name:          "Atlantic Avenue",
		Group:         "yellow",
		Position:      26,
		Price:         260,
		RentPrices:    []int{22, 110, 330, 800, 975, 1150},
		HouseCost:     150,
		MortgageValue: 130,
		Images:        []string{"atlantic-avenue-front"},
	},
	{
		Name:          "Ventnor Avenue",
		Group:         "yellow",
		Position:      27,
		Price:         260,
		RentPrices:    []int{22, 110, 330, 800, 975, 1150},
		HouseCost:     150,
		MortgageValue: 130,
		Images:        []string{"ventnor-avenue-front"},
	},
	{
		Name:          "Marvin Gardens",
		Group:         "yellow",
		Position:      29,
		Price:         280,
		RentPrices:    []int{24, 120, 360, 850, 1025, 1200},
		HouseCost:     150,
		MortgageValue: 140,
		Images:        []string{"marvin-gardens-front"},
	},

	// Green properties
	{
		Name:          "Pacific Avenue",
		Group:         "green",
		Position:      31,
		Price:         300,
		RentPrices:    []int{26, 130, 390, 900, 1100, 1275},
		HouseCost:     200,
		MortgageValue: 150,
		Images:        []string{"pacific-avenue-front"},
	},
	{
		Name:          "North Carolina Avenue",
		Group:         "green",
		Position:      32,
		Price:         300,
		RentPrices:    []int{26, 130, 390, 900, 1100, 1275},
		HouseCost:     200,
		MortgageValue: 150,
		Images:        []string{"north-carolina-avenue-front"},
	},
	{
		Name:          "Pennsylvania Avenue",
		Group:         "green",
		Position:      34,
		Price:         320,
		RentPrices:    []int{28, 150, 450, 1000, 1200, 1400},
		HouseCost:     200,
		MortgageValue: 160,
		Images:        []string{"pennsylvania-avenue-front"},
	},

	// Dark Blue properties
	{
		Name:          "Park Place",
		Group:         "dark-blue",
		Position:      37,
		Price:         350,
		RentPrices:    []int{35, 175, 500, 1100, 1300, 1500},
		HouseCost:     200,
		MortgageValue: 175,
		Images:        []string{"park-place-front"},
	},
	{
		Name:          "Boardwalk",
		Group:         "dark-blue",
		Position:      39,
		Price:         400,
		RentPrices:    []int{50, 200, 600, 1400, 1700, 2000},
		HouseCost:     200,
		MortgageValue: 200,
		Images:        []string{"boardwalk-front"},
	},
	{
		Name:          "Reading Railroad",
		Group:         "railroad",
		Position:      5,
		Price:         200,
		RentPrices:    []int{25, 50, 100, 200},
		MortgageValue: 100,
		Images:        []string{"reading-railroad-front"},
	},
	{
		Name:          "Pennsylvania Railroad",
		Group:         "railroad",
		Position:      15,
		Price:         200,
		RentPrices:    []int{25, 50, 100, 200},
		MortgageValue: 100,
		Images:        []string{"pennsylvania-railroad-front"},
	},
	{
		Name:          "B. & O. Railroad",
		Group:         "railroad",
		Position:      25,
		Price:         200,
		RentPrices:    []int{25, 50, 100, 200},
		MortgageValue: 100,
		Images:        []string{"b-o-railroad-front"},
	},
	{
		Name:          "Short Line",
		Group:         "railroad",
		Position:      35,
		Price:         200,
		RentPrices:    []int{25, 50, 100, 200},
		MortgageValue: 100,
		Images:        []string{"short-line-railroad-front"},
	},

	// Utilities
	{
		Name:          "Electric Company",
		Group:         "utility",
		Position:      12,
		Price:         150,
		RentPrices:    []int{4, 10}, // multiplier for dice roll (4x if 1 utility, 10x if both)
		MortgageValue: 75,
		Images:        []string{"electric-company-front"},
	},
	{
		Name:          "Water Works",
		Group:         "utility",
		Position:      28,
		Price:         150,
		RentPrices:    []int{4, 10}, // Multiplier for dice roll (4x if 1 utility, 10x if both)
		MortgageValue: 75,
		Images:        []string{"water-works-front"},
	},
}
