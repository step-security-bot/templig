package templig

import "fmt"

// wrapError returns a formatted error with the given text if an error is provided, otherwise it returns nil.
func wrapError(text string, err error) error {
	if len(text) == 0 {
		text = "%w"
	}

	if err != nil {
		return fmt.Errorf(text, err) //nolint:err113
	}

	return nil
}
