package icrawler

import (
	// "crawler/ucrawler"

	"io/ioutil"
)

func writeFile(path string, data []byte) (err error) {
	err = ioutil.WriteFile(path, data, 0644)
	return
}
