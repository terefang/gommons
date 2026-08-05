package xcrypt

import (
    "crypto/sha1"
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
