# funalyser

> a command-line tool that analyses the time and space complexity of functions written in **Go** — giving you a quick glance into how heavy your code is

<br />

## 🪡 Simple Usage
type `funalyser analyse` and add route to the file you want to get an analysis of. That's it 🙌

gif

## 🧰 Functionality

- Recursive calls detection
- Loop-based iteration
- Memory allocation patterns (`make`, `append`, etc.)
- Fan-out factor (number of recursive calls per invocation)
- A basic complexity estimate (`O(n)`, `O(log n)`, `O(n log n)`...)

## ⚙️ Options

`funalyser` has flags:

- `--func` specify if you want an analysis for a specific function
- `--json` [in development] outputs the analysis in json format 

#### ⌨️ Usage:

- `funalyser analyse test/test_data/space_samples.go` 
- `funalyser analyse test/test_data/time_samples.go --func recursion`

### ⬇️ Download

`go install https://github.com/DanyloPiatyhorets/funalyser@latest`

## 🛠️ Supported Languages
- Goland

> Note: as part of the future roadmap, connecting parsers in other languages like
> - Java
> - Python
> - JavaScript / Typescript 


## 🎯 Personal Motivation

As part of our **Language Processors** module at City, University of London, I was already neck-deep building a parser in Java — working with ASTs, understanding function calls and how compilation process works. That coursework planted the seed

But I wanted to go further and try applying the same principles to a **real-world static analyser** and do it in a language I wasn’t already fluent in - Go!

So I picked **Go**, and along with it, **cobra** for building the CLI. The result? This project — a tool that reads Go (more options in the future 🚀) code, parses it into an AST, and inspects functions 
