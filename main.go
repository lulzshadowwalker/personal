package main

import (
	"os"
	"context"
)

func main() {
	component := hello("lulzie")
	component.Render(context.Background(), os.Stdout)
}
