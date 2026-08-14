package xcrypt

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strings"

	"github.com/terefang/gommons/pkg/xconvert"
	"github.com/terefang/gommons/pkg/xcrypto"
)

// PasswordComplexity
const (
	PasswordComplexityLevel0 = 1 // only uppercase
	PasswordComplexityLevel1 = 2 // adds lowercase
	PasswordComplexityLevel2 = 3 // adds digits
	PasswordComplexityLevel3 = 4 // adds compatible symbols
	PasswordComplexityLevel4 = 5 // adds more symbols

	PasswordComplexityDefault = PasswordComplexityLevel2
	PasswordComplexitySafer   = PasswordComplexityLevel3
	PasswordComplexityMax     = PasswordComplexityLevel4
)

// PasswordSymbolSets list of internal sets
var PasswordSymbolSets map[string][]string = make(map[string][]string)

var PasswordSymbolSetComplex = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "_!$/=#*+", "%§&?")
var PasswordSymbolSetComplex2 = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "-_!$%/=?*+", ",.;:§&#")
var PasswordSymbolSetMainframe = xconvert.AsArray("ABCDEFGHIJKLM", "02468", "NOPQRSTUVWXYZ", "13579")
var PasswordSymbolSetFlickrBase58 = xconvert.AsArray("ABCDEFGHJKLMNPQRSTUVWXYZ", "abcdefghijkmnopqrstuvwxyz", "123456789")
var PasswordSymbolSetCookieBase90 = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "!#$%&'()*+-./:<=>?@[]^_`{|}~")
var PasswordSymbolSetBase85 = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "@[\\]^_`!\"#$%&'()*+,-./:;{<|=}>~?")
var PasswordSymbolSetBase62 = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789")
var PasswordSymbolSetRfc4648Base64 = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "+/")
var PasswordSymbolSetRfc4648Base32 = xconvert.AsArray("ABCDEFGHIJKLM", "NOPQRSTUVWXYZ", "234567")
var PasswordSymbolSetRfc4648Base32SafeSet = xconvert.AsArray("ACDEFGHJKLM", "NPQRTUVWXY", "3467")
var PasswordSymbolSetWordSafe = xconvert.AsArray("RVWXcfghjmpqrvwx", "23456789", "CFGHJMPQ")
var PasswordSymbolSetSafer = xconvert.AsArray("ACDEFGHJKLMNPQRTUVWXY", "abcdefghjkmnpqrtuvwxy", "34679", "$%&*+")
var PasswordSymbolSetExtensive = xconvert.AsArray("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789", "!#%&()*+,-./:;<=>?@[\\]^_{|}~")

var PasswordSymbolSetDefault = PasswordSymbolSetCookieBase90

var WordSet0A = []string{"Colonel", "General", "Captain", "Doctor", "Professor", "Magister", "Lord", "Baron", "Count", "Marquis", "Sir", "Master", "Archon", "Esquire", "Reverend", "Monsieur"}

var WordSet0B = []string{"Lady", "Countess", "Baroness", "Marquess", "Madame", "Mistress", "Magistra", "Abbess", "Oracle", "Empress", "Priestess", "Archess", "Widow"}

var WordSet1A = []string{"Viktor", "Alar", "Mordechai", "Cassian", "Lucien", "Severin", "Octavian", "Aurelius", "Malachai"}

var WordSet1B = []string{"Lilith", "Selene", "Isolde", "Morgana", "Vivienne", "Lenore", "Cassandra", "Ophelia", "Elara", "Seraphine", "Belladonna", "Nyx", "Helena", "Viola", "Adrienne", "Noa", "Aya", "Zara", "Mira", "Nia", "Vesper", "Lux", "Nyra", "Sora", "Vaelis", "Xyra", "Morren", "Elyra", "Kaelis", "Vireya", "Sylvara", "Neris", "Zyrene", "Arissa"}

