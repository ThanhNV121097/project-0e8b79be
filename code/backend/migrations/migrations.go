package migrations

import "embed"

// Files contains all upward SQL migrations.
//
//go:embed *.up.sql
var Files embed.FS
