## Wordle Guesser

Do you like wordle, but you're stuck on which word to guess next? This tool can give you some options to narrow down your search for the right word.

### Prerequisites

[Go 1.17](https://go.dev/dl/) or greater.

### Usage

I make the wild assumption that if you use this tool, it will be only after you've made at least one guess. (I know, I know, "when you assume, you make ..."). This means you should have some information right away: namely, at least the letters to exclude.

So although I created the flags below, if you run anything without the `exclude` flag, you won't get much use out of the tool. Sorry.

#### Flags
Enter letters to:

* **exclude** (the wrong ones) <img width="100" alt="excluded letters" src="https://user-images.githubusercontent.com/305137/161393292-1331fa3c-cb76-447e-9766-2630137229cf.png">

`go run main.go -exclude=dieu`

* **include** (the right ones) <img width="20" alt="Screen Shot 2022-03-30 at 07 27 36" src="https://user-images.githubusercontent.com/305137/161393281-bc8dbc15-e94a-495e-8829-a012787083fc.png">

`go run main.go -include=a`

* **pattern** (when you're on the right track) <img width="100" alt="Screen Shot 2022-03-29 at 12 49 19" src="https://user-images.githubusercontent.com/305137/161393329-2782814b-1d1b-4bdb-b85f-34cfc44c0d63.png">

`go run main.go -pattern=sha--`

* **antipattern** (the pattern that definitely doesn't work) <img width="100" alt="Screen Shot 2022-04-02 at 10 22 10" src="https://user-images.githubusercontent.com/305137/161394146-dd025f5b-4a15-494c-9c08-336718cb2bd2.png">

`go run main.go -antipattern=a----`

#### Combining Flags
Of course, you can also combine the flags, e.g.: `go run main.go -pattern=sha-- -include=a -exclude=dieu -antipattern=a----`

### Installation

To install, you can:
* download the [latest release](https://github.com/jeveleth/wordle-helper/releases), or
* you can build the source code: e.g., `go build -o wh`.


