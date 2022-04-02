package arweave

import "time"

type ArAccount string

const (
	MirrorUploader ArAccount = "Ky1c1Kkt-jZ9sY1hvLF5nCf6WWdBhIU5Un_BMYh-t3c"
)

const (
	DefaultCrawlStep     = 50
	DefaultFromHeight    = 1
	DefaultConfirmations = 10
	DefaultCrawlMinStep  = 10
)

var DefaultCrawlConfig = &crawlConfig{
	DefaultFromHeight,
	DefaultConfirmations,
	DefaultCrawlStep,
	DefaultCrawlMinStep,
	time.Duration(600),
}
