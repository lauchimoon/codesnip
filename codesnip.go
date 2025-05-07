package main

import (
    "bufio"
    "os"
    "fmt"
    "strconv"
    "strings"

    "image"
    "image/color"
    "image/draw"
    "image/png"

    "golang.org/x/image/math/fixed"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"

    "github.com/lauchimoon/codesnip/lexer"
)

const (
    FONT_PATH = "./resources/fonts/SourceCodePro-Regular.ttf"
    DPI = 120.0
    PROGRAM_NAME = "codesnip"
    OUTPUT_FILE = "snippet.png"
)

type File struct {
    Path        string
    Content     []string
    NumLines    int
}

type Canvas struct {
    FontSize        int
    Bounds          image.Rectangle
    Handle          *image.RGBA
    Font            *truetype.Font
    FreetypeContext *freetype.Context
}

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("usage: %s <file> [num1-num2]\n", PROGRAM_NAME)
        fmt.Println("  if no range is given, screenshot the whole file")
        return
    }

    path := os.Args[1]
    file, err := ReadFile(path)
    if err != nil {
        panic(err)
    }

    num1, num2 := 0, 0
    if len(os.Args) == 3 {
        num1, num2 = SplitRange(os.Args[2])
        if !ValidRange(num1, num2, file) {
            fmt.Printf("error: range must be between 1 and %d\n", file.NumLines)
        }
    } else {
        num1 = 1
        num2 = file.NumLines
    }

    snippet := file.Content[num1-1:num2]
    longestLine := GetLongestLine(snippet)

    lex := lexer.LexerNew(strings.Join(snippet, ""))
    tokens := lex.Lex()

    canvas, err := CreateCanvas(24, longestLine, len(snippet))
    if err != nil {
        panic(err)
    }

    lineNum := 0
    textX := 5

    for _, token := range tokens {
        textY := 20 + lineNum*23
        textColor := TextColorByToken(token.Kind)
        canvas.DrawText(token.Text, textX, textY, textColor)

        addedWidth := MeasureString(canvas.Font, canvas.FontSize, token.Text)
        textX += addedWidth

        if token.Kind == lexer.TOKEN_NEWLINE {
            textX = 5
            lineNum += 1
        }
    }
    canvas.Export()
}

func SplitRange(s string) (int, int) {
    split := strings.Split(s, "-")
    num1, err := strconv.Atoi(split[0])
    if err != nil {
        return -1, -1
    }
    num2, err := strconv.Atoi(split[1])
    if err != nil {
        return -1, -1
    }

    return num1, num2
}

func ValidRange(num1, num2 int, f File) bool {
    return 1 <= num1 &&
         num1 <= f.NumLines &&
         1 <= num2 &&
         num2 <= f.NumLines
}

func ReadFile(path string) (File, error) {
    file := File{}

    handle, err := os.Open(path)
    if err != nil {
        return file, err
    }
    defer handle.Close()

    file.Path = path
    scanner := bufio.NewScanner(handle)
    for scanner.Scan() {
        text := scanner.Text()
        text = strings.ReplaceAll(text, "\t", "    ")
        text += "\n"
        file.Content = append(file.Content, text)
    }

    file.NumLines = len(file.Content)

    return file, nil
}

func GetLongestLine(lines []string) string {
    longest := ""
    maxLen := 0

    for _, line := range lines {
        lineLen := len(line)
        if maxLen < lineLen {
            maxLen = lineLen
            longest = line
        }
    }

    return longest
}

func CreateCanvas(fontSize int, longestLine string, numLines int) (Canvas, error) {
    canvas := Canvas{}

    canvas.FontSize = fontSize
    canvas.Bounds = image.Rect(0, 0, (fontSize - 8)*len(longestLine), fontSize*numLines)

    canvas.Handle = image.NewRGBA(canvas.Bounds)
    backgroundColor := color.RGBA{ 60, 60, 60, 255 }
    draw.Draw(canvas.Handle, canvas.Handle.Bounds(), &image.Uniform{backgroundColor},
            image.Point{0, 0}, draw.Src)

    font, err := LoadFont(FONT_PATH)
    if err != nil {
        return Canvas{}, err
    }

    canvas.Font = font

    canvas.FreetypeContext = freetype.NewContext()
    canvas.FreetypeContext.SetDPI(DPI)
    canvas.FreetypeContext.SetFont(font)
    canvas.FreetypeContext.SetClip(canvas.Bounds)
    canvas.FreetypeContext.SetDst(canvas.Handle)

    return canvas, nil
}

func TextColorByToken(tokenKind lexer.TokenKind) color.RGBA {
    switch (tokenKind) {
        case lexer.TOKEN_KEYWORD: return color.RGBA{ 255, 203, 0, 255 }
        case lexer.TOKEN_NUMBER: return color.RGBA{ 255, 105, 105, 255 }
        case lexer.TOKEN_STRING, lexer.TOKEN_CHAR: return color.RGBA{ 132, 255, 105, 255 }
        case lexer.TOKEN_PREPROC: return color.RGBA{ 144, 178, 255, 255 }
        case lexer.TOKEN_COMMENT: return color.RGBA{ 211, 176, 131, 255 }
        default: return color.RGBA{ 255, 255, 255, 255 }
    }
}

func MeasureString(font *truetype.Font, fontSize int, text string) int {
    width := 0

    for _, c := range text {
        index := font.Index(c)
        hmetric := font.HMetric(fixed.Int26_6(fontSize - 4), index)
        width += int(hmetric.AdvanceWidth)
    }

    return width
}

func (canvas Canvas) DrawText(text string, x, y int, clr color.RGBA) {
    canvas.FreetypeContext.SetSrc(&image.Uniform{clr})
    pt := freetype.Pt(x, y)

    if strings.Contains(text, "\n") {
        text = strings.ReplaceAll(text, "\n", "")
    }
    canvas.FreetypeContext.DrawString(text, pt)
}

func (canvas Canvas) Export() {
    outFile, err := os.Create(OUTPUT_FILE)
    if err != nil {
        panic(err)
    }

    img := canvas.Handle.SubImage(canvas.Handle.Bounds())
    png.Encode(outFile, img)
}

func LoadFont(path string) (*truetype.Font, error) {
    fontBytes, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    font, err := freetype.ParseFont(fontBytes)
    if err != nil {
        return nil, err
    }

    return font, nil
}
