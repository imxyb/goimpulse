package sender

import (
	"sync"
	"time"

	"strconv"
	"goimpulse/conf"
)

type Sequence struct {
	typeName  string
	cache     chan int64
	lock      sync.Mutex
	curLastId int
}

var sequence *Sequence
var batch int = conf.Cfg.Batch
var i int

func NewSeq(typeName string) *Sequence {
	sequence = &Sequence{
		cache:    make(chan int64, batch),
		typeName: typeName,
	}

	lastId, _ := strconv.Atoi(GetLastId(typeName, "1"))
	for i = lastId; i < lastId+batch; i++ {
		sequence.cache <- int64(i)
	}
	sequence.curLastId = i - 1
	SaveLastId(typeName, strconv.Itoa(sequence.curLastId))

	go sequence.expand()

	return sequence
}

func (this *Sequence) GetId() interface{} {
	return <-this.cache
}

func (this *Sequence) expand() {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			if len(this.cache) <= (batch / 2) {
				for i = this.curLastId + 1; i <= this.curLastId+batch/2; i++ {
					this.cache <- int64(i)
				}
				this.curLastId = i - 1
				SaveLastId(this.typeName, strconv.Itoa(this.curLastId))
			}
		}
	}
}
