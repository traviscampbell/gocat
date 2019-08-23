// GoCat started hanging around those bad eggs from the chroma pkg who lead him down a not so dark or lonely road.
// He's now a full blown addict and can't do anything without getting high on Nip first, cause in his own words,
// 'it makes shit bearable...'
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
)

func init() {
	// overwrite the fallback/plaintext lexer & style
	// adds a couple rules to highlight potentially common things in hopes
	// it'll at least break up the monotony of boring-as-fuck plain text somewhat.
	lexers.Fallback = chroma.MustNewLexer(
		&chroma.Config{
			Name:      "fallback",
			Aliases:   []string{"text", "plain", "fallback"},
			Filenames: []string{"*"},
		},
		chroma.Rules{
			"root": {
				chroma.Rule{Pattern: `\n`, Type: chroma.Text, Mutator: nil},
				chroma.Rule{Pattern: `[^\S\n]+`, Type: chroma.Text, Mutator: nil},
				chroma.Rule{Pattern: `\\\n`, Type: chroma.Text, Mutator: nil},
				chroma.Rule{Pattern: `\\`, Type: chroma.Text, Mutator: nil},
				chroma.Rule{Pattern: `[()\[\]'",:;$/?{}]|\.{1,3}`, Type: chroma.Punctuation, Mutator: nil},
				chroma.Rule{Pattern: `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`, Type: chroma.NameVariable, Mutator: nil},
				chroma.Rule{Pattern: `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`, Type: chroma.Keyword, Mutator: nil},
				chroma.Rule{Pattern: `\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`, Type: chroma.NameAttribute, Mutator: nil},
				chroma.Rule{Pattern: `(?i)(?:[0-3]?\d(?:st|nd|rd|th)?\s+(?:of\s+)?(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)|(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)\s+[0-3]?\d(?:st|nd|rd|th)?)(?:\,)?\s*(?:\d{4})?|[0-3]?\d[-\./][0-3]?\d[-\./]\d{2,4}`, Type: chroma.LiteralDate, Mutator: nil},
				chroma.Rule{Pattern: `[-+/*%#@=<>&^|!\\~]`, Type: chroma.Operator, Mutator: nil},
				chroma.Rule{Pattern: `\b[\d][\d_-]*\b`, Type: chroma.Number, Mutator: nil},
				chroma.Rule{Pattern: `\b[A-Z]{2,}\b`, Type: chroma.NameConstant, Mutator: nil},
			},
		},
	)

	// shitty style fallback to accompany the shitty random file lexer above...
	styles.Fallback = chroma.MustNewStyle("fallback", chroma.StyleEntries{
		chroma.Keyword:       "#80ff00 bold", // link-ish
		chroma.NameConstant:  "#efefef bold", // ALL CAPS
		chroma.NameAttribute: "#0077ff bold", // phone#-ish
		chroma.NameVariable:  "#ff7700 bold", // email-ish
		chroma.Date:          "#e04014 bold",
		chroma.Number:        "#0099ff bold",
		chroma.Operator:      "#fff700 bold",
		chroma.Punctuation:   "#fff700 bold",
		chroma.Text:          "#bababa",
	})
}

func printBanner() {
	w := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	b := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	y := color.New(color.FgHiYellow, color.Bold).SprintFunc()

	// figlet font => fbr2____.flf -the day will come when I want to know that. and, on that day, im'ma be so fucking proud of myself!
	color.Set(color.FgHiYellow, color.Bold)
	fmt.Println()
	fmt.Printf(" %s         %s         \n", y("#####    #####    ######     #### ########"), b(",_---~~~~~----._"))
	fmt.Printf("%s    %s \n", y("###  ##  ###  ##  ###  ##    ##### #  ###"), b("_,,_,*^____      _____``*g*\"*,"))
	fmt.Printf("%s   %s\n", y("###      ###  ##  ###       ## ###    ###"), b("/ __/ /'     ^.  /      \\ ^@q  f"))
	fmt.Printf("%s  %s \n", y("### ###  ###  ##  ###      ##  ###    ###"), b("[  @f | @))    |  | @))   l  0 _/"))
	fmt.Printf("%s   %s\n", y("###  ##  ###  ##  ###     ########    ###"), b("\\`/   \\~____  /__ \\_____/     \\"))
	fmt.Printf("%s    %s  \n", y("###  ##  ###  ##  ###  ## ##   ###    ###"), b("|           _l__l_           I"))
	fmt.Printf(" %s    %s \n", y("#####    #####    #####  ##   ###    ###"), b("}          [______]           I"))
	fmt.Printf("%s %s %s   %s \n", w("==========================="), y("###"), w("=========="), b("]            | | |            |"))
	fmt.Printf("  %s This is your brain on Nip %s      %s \n\n", y("✧ﾟ･:*"), y("*:･ﾟ✧"), b("]             ~ ~             |"))
	fmt.Println()
	color.Unset()
}

