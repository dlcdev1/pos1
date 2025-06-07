package main

import (
	"fmt"
	"github.com/dlcdev1/pos1/CURSO-GO/11-EVENTOS/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}
