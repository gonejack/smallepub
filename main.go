package main

import (
	"log"

	"smallepub/smallepub"
)

func main() {
	cmd := smallepub.SmallEpub{
		Options: smallepub.MustParseOptions(),
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
