package randgen

import (
	"fmt"
	"math/rand"
	"time"
)

var adjectives = []string{
	"Fast",
	"Slow",
	"Quick",
	"Speedy",
	"Rapid",
	"Swift",
	"Brisk",
	"Active",
	"Angry",
	"Big",
	"Brave",
	"Calm",
	"Clever",
	"Dark",
	"Easy",
	"Famous",
	"Gifted",
	"Holy",
	"Jolly",
	"Happy",
	"Kind",
	"Lively",
	"Nice",
	"Old",
	"Sad",
	"Tiny",
	"Wild",
	"Young",
	"Zany",
	"Adorable",
	"Beautiful",
	"Clean",
	"Drab",
	"Elegant",
	"Fancy",
	"Glamorous",
	"Handsome",
	"Long",
	"Magnificent",
	"Old-fashioned",
	"Plain",
	"Quaint",
	"Sparkling",
	"Ugliest",
	"Unsightly",
	"Alive",
	"Better",
	"Careful",
	"Clever",
	"Dead",
}

var nouns = []string{
	"Person",
	"Place",
	"Thing",
	"Animal",
	"Plant",
	"Food",
	"Mineral",
	"Car",
	"House",
	"Boat",
	"Plane",
	"Train",
	"Computer",
	"Phone",
	"Table",
	"Chair",
	"Desk",
	"Bed",
	"Shirt",
	"Shoes",
	"Jacket",
	"Coat",
	"Book",
	"Magazine",
	"Newspaper",
	"Letter",
	"Number",
	"Shape",
	"Color",
	"Sound",
	"Music",
}

func generateRandomUsername() string {
	rand.Seed(time.Now().UnixNano())
	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	number := rand.Intn(1000)
	return fmt.Sprintf("%s%s%d", adjective+"_", noun, number)
}
