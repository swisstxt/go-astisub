package astisub

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Bytes
var (
	BytesBOM           = []byte{239, 187, 191}
	bytesLineSeparator = []byte("\n")
	bytesSpace         = []byte(" ")
)

// Colors
var (
	ColorBlack   = &Color{}
	ColorBlue    = &Color{Blue: 255}
	ColorCyan    = &Color{Blue: 255, Green: 255}
	ColorGreen   = &Color{Green: 255}
	ColorMagenta = &Color{Blue: 255, Red: 255}
	ColorRed     = &Color{Red: 255}
	ColorYellow  = &Color{Green: 255, Red: 255}
	ColorWhite   = &Color{Blue: 255, Green: 255, Red: 255}
)

// Errors
var (
	ErrInvalidExtension   = errors.New("astisub: invalid extension")
	ErrNoSubtitlesToWrite = errors.New("astisub: no subtitles to write")
)

// Now allows testing functions using it
var Now = func() time.Time {
	return time.Now()
}

// Options represents open or write options
type Options struct {
	Filename string
	Teletext TeletextOptions
}

// Open opens a subtitle reader based on options
func Open(o Options) (s *Subtitles, err error) {
	// Open the file
	var f *os.File
	if f, err = os.Open(o.Filename); err != nil {
		err = errors.Wrapf(err, "astisub: opening %s failed", o.Filename)
		return
	}
	defer f.Close()

	// Parse the content
	switch filepath.Ext(o.Filename) {
	case ".srt":
		s, err = ReadFromSRT(f)
	case ".ssa", ".ass":
		s, err = ReadFromSSA(f)
	case ".stl":
		s, err = ReadFromSTL(f)
	case ".ts":
		s, err = ReadFromTeletext(f, o.Teletext)
	case ".ttml":
		s, err = ReadFromTTML(f)
	case ".vtt":
		s, err = ReadFromWebVTT(f)
	default:
		err = ErrInvalidExtension
	}
	return
}

// OpenFile opens a file regardless of other options
func OpenFile(filename string) (*Subtitles, error) {
	return Open(Options{Filename: filename})
}

// Subtitles represents an ordered list of items with formatting
type Subtitles struct {
	Items    []*Item
	Metadata *Metadata
	Regions  map[string]*Region
	Styles   map[string]*Style
}

// NewSubtitles creates new subtitles
func NewSubtitles() *Subtitles {
	return &Subtitles{
		Regions: make(map[string]*Region),
		Styles:  make(map[string]*Style),
	}
}

// Item represents a text to show between 2 time boundaries with formatting
type Item struct {
	Comments    []string
	EndAt       time.Duration
	InlineStyle *StyleAttributes
	Lines       []Line
	Region      *Region
	StartAt     time.Duration
	Style       *Style
}

// String implements the Stringer interface
func (i Item) String() string {
	var os []string
	for _, l := range i.Lines {
		os = append(os, l.String())
	}
	return strings.Join(os, " - ")
}

// Color represents a color
type Color struct {
	Alpha, Blue, Green, Red uint8
}

// newColorFromString builds a new color based on a string
func newColorFromString(s string, base int) (c *Color, err error) {
	var i int64
	if i, err = strconv.ParseInt(s, base, 64); err != nil {
		err = errors.Wrapf(err, "parsing int %s with base %d failed", s, base)
		return
	}
	c = &Color{
		Alpha: uint8(i>>24) & 0xff,
		Blue:  uint8(i>>16) & 0xff,
		Green: uint8(i>>8) & 0xff,
		Red:   uint8(i) & 0xff,
	}
	return
}

// String expresses the color as a string for a specific base
func (c *Color) String(base int) string {
	var i = uint32(c.Alpha)<<24 | uint32(c.Blue)<<16 | uint32(c.Green)<<8 | uint32(c.Red)
	if base == 16 {
		return fmt.Sprintf("%.8x", i)
	}
	return strconv.Itoa(int(i))
}

