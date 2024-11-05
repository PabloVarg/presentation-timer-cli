package cli

import "github.com/PabloVarg/presentation-timer-cli/internal/api"

type APIModel struct {
	api api.APIClient
}

type APIResponseModel[T any] struct {
	done bool
	err  error
	data T
}

func (m *APIResponseModel[T]) SetErr(err error) {
	m.done = true
	m.err = err
}

func (m *APIResponseModel[T]) SetData(data T) {
	m.done = true
	m.data = data
}

func (m APIResponseModel[T]) Loading() bool {
	return !m.done
}