var WordSet2 = []string{"Acheron", "Albret", "Albret", "Angus", "Anjou", "Apremont", "Arcanis", "Armagnac", "Arnim", "Arundel", "Ash", "Ashcroft", "Ashen", "Atholl", "Audley", "Aurelian", "Azarin", "Bar", "Bardolf", "Baux", "Bearn", "Beauchamp", "Beaufort", "Beaujeu", "Beaumont", "Beaumont", "Belladore", "Bellgrave", "Berkeley", "Bethune", "Bigorre", "Bismarck", "Blackfang", "Blackthorn", "Blackveil", "Blackwood", "Blois", "Bloodmoon", "Bloodthorn", "Blucher", "Bohun", "Borcke", "Bourbon", "Branicki", "Braose", "Brienne", "Brienne", "Bruce", "Bulow", "Campbell", "Capet", "Carrion", "Carrouges", "Cassadine", "Chalon", "Chatillon", "Chodkiewicz", "Clare", "Clermont", "Clifford", "Clinton", "Comminges", "Comminges", "Corvin", "Coucy", "Courcy", "Courtenay", "Courtenay", "Craon", "Crow", "Crowe", "Crowhurst", "Czartoryski", "Dacre", "Dammartin", "Danton", "Darkwell", "Daubigne", "Delacroix", "Denhoff", "Despenser", "DeValoire", "Dohna", "Douglas", "Draven", "Drax", "Dray", "Drazan", "Dreadmoor", "Dreux", "Drucki", "Dusk", "Ebonveil", "Evernight", "Exoudun", "Fang", "Ferrers", "Fiennes", "Fife", "Finkenstein", "Firlej", "FitzAlan", "FitzGerald", "FitzWalter", "Flint", "Foix", "Fraser", "Frost", "Giffard", "Gneisenau", "Goltz", "Gordon", "Gore", "Graham", "Grahame", "Graves", "Grey", "Greystoke", "Grimm", "Grimshaw", "Guise", "Harcourt", "Harrow", "Hastings", "Hawthorne", "Hex", "Hollow", "Holloway", "Hungerford", "Ironclaw", "Issoudun", "Itzenplitz", "Jabłonowski", "Joinville", "Joinville", "Kaine", "Kaligor", "Kalinowski", "Kassar", "Khan", "Kharov", "Kiszka", "Kleist", "Koniecpolski", "Korvin", "KorwinKossakowski", "Krell", "Kuro", "Lacy", "LaGuiche", "LaMarche", "Latimer", "LaTremoille", "Laval", "LaVauguyon", "Laveau", "Leslie", "Limoges", "Lindsay", "Lisle", "Lorraine", "Lovell", "Lubomirski", "Lucy", "Lumley", "Lusignan", "Mabuse", "MacDonald", "MacDuff", "Maleficar", "Malovar", "Maltravers", "Maltzahn", "Manchu", "Mar", "Marceau", "Margrave", "Massalski", "McShay", "Mniszech", "Moltke", "Montacute", "Montagu", "Montaigu", "Montauban", "Montclair", "Montdidier", "Montfort", "Montfort", "Montmorency", "Moray", "Morcant", "Mori", "Moroz", "Morrow", "Mortimer", "Mortimer", "Moru", "Morvain", "Mowbray", "Nemours", "Nemours", "Netherby", "Neville", "Night", "Nightingale", "Nightshade", "Nocturne", "Noir", "Nyx", "Obsidian", "Oginski", "Olelkowicz", "Orlov", "Ossolinski", "Osten", "Ostrogski", "Pac", "Percy", "Perigord", "Perigord", "Plessen", "Polignac", "Poniatowski", "Pons", "Potocki", "Poynings", "Radziwill", "Ravager", "Ravane", "Ravelle", "Raven", "Ravenscroft", "Ravik", "Raze", "Rethel", "Rieux", "Riven", "Rochechouart", "Rohan", "Rook", "Roos", "Ros", "Rosencourt", "Ross", "Runevale", "Sable", "Sacken", "SaintPol", "Sancerre", "Sanguszko", "Sapieha", "Scales", "Scrope", "Scythe", "Segrave", "Semur", "Serpent", "Serpentis", "Severaine", "Shade", "Shadowend", "Shen", "Shenkar", "Shiro", "Sinclair", "Soissons", "Soryn", "Stafford", "Stargrave", "Stewart", "StJohn", "Strix", "Sully", "Talbot", "Talon", "Tancarville", "Tancarville", "Thorn", "Thornfang", "Thornfield", "Toulouse", "Tyszkiewicz", "Ufford", "Umbra", "Valash", "Vale", "Valecourt", "Valemont", "Valence", "Valois", "Valtieri", "Valtor", "Vane", "Vanehurst", "Varr", "Vaudrin", "Veil", "Vendome", "Venom", "Ventadour", "Ventadour", "Ventress", "Vere", "Vermandois", "Vesper", "Vey", "Veylan", "Veyron", "Vienne", "Viper", "Viremont", "Voidmere", "Volkov", "Voren", "Voss", "Warenne", "Wartberg", "Wartensleben", "Wedel", "Whitlock", "Willoughby", "Winter", "Winterbourne", "Winterfeldt", "Wisniowiecki", "Witzleben", "Wolfbane", "Wolfsbane", "Wraith", "Wrangel", "Wycliffe", "York", "Zamoyski", "Zarek", "Zarin", "Zaslawski", "Zbaraski", "Zhan", "Zieten", "Zouche"}
var WordSet3 = []string{"Crimson", "Jade", "Black", "Scarlet", "Ivory", "Silent", "Velvet", "White", "Veiled", "Midnight", "Red", "Orange", "Yellow", "Green", "Blue", "Indigo", "Violet", "Black", "White", "Grey", "Velvet", "Maroon", "Tan", "Brown", "Silver", "Golden", "Bronze", "Granite", "Gleaming", "Shining", "Dull", "Gilded", "Beaming", "Blazing", "Burning", "Bright", "Dim", "Dark", "Tall", "Short", "Fat", "Skinny", "Small", "Large", "Grand", "Lowly", "Gangly", "Thin", "Wide", "Drunken", "Sober", "Wet", "Dry", "High", "Low", "Royal", "Favoured", "Smart", "Foolish", "Tiring", "Full", "Empty", "New", "Old", "Brandished", "Polished", "Sparkling", "Matted", "Broken", "Mended", "Fallen", "Risen", "Lifted", "Emerald", "Saphire", "Diamond", "Pink", "Cerulean", "Lavish", "Glazed", "Flattened", "Tipped", "Sheltered", "Stray", "Lone", "Lonely", "Vile", "Lucky", "Cunning", "Sultry", "OneEyed", "Lazy", "Fair", "Bountiful", "Jolly", "Blind", "Hidden", "Frozen"}
var WordSet4 = []string{"Widow", "Empress", "Orchid", "Serpent", "Queen", "Lotus", "Fang", "Raven", "Oracle", "Duchess"}
var WordSet5 = "_!%/=?*+,."

