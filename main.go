package main
import (
	"os"
	"io/ioutil"
	"strings"
	"fmt"
)

type Position interface {
	X() uint
	Y() uint
}
type position struct {
	x, y uint
}
func (p *position) X() uint {
	return p.x
}
func (p *position) Y() uint {
	return p.y
}

type Map interface {
	Width() uint
	Height() uint
	Set(X, Y uint, id rune)
	Get(X, Y uint) rune
}

const EMPTY = '-'
const NEW_LINE = '\n'
type Features map[rune]Position
type matrix struct {
	width, height uint
	features      Features
}
func (m *matrix) Width() uint {
	return m.width
}
func (m *matrix) Height() uint {
	return m.height
}
func (m *matrix) Get(x, y uint) rune {
	if m.features != nil {
		for mark, pos := range m.features {
			if pos.X() == x && pos.Y() == y {
				return mark;
			}
		}
	}
	return EMPTY
}
func (m *matrix) Set(x, y uint, value rune) {
	if m.features == nil {
		m.features = make(Features)
	}
	m.features[value] = &position{x, y}
}


func NewMap(width, height uint) *matrix {
	if width == 0 || height == 0 {
		width = 0
		height = 0
	}
	return &matrix{width:width, height:height, features:make(Features)}
}

func Parse(input string) (*matrix, error) {
	features := Features{}
	var width, height uint

	if input != "" {
		lines := strings.Split(input, string(NEW_LINE))
		height = uint(len(lines))
		for y, line := range lines {
			if l := uint(len(line)); l == 0 {
				return nil, fmt.Errorf("Empty row at index: %d", y)
			} else {
				if width > 0 {
					if l != width {
						return nil, fmt.Errorf("Each line must have uniform length (expected: %d, current %d)", width, l)
					}
				} else {
					width = l
				}
			}
			for x, field := range line {
				if field != EMPTY {
					features[field] = &position{uint(x), uint(y)}
				}
			}
		}
	}
	result := NewMap(width, height)
	result.features = features
	return result, nil
}
func (m *matrix)String() string {
	var result []string
	var y, x uint
	if m.width > 0 {
		for y = 0; y < m.height; y++ {
			var line []byte
			for x = 0; x < m.width; x++ {
				line = append(line, byte(m.Get(x, y)))
			}
			result = append(result, string(line))
		}
	}
	return strings.Join(result, string(NEW_LINE))
}

func (m *matrix)Route(from, to rune) (string, error) {
	f, ok := m.features[from]
	if !ok {
		return "", fmt.Errorf("Feature not found: %s", from)
	}
	t, ok := m.features[to]
	if !ok {
		return "", fmt.Errorf("Feature not found: %s", to)
	}
	dx := int(t.X()) - int(f.X())
	dy := int(t.Y()) - int(f.Y())

	result := ""
	if dx > 0 {
		result += strings.Repeat("RIGHT\n", dx)
	}
	if dx < 0 {
		result += strings.Repeat("LEFT\n", -dx)
	}
	if dy > 0 {
		result += strings.Repeat("DOWN\n", dy)
	}
	if dy < 0 {
		result += strings.Repeat("UP\n", -dy)
	}

	return result, nil
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if (err != nil) {
		print(err.Error())
		return
	}
	lines := strings.Split(string(input), "\n")

	table, err := Parse(strings.Join(lines[1:], "\n"))
	if (err != nil) {
		print(err.Error())
		return
	}
	r, err := table.Route('m', 'p')
	if (err != nil) {
		print(err.Error())
		return
	}
	fmt.Println(r)
}