// StyleAttributes represents style attributes
// TODO Convert styles+inline styles form different formats as well
type StyleAttributes struct {
	SSAAlignment         *int
	SSAAlphaLevel        *float64
	SSAAngle             *float64 // degrees
	SSABackColour        *Color
	SSABold              *bool
	SSABorderStyle       *int
	SSAEffect            string
	SSAEncoding          *int
	SSAFontName          string
	SSAFontSize          *float64
	SSAItalic            *bool
	SSALayer             *int
	SSAMarginLeft        *int // pixels
	SSAMarginRight       *int // pixels
	SSAMarginVertical    *int // pixels
	SSAMarked            *bool
	SSAOutline           *int // pixels
	SSAOutlineColour     *Color
	SSAPrimaryColour     *Color
	SSAScaleX            *float64 // %
	SSAScaleY            *float64 // %
	SSASecondaryColour   *Color
	SSAShadow            *int // pixels
	SSASpacing           *int // pixels
	SSAStrikeout         *bool
	SSAUnderline         *bool
	TeletextColor        *Color
	TeletextDoubleHeight *bool
	TeletextDoubleSize   *bool
	TeletextDoubleWidth  *bool
	TeletextSpacesAfter  *int
	TeletextSpacesBefore *int
	TTMLBackgroundColor  string
	TTMLColor            string
	TTMLDirection        string
	TTMLDisplay          string
	TTMLDisplayAlign     string
	TTMLExtent           string
	TTMLFontFamily       string
	TTMLFontSize         string
	TTMLFontStyle        string
	TTMLFontWeight       string
	TTMLLineHeight       string
	TTMLOpacity          string
	TTMLOrigin           string
	TTMLOverflow         string
	TTMLPadding          string
	TTMLShowBackground   string
	TTMLTextAlign        string
	TTMLTextDecoration   string
	TTMLTextOutline      string
	TTMLUnicodeBidi      string
	TTMLVisibility       string
	TTMLWrapOption       string
	TTMLWritingMode      string
	TTMLZIndex           int
	WebVTTAlign          string
	WebVTTLine           string
	WebVTTLines          int
	WebVTTPosition       string
	WebVTTRegionAnchor   string
	WebVTTScroll         string
	WebVTTSize           string
	WebVTTVertical       string
	WebVTTViewportAnchor string
	WebVTTWidth          string
}

// Metadata represents metadata
type Metadata struct {
	Comments  []string
	Copyright string
	Framerate int
	Language  string
	Title     string
}

// Region represents a subtitle's region
type Region struct {
	ID          string
	InlineStyle *StyleAttributes
	Style       *Style
}

// Style represents a subtitle's style
type Style struct {
	ID          string
	InlineStyle *StyleAttributes
	Style       *Style
}

// Line represents a set of formatted line items
type Line struct {
	Items     []LineItem
	VoiceName string
}

// String implement the Stringer interface
func (l Line) String() string {
	var texts []string
	for _, i := range l.Items {
		texts = append(texts, i.Text)
	}
	return strings.Join(texts, " ")
}

// LineItem represents a formatted line item
type LineItem struct {
	InlineStyle *StyleAttributes
	Style       *Style
	Text        string
}

// Add adds a duration to each time boundaries. As in the time package, duration can be negative.
func (s *Subtitles) Add(d time.Duration) {
	for _, v := range s.Items {
		v.EndAt += d
		v.StartAt += d
	}
}

// Duration returns the subtitles duration
func (s Subtitles) Duration() time.Duration {
	if len(s.Items) == 0 {
		return time.Duration(0)
	}
	return s.Items[len(s.Items)-1].EndAt
}

// ForceDuration updates the subtitles duration.
// If requested duration is bigger, then we create a dummy item.
// If requested duration is smaller, then we remove useless items and we cut the last item or add a dummy item.
func (s *Subtitles) ForceDuration(d time.Duration) {
	// Requested duration is the same as the subtitles'one
	if s.Duration() == d {
		return
	}

	// Requested duration is bigger than subtitles'one
	if s.Duration() > d {
		// Find last item before input duration and update end at
		var lastIndex = -1
		for index, i := range s.Items {
			// Start at is bigger than input duration, we've found the last item
			if i.StartAt >= d {
				lastIndex = index
				break
			} else if i.EndAt > d {
				s.Items[index].EndAt = d
			}
		}

		// Last index has been found
		if lastIndex != -1 {
			s.Items = s.Items[:lastIndex]
		}
	}

	// Add dummy item with the minimum duration possible
	if s.Duration() < d {
		s.Items = append(s.Items, &Item{EndAt: d, Lines: []Line{{Items: []LineItem{{Text: "..."}}}}, StartAt: d - time.Millisecond})
	}
}

