package cedar

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

type item struct {
	key   []byte
	value int
}

var (
	dict []item

	trie = New()
)

func loadDict() {
	testFile := "testdata/dict.txt"
	f, err := os.Open(testFile)
	if err != nil {
		log.Fatal("failed to open ", testFile)
	}
	defer f.Close()

	in := bufio.NewReader(f)
	added := make(map[string]struct{})

	var (
		key  string
		freq int
		pos  string
	)

	for {
		_, err := fmt.Fscanln(in, &key, &freq, &pos)
		if err == io.EOF {
			break
		}

		if _, ok := added[string(key)]; !ok {
			dict = append(dict, item{[]byte(key), freq})
			added[string(key)] = struct{}{}
		}
	}
}

func exist(i int) {
	item := dict[i]
	// fmt.Println(i, string(item.key))
	id, err := trie.Jump(item.key, 0)
	failIfError(err)

	key, err := trie.Key(id)
	failIfError(err)

	value, err := trie.Value(id)
	failIfError(err)

	if string(key) != string(item.key) || value != item.value {
		v, _ := trie.Get(item.key)
		fmt.Println("exist but no equal: ", i, string(key), string(item.key), value, item.value, v)
		panic("large dict test fail: no equal")
	}
}

func notExist(i int) {
	_, err := trie.Get(dict[i].key)
	// fmt.Println(i, err)
	if err != ErrNoPath && err != ErrNoValue {
		panic("large dict test fail: should not exist")
	}
}

func checkSize(exp int) {
	if keys, _, _, _ := trie.Status(); keys != exp {
		panic("not correct status")
	}
}

func insertDict(size int) {
	// Insert the first half of the dict.
	for i := 0; i < size/2; i++ {
		item := dict[i]
		if i%2 == 0 {
			if err := trie.Insert(item.key, item.value); err != nil {
				panic(err)
			}
		} else {
			if err := trie.Update(item.key, item.value); err != nil {
				panic(err)
			}
		}
	}
	checkSize(size / 2)
}

func checkDict(size int) {
	// Check the first half of the dict.
	for i := 0; i < size/2; i++ {
		exist(i)
	}
	log.Println("first half OK")

	// Delete even items in the first half.
	for i := 0; i < size/2; i += 2 {
		err := trie.Delete(dict[i].key)
		failIfError(err)
	}
	checkSize(size / 2 / 2)

	// Make sure even items were deleted, and the rest are fine.
	for i := 0; i < size/2; i++ {
		if i%2 == 0 {
			notExist(i)
		} else {
			exist(i)
		}
	}
	log.Println("first half odd OK")

	// Insert the second half of the dict.
	for i := size / 2; i < size; i++ {
		item := dict[i]
		trie.Insert(item.key, item.value)
	}
	checkSize(size/2/2 + (size - size/2))
}

func odd(size int) {
	for i := 0; i < size/2; i++ {
		if i%2 == 0 {
			notExist(i)
		} else {
			exist(i)
		}
	}
	log.Println("first half odd still OK")

	// Delete even items in the second half.
	half := size / 2
	if half%2 == 1 {
		half++
	}

	for i := half; i < size; i += 2 {
		err := trie.Delete(dict[i].key)
		failIfError(err)
	}

	// Make sure even items were deleted, and the rest are fine.
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			notExist(i)
		} else {
			exist(i)
		}
	}

	log.Println("odd OK")
}

func even(size int) {
	// Insert all even terms.
	for i := 0; i < size; i += 2 {
		item := dict[i]
		notExist(i)
		trie.Update([]byte(item.key), item.value)
	}

	for i := 0; i < size; i++ {
		exist(i)
	}
	log.Println("all OK")

	// Insert every item again, should be fine.
	for i := 1; i < size; i++ {
		item := dict[i]
		trie.Insert([]byte(item.key), item.value)
	}

	for i := 1; i < size; i += 2 {
		exist(i)
	}
	log.Println("still OK")
}

func TestLargeDict(t *testing.T) {
	loadDict()
	size := len(dict)
	log.Println("dict size:", size)

	insertDict(size)
	checkDict(size)

	odd(size)
	even(size)
}

func failIfError(err error) {
	if err != nil {
		panic(err)
	}
}
