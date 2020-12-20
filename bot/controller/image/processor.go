package image

import (
	"io/ioutil"
	"math/rand"
	"strconv"
	"tbot/model"
	"time"
)

type OutcomeImageHandler interface {
	GetImage(imageSet model.ImageSet) []byte
}

type Processor struct {
	imageSets              chan model.ImageSet
	close                  chan struct{}
	outcomeImageSetHandler OutcomeImageHandler
}

func NewProcessor(close chan struct{}, handler OutcomeImageHandler) *Processor {
	p := &Processor{make(chan model.ImageSet, 10), close, handler}
	//TODO: config
	go p.extractImage()
	return p
}

func (p *Processor) Put(images model.ImageSet) {
	p.imageSets <- images
}

//TODO: revise solution. Consider a pool of routines to Put images
func (p *Processor) extractImage() {
	for {
		select {
		case imageSet := <-p.imageSets:
			image := p.outcomeImageSetHandler.GetImage(imageSet)
			//TODO: put image to queue
			ioutil.WriteFile(strconv.Itoa((int)(rand.NewSource(time.Now().UTC().UnixNano()).Int63()))+".png", image, 0644)

		case <-p.close:
			return
		}
	}
}
