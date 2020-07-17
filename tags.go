package cloud66

import (
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