var WordSet []string = make([]string, 0)

func init() {
	PasswordSymbolSets["Complex"] = PasswordSymbolSetComplex
	PasswordSymbolSets["Complex2"] = PasswordSymbolSetComplex2
	PasswordSymbolSets["Mainframe"] = PasswordSymbolSetMainframe
	PasswordSymbolSets["FlickrBase58"] = PasswordSymbolSetFlickrBase58
	PasswordSymbolSets["CookieBase90"] = PasswordSymbolSetCookieBase90
	PasswordSymbolSets["Base85"] = PasswordSymbolSetBase85
	PasswordSymbolSets["Base62"] = PasswordSymbolSetBase62
	PasswordSymbolSets["Rfc4648Base64"] = PasswordSymbolSetRfc4648Base64
	PasswordSymbolSets["Rfc4648Base32"] = PasswordSymbolSetRfc4648Base32
	PasswordSymbolSets["Rfc4648Base32SafeSet"] = PasswordSymbolSetRfc4648Base32SafeSet
	PasswordSymbolSets["WordSafe"] = PasswordSymbolSetWordSafe
	PasswordSymbolSets["Safer"] = PasswordSymbolSetSafer
	PasswordSymbolSets["Extensive"] = PasswordSymbolSetExtensive
	PasswordSymbolSets["Default"] = PasswordSymbolSetDefault
	WordSet = append(WordSet, WordSet0A...)
	WordSet = append(WordSet, WordSet1A...)
	WordSet = append(WordSet, WordSet0B...)
	WordSet = append(WordSet, WordSet1B...)
	WordSet = append(WordSet, WordSet2...)
	WordSet = append(WordSet, WordSet3...)
	WordSet = append(WordSet, WordSet4...)
}