func main() {
	l := flag.String("l", "", "Specific Lexer to use")
	s := flag.String("s", "native", "Specific Style to use")

	ll := flag.Bool("ll", false, "Print available lexers and exit")
	sl := flag.Bool("sl", false, "Print available styles and exit")

	debug := flag.Bool("debug", false, "Debug Mode (Prints filename and lexer/style used)")
	flag.Parse()

	if *ll {
		fmt.Printf("\n%s\n\n", strings.Join(lexers.Names(false), ","))
		os.Exit(0)
	}

	if *sl {
		fmt.Printf("\n%s\n\n", strings.Join(styles.Names(), ","))
		os.Exit(0)
	}

	// no args given and nothing on stdin
	if flag.NArg() < 1 && !isPiped() {
		printBanner()
		fmt.Printf("  Usage: gocat [Options] [File|-]...\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// typical day in the life of a gocat. get
	// dope, get friends, get high; rinse & repeat
	dayInTheLife := func(in string) {
		gocat := GoCat{}

		if err := gocat.scoreSomeNip(in); err != nil {
			log.Fatalln(err.Error())
		}

		gocat.gatherHomies(in, *l, *s)

		if err := gocat.getNippedAF(); err != nil {
			log.Println(err.Error())
		}

		if *debug {
			log.Printf("Lexer: %s Style: %s\n", gocat.lexy.Config().Name, gocat.steez.Name)
		}
	}

	if flag.NArg() < 1 {
		dayInTheLife("-")
	} else {
		for _, farg := range flag.Args() {
			dayInTheLife(farg)
		}
	}
}

// GoCat just wants to perceive things in a more colorful
// light, and his favorite way of making that happen is by
// getting nipped-out-of-his-mind with his homies.
type GoCat struct {
	nip []byte

	lexy    chroma.Lexer
	steez   *chroma.Style
	fmatter chroma.Formatter
}

// scoreSomeNip locally from a file-boy or have
// some piped-in asap if the streets are drying up
func (gc *GoCat) scoreSomeNip(farg string) error {
	var err error

	switch {
	case farg == "-":
		gc.nip, err = ioutil.ReadAll(os.Stdin)

	default:
		fpath, _ := filepath.Abs(farg)
		gc.nip, err = ioutil.ReadFile(fpath)
	}

	return err
}

// hit up lexy, steez, and fmatter to let em know bout
// the bomb-ass Nip you just scored and that Time of the
// Nip is fast approaching so they better come the fuck on..
func (gc *GoCat) gatherHomies(f, l, s string) {

	// attempt to see if teminal supports truecolour
	// Reference: https://gist.github.com/XVilka/8346728 
    if ct := os.Getenv("COLORTERM"); ct == "truecolour" || ct == "24bit" {
		gc.fmatter = formatters.Get("terminal16m")
		formatters.Fallback = formatters.Get("terminal256")
	} else {
		gc.fmatter = formatters.Get("terminal256")
		formatters.Fallback = formatters.Get("terminal")
	}

	// -----

	lexy := lexers.Get(l)
	if lexy == nil && !strings.HasSuffix(f, ".txt") {
		lexy = lexers.Match(f)
	}
	if lexy == nil {
		lexy = lexers.Analyse(string(gc.nip))
	}
	if lexy == nil {
		lexy = lexers.Fallback
	}

	gc.lexy = chroma.Coalesce(lexy)

	// -----

	steez := styles.Fallback
	if ss := styles.Get(s); ss != nil && lexy.Config().Name != "fallback" {
		steez = tweakStyleForTerm(ss)
	}

	gc.steez = steez
}

// Spark It Up. It's Time to get NIPPED AS FUCK!!! **insert filthy drop**
func (gc *GoCat) getNippedAF() error {
	it, err := gc.lexy.Tokenise(&chroma.TokeniseOptions{Nested: true, State: "root"}, string(gc.nip))
	if err != nil {
		return err
	}
	return gc.fmatter.Format(os.Stdout, gc.steez, it)
}

// T T Y you ain't got no alibi, you ugly. Hey, Yeah! You ugly!
// -----
// tweak the style so it looks better in a terminal.
// basically just unsetting the background, which
// ends up giving the background color of the terminal
// which is generally dark. So, then we look for any other types
// that have dark colors set, and kinda invert it or w/e.
func tweakStyleForTerm(s *chroma.Style) *chroma.Style {
	se := make(chroma.StyleEntries)
	for _, tt := range s.Types() {
		t := s.Get(tt)

		var (
			cb = t.Colour.Brightness()
			bb = t.Background.Brightness()
		)

		switch {
		case tt == chroma.Background && bb > 0.15:
			fallthrough

		case tt != chroma.Background && cb < 0.375:
			t.Colour = t.Colour.Brighten((1.0 - cb) * 0.625)
			fallthrough

		default:
			t.Background = 0
			se[tt] = t.String()
		}
	}
	return chroma.MustNewStyle(s.Name, se)
}

// (V) (°,,,,°) (V) - why not, stdin..?
func isPiped() bool {
	stat, err := os.Stdin.Stat()
	return err == nil && (stat.Mode()&os.ModeCharDevice) != os.ModeCharDevice
}