// Fragment fragments subtitles with a specific fragment duration
func (s *Subtitles) Fragment(f time.Duration) {
	// Nothing to fragment
	if len(s.Items) == 0 {
		return
	}

	// Here we want to simulate fragments of duration f until there are no subtitles left in that period of time
	var fragmentStartAt, fragmentEndAt = time.Duration(0), f
	for fragmentStartAt < s.Items[len(s.Items)-1].EndAt {
		// We loop through subtitles and process the ones that either contain the fragment start at,
		// or contain the fragment end at
		//
		// It's useless processing subtitles contained between fragment start at and end at
		//             |____________________|             <- subtitle
		//           |                        |
		//   fragment start at        fragment end at
		for i, sub := range s.Items {
			// Init
			var newSub = &Item{}
			*newSub = *sub

			// A switch is more readable here
			switch {
			// Subtitle contains fragment start at
			// |____________________|                         <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentStartAt && sub.EndAt > fragmentStartAt:
				sub.StartAt = fragmentStartAt
				newSub.EndAt = fragmentStartAt
			// Subtitle contains fragment end at
			//                         |____________________| <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentEndAt && sub.EndAt > fragmentEndAt:
				sub.StartAt = fragmentEndAt
				newSub.EndAt = fragmentEndAt
			default:
				continue
			}

			// Insert new sub
			s.Items = append(s.Items[:i], append([]*Item{newSub}, s.Items[i:]...)...)
		}

		// Update fragments boundaries
		fragmentStartAt += f
		fragmentEndAt += f
	}

	// Order
	s.Order()
}

// IsEmpty returns whether the subtitles are empty
func (s Subtitles) IsEmpty() bool {
	return len(s.Items) == 0
}

// Merge merges subtitles i into subtitles
func (s *Subtitles) Merge(i *Subtitles) {
	// Append items
	s.Items = append(s.Items, i.Items...)
	s.Order()

	// Add regions
	for _, region := range i.Regions {
		if _, ok := s.Regions[region.ID]; !ok {
			s.Regions[region.ID] = region
		}
	}

	// Add styles
	for _, style := range i.Styles {
		if _, ok := s.Styles[style.ID]; !ok {
			s.Styles[style.ID] = style
		}
	}
}

// Order orders items
func (s *Subtitles) Order() {
	// Nothing to do if less than 1 element
	if len(s.Items) <= 1 {
		return
	}

	// Order
	var swapped = true
	for swapped {
		swapped = false
		for index := 1; index < len(s.Items); index++ {
			if s.Items[index-1].StartAt > s.Items[index].StartAt {
				var tmp = s.Items[index-1]
				s.Items[index-1] = s.Items[index]
				s.Items[index] = tmp
				swapped = true
			}
		}
	}
}

// Unfragment unfragments subtitles
func (s *Subtitles) Unfragment() {
	// Nothing to do if less than 1 element
	if len(s.Items) <= 1 {
		return
	}

	// Loop through items
	for i := 0; i < len(s.Items)-1; i++ {
		for j := i + 1; j < len(s.Items); j++ {
			// Items are the same
			if s.Items[i].String() == s.Items[j].String() && s.Items[i].EndAt == s.Items[j].StartAt {
				s.Items[i].EndAt = s.Items[j].EndAt
				s.Items = append(s.Items[:j], s.Items[j+1:]...)
				j--
			}
		}
	}

	// Order
	s.Order()
}

