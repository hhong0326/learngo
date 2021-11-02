package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string

var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errExists     = errors.New("Exists")
)

// Search
func (d Dictionary) Search(word string) (string, error) {

	v, ok := d[word]

	if ok {
		return v, nil
	}

	return "", errNotFound

}

// Add
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errExists
	}

	return nil
}

// Update
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}

	return nil
}

// Delete
func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errNotFound
	}
	delete(d, word)
	return nil
}
