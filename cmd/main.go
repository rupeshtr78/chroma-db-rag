package main

import (
	"chroma-db/cmd/db"
	"context"
)

func main() {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	db.RunVectorDb(ctx)

	// chat.ChatOllama(ctx)
}