// Write writes subtitles to a file
func (s Subtitles) Write(dst string) (err error) {
	// Create the file
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		err = errors.Wrapf(err, "astisub: creating %s failed", dst)
		return
	}
	defer f.Close()

	// Write the content
	switch filepath.Ext(dst) {
	case ".srt":
		err = s.WriteToSRT(f)
	case ".ssa", ".ass":
		err = s.WriteToSSA(f)
	case ".stl":
		err = s.WriteToSTL(f)
	case ".ttml":
		err = s.WriteToTTML(f)
	case ".vtt":
		err = s.WriteToWebVTT(f)
	default:
		err = ErrInvalidExtension
	}
	return
}

// parseDuration parses a duration in "00:00:00.000", "00:00:00,000" or "0:00:00:00" format
func parseDuration(i, millisecondSep string, numberOfMillisecondDigits int) (o time.Duration, err error) {
	// Split milliseconds
	var parts = strings.Split(i, millisecondSep)
	var milliseconds int
	var s string
	if len(parts) >= 2 {
		// Invalid number of millisecond digits
		s = strings.TrimSpace(parts[len(parts)-1])
		if len(s) > 3 {
			err = fmt.Errorf("astisub: Invalid number of millisecond digits detected in %s", i)
			return
		}

		// Parse milliseconds
		if milliseconds, err = strconv.Atoi(s); err != nil {
			err = errors.Wrapf(err, "astisub: atoi of %s failed", s)
			return
		}
		milliseconds *= int(math.Pow10(numberOfMillisecondDigits - len(s)))
		s = strings.Join(parts[:len(parts)-1], millisecondSep)
	} else {
		s = i
	}

	// Split hours, minutes and seconds
	parts = strings.Split(strings.TrimSpace(s), ":")
	var partSeconds, partMinutes, partHours string
	if len(parts) == 2 {
		partSeconds = parts[1]
		partMinutes = parts[0]
	} else if len(parts) == 3 {
		partSeconds = parts[2]
		partMinutes = parts[1]
		partHours = parts[0]
	} else {
		err = fmt.Errorf("astisub: No hours, minutes or seconds detected in %s", i)
		return
	}

	// Parse seconds
	var seconds int
	s = strings.TrimSpace(partSeconds)
	if seconds, err = strconv.Atoi(s); err != nil {
		err = errors.Wrapf(err, "astisub: atoi of %s failed", s)
		return
	}

	// Parse minutes
	var minutes int
	s = strings.TrimSpace(partMinutes)
	if minutes, err = strconv.Atoi(s); err != nil {
		err = errors.Wrapf(err, "astisub: atoi of %s failed", s)
		return
	}

	// Parse hours
	var hours int
	if len(partHours) > 0 {
		s = strings.TrimSpace(partHours)
		if hours, err = strconv.Atoi(s); err != nil {
			err = errors.Wrapf(err, "astisub: atoi of %s failed", s)
			return
		}
	}

	// Generate output
	o = time.Duration(milliseconds)*time.Millisecond + time.Duration(seconds)*time.Second + time.Duration(minutes)*time.Minute + time.Duration(hours)*time.Hour
	return
}

// formatDuration formats a duration
func formatDuration(i time.Duration, millisecondSep string, numberOfMillisecondDigits int) (s string) {
	// Parse hours
	var hours = int(i / time.Hour)
	var n = i % time.Hour
	if hours < 10 {
		s += "0"
	}
	s += strconv.Itoa(hours) + ":"

	// Parse minutes
	var minutes = int(n / time.Minute)
	n = i % time.Minute
	if minutes < 10 {
		s += "0"
	}
	s += strconv.Itoa(minutes) + ":"

	// Parse seconds
	var seconds = int(n / time.Second)
	n = i % time.Second
	if seconds < 10 {
		s += "0"
	}
	s += strconv.Itoa(seconds) + millisecondSep

	// Parse milliseconds
	var milliseconds = float64(n/time.Millisecond) / float64(1000)
	s += fmt.Sprintf("%."+strconv.Itoa(numberOfMillisecondDigits)+"f", milliseconds)[2:]
	return
}

// appendStringToBytesWithNewLine adds a string to bytes then adds a new line
func appendStringToBytesWithNewLine(i []byte, s string) (o []byte) {
	o = append(i, []byte(s)...)
	o = append(o, bytesLineSeparator...)
	return
}
