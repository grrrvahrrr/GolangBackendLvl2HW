package elastic

import (
	"GoBeLvl2/enteties"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type PostStorage struct {
	elastic Elastic
	timeout time.Duration
}

func NewPostStorage(elastic Elastic) (PostStorage, error) {
	return PostStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (p PostStorage) Insert(ctx context.Context, user enteties.User) error {
	bdy, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("insert: marshall: %w", err)
	}

	// res, err := p.elastic.client.Create()
	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: user.Login,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return enteties.ErrConflict
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (p PostStorage) FindOne(ctx context.Context, login string) (enteties.User, error) {
	// res, err := p.elastic.client.Get()
	req := esapi.GetRequest{
		Index:      p.elastic.alias,
		DocumentID: login,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return enteties.User{}, fmt.Errorf("find one: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return enteties.User{}, enteties.ErrNotFound
	}

	if res.IsError() {
		return enteties.User{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var (
		user enteties.User
		body document
	)
	body.Source = &user

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return enteties.User{}, fmt.Errorf("find one: decode: %w", err)
	}

	return user, nil
}

func (p PostStorage) SearchUser(ctx context.Context, key string, value string) (map[string]interface{}, error) {
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		},
	}
	err := json.NewEncoder(&buffer).Encode(query)
	if err != nil {
		return nil, err
	}
	response, err := p.elastic.client.Search(p.elastic.client.Search.WithIndex(p.elastic.index),
		p.elastic.client.Search.WithBody(&buffer))
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		users :=
			hit.(map[string]interface{})["_source"].(map[string]interface{})
		return users, nil
	}

	return nil, nil

}