func GenerateWordPass(length int) string {
	_sb := strings.Builder{}
	_s := rand.Intn(64)
	if false {
		if _s&0x1 == 1 {
			_s = rand.Intn(len(WordSet0A))
			_sb.WriteString(WordSet0A[_s])
			_s = rand.Intn(len(WordSet1A))
			_sb.WriteString(WordSet1A[_s])
		} else {
			_s = rand.Intn(len(WordSet0B))
			_sb.WriteString(WordSet0B[_s])
			_s = rand.Intn(len(WordSet1B))
			_sb.WriteString(WordSet1B[_s])
		}
		_s = rand.Intn(len(WordSet5))
		_sb.WriteByte(WordSet5[_s])
		_s = rand.Intn(len(WordSet2))
		_sb.WriteString(WordSet2[_s])
		_s = rand.Intn(len(WordSet5))
		_sb.WriteByte(WordSet5[_s])
		_s = rand.Intn(len(WordSet3))
		_sb.WriteString(WordSet3[_s])
		_s = rand.Intn(len(WordSet4))
		_sb.WriteString(WordSet4[_s])
		_s = rand.Intn(256)
		_sb.WriteString(fmt.Sprintf("%d", _s))
	}
	for len(_sb.String()) < length {
		_s = rand.Intn(len(WordSet))
		_sb.WriteString(WordSet[_s])
		_s = rand.Intn(10)
		_sb.WriteString(fmt.Sprintf("%d", _s))
		_s = rand.Intn(len(WordSet5))
		_sb.WriteByte(WordSet5[_s])
	}
	return _sb.String()
}

func GeneratePassword(length int) string {
	return GeneratePasswordWithSym(PasswordSymbolSetDefault, length)
}

func GeneratePasswordWithSym(symbolset []string, length int) string {
	return GeneratePasswordWithSeedSymLevel(-1, symbolset, length, len(symbolset))
}

func GeneratePasswordWithSymLevel(symbolset []string, length int, level int) string {
	return GeneratePasswordWithSeedSymLevel(-1, symbolset, length, level)
}

func GeneratePasswordWithSeedSymLevel(seed int64, symbolset []string, length int, level int) string {
	if seed == -1 {
		seed = rand.Int63()
	}
	return GeneratePasswordWithRngSymLevelAlgo(rand.New(rand.NewSource(seed)), symbolset, length, level, true)
}

func GeneratePasswordWithRngSymLevelAlgo(_rng *rand.Rand, symbolset []string, length int, level int, altAlgo bool) string {
	// Make a buffer and fill it with all characters
	buf := make([]byte, length)
	lensym := len(symbolset)
	lens := make([]int, lensym)
	for i, s := range symbolset {
		lens[i] = len(s)
	}
	lensym = min(lensym, level)
	if altAlgo {
		for i := 0; i < length; i++ {
			buf[i] = symbolset[i%lensym][int(_rng.Int63())%lens[i%lensym]]
		}
		// Shuffle
		_rng.Shuffle(len(buf)-lensym, func(i, j int) {
			buf[i+lensym], buf[j+lensym] = buf[j+lensym], buf[i+lensym]
		})
	} else {
		_sb := strings.Builder{}
		for i := 0; i < lensym; i++ {
			_sb.WriteString(symbolset[i])
		}
		sset := _sb.String()
		slen := len(sset)
		for i := 0; i < lensym; i++ {
			buf[i] = symbolset[i][int(_rng.Int63())%lens[i]]
		}
		for i := lensym; i < length; i++ {
			buf[i] = sset[int(_rng.Int63())%slen]
		}
		// Shuffle
		_rng.Shuffle(len(buf), func(i, j int) {
			buf[i], buf[j] = buf[j], buf[i]
		})
	}
	return string(buf)
}

func GeneratePasswordWithKdfSymLevel(_perm string, symbolset []string, length int, level int) string {
	_dk := xcrypto.Kdf(sha1.New, length, []byte(_perm))
	return GeneratePasswordWithPermSymLevel(_dk, symbolset, length, level)
}
func GeneratePasswordWithPermSymLevel(_perm []byte, symbolset []string, length int, level int) string {
	plen := len(_perm)
	buf := make([]byte, length)
	lensym := len(symbolset)
	lens := make([]int, lensym)
	for i, s := range symbolset {
		lens[i] = len(s)
	}
	lensym = min(lensym, level)

	for i := 0; i < length; i++ {
		buf[i] = symbolset[i%lensym][int(_perm[i%plen])%lens[i%lensym]]
	}
	// Shuffle
	ti := 0
	for i := lensym; i < length; i++ {
		ti += int(_perm[i%plen])
		tmp := buf[i]
		buf[i] = buf[ti%length]
		buf[ti%length] = tmp
	}
	return string(buf)
}
