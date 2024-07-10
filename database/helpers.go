package database

import "encoding/json"

func MarshalObject(data any, entityId string, entityType string) (*Entity, error) {

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	entity := Entity{
		Id:   entityId,
		Type: entityType,
		Data: string(b),
	}

	return &entity, nil

}

func UnmarshalObject(entity Entity, v any) error {

	err := json.Unmarshal([]byte(entity.Data), &v)
	if err != nil {
		return err
	}

	return nil
}
