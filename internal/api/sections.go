package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/helpers"
)

type Section struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Position int           `json:"position"`
}

func GetSections(client APIClient, get KeyValueRetriever, presentationID int) ([]Section, error) {
	path, err := url.JoinPath(
		client.Url(get),
		fmt.Sprintf("presentations/%d/sections", presentationID),
	)
	if err != nil {
		return nil, err
	}
	path += "?sort_by=position&page_size=100"

	res, err := client.HTTPClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		response, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", string(response))
	}

	var result PaginatedResponse[[]Section]
	if err := helpers.ReadJSON(res.Body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func GetSection(client APIClient, get KeyValueRetriever, ID int) (Section, error) {
	path, err := url.JoinPath(client.Url(get), fmt.Sprintf("/sections/%d", ID))
	if err != nil {
		return Section{}, err
	}

	res, err := client.HTTPClient.Get(path)
	if err != nil {
		return Section{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		response, err := io.ReadAll(res.Body)
		if err != nil {
			return Section{}, err
		}

		return Section{}, fmt.Errorf("%s", string(response))
	}

	var result Section
	if err := helpers.ReadJSON(res.Body, &result); err != nil {
		return Section{}, err
	}

	return result, nil
}

type CreateSectionMsg struct {
	PresentationID int           `json:"-"`
	Name           string        `json:"name"`
	Duration       time.Duration `json:"duration"`
}

func CreateSection(client APIClient, get KeyValueRetriever, msg CreateSectionMsg) error {
	path, err := url.JoinPath(
		client.Url(get),
		fmt.Sprintf("/presentations/%d/sections", msg.PresentationID),
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res, err := client.HTTPClient.Post(path, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return ExtractErrorMsg(res.Body)
	}

	return nil
}

type EditSectionMsg struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

func UpdateSection(
	client APIClient,
	get KeyValueRetriever,
	p EditSectionMsg,
) error {
	path, err := url.JoinPath(client.Url(get), fmt.Sprintf("/sections/%d", p.ID))
	if err != nil {
		return err
	}

	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return ExtractErrorMsg(res.Body)
	}

	return nil
}

func DeleteSection(client APIClient, get KeyValueRetriever, ID int) error {
	path, err := url.JoinPath(client.Url(get), fmt.Sprintf("/sections/%d", ID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ExtractErrorMsg(res.Body)
	}

	return nil
}
