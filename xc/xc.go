package xc

import (
	"os"
	"sync"
)

var Wg sync.WaitGroup
var FileXc []*os.File
var Mutex = &sync.Mutex{}
