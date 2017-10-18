package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/engine/docker"
	"github.com/drone/drone-runtime/engine/plugin"
	"github.com/drone/drone-runtime/runtime"
	"github.com/drone/drone-runtime/runtime/chroot"
	"github.com/drone/drone-runtime/runtime/term"
	"github.com/drone/drone-runtime/version"
	"github.com/drone/signal"
)

var tty = isatty.IsTerminal(os.Stdout.Fd())

func main() {
	b := flag.String("chroot", "", "")
	p := flag.String("plugin", "", "")
	t := flag.Duration("timeout", time.Hour, "")
	v := flag.Bool("version", false, "")
	h := flag.Bool("help", false, "")

	flag.BoolVar(h, "h", false, "")
	flag.BoolVar(v, "v", false, "")
	flag.Usage = usage
	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if *v {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	var source string
	if flag.NArg() > 0 {
		source = flag.Args()[0]
	}

	config, err := engine.ParseFile(source)
	if err != nil {
		log.Fatalln(err)
	}

	var engine engine.Engine
	if *p == "" {
		engine, err = docker.NewEnv()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		engine, err = plugin.Open(*p)
		if err != nil {
			log.Fatalln(err)
		}
	}

	hooks := &runtime.Hook{}
	hooks.GotLine = term.WriteLine(os.Stdout)
	if tty {
		hooks.GotLine = term.WriteLinePretty(os.Stdout)
	}

	var fs runtime.FileSystem
	if *b != "" {
		fs, err = chroot.New(*b)
		if err != nil {
			log.Fatalln(err)
		}
	}

	r := runtime.New(
		runtime.WithFileSystem(fs),
		runtime.WithEngine(engine),
		runtime.WithConfig(config),
		runtime.WithHooks(hooks),
	)

	ctx, cancel := context.WithTimeout(context.Background(), *t)
	ctx = signal.WithContext(ctx)
	defer cancel()

	err = r.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Println(`Usage: drone-runtime [OPTION]... [SOURCE]
      --plugin    loads a runtime engine from a .so file
      --timeout   sets an execution timeout
  -v, --version   display the version exit
  -h, --help      display this help and exit`)
}
