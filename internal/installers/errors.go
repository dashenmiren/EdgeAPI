// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package installers

type GrantError struct {
	err string
}

func newGrantError(err string) *GrantError {
	return &GrantError{err: err}
}

func (this *GrantError) Error() string {
	return this.err
}

func (this *GrantError) String() string {
	return this.err
}

func IsGrantError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*GrantError)
	return ok
}
