package cloud66

import (
	"strconv"
	"time"
)

type EntityTags struct {
	EntityType string      `json:"entity_type"`
	EntityId   string      `json:"entity_id"`
	Tags       []EntityTag `json:"tags"`
	CreatedAt  time.Time   `json:"created_at_iso"`
	UpdatedAt  time.Time   `json:"updated_at_iso"`
}

type EntityTag struct {
	Value    string `json:"value"`
	Reserved bool   `json:"reserved"`
}

func (c *Client) AllTags() ([]EntityTags, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []EntityTags
	var intermediateResult []EntityTags

	for {
		req, err := c.NewRequest("GET", "/tags", nil, query_strings)
		if err != nil {
			return nil, err
		}

		intermediateResult = nil
		err = c.DoReq(req, &intermediateResult, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, intermediateResult...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}
	return result, nil
}

func (c *Client) EntityTags(entity, id string) (*EntityTags, error) {
	req, err := c.NewRequest("GET", "/tags/"+entity+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}

	var entityTags *EntityTags
	err = c.DoReq(req, &entityTags, nil)
	if err != nil {
		return nil, err
	}

	return entityTags, nil
}

func (c *Client) PatchEntityTags(entity, id string, tagsToAdd, tagsToDelete []string) (*EntityTags, error) {
	type tagOperation struct {
		Operation string   `json:"op"`
		Tags      []string `json:"tags"`
	}

	addOperation := tagOperation{
		Operation: "add",
		Tags:      tagsToAdd,
	}

	deleteOperation := tagOperation{
		Operation: "delete",
		Tags:      tagsToDelete,
	}

	operations := struct {
		Operations []tagOperation `json:"operations"`
	}{
		Operations: []tagOperation{addOperation, deleteOperation},
	}

	req, err := c.NewRequest("PATCH", "/tags/"+entity+"/"+id, operations, nil)
	if err != nil {
		return nil, err
	}
	var entityTags *EntityTags
	return entityTags, c.DoReq(req, &entityTags, nil)
}
