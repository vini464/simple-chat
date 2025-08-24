package utils

import (
  "encoding/json"
)

type serializable interface {
  // TODO: place all message types here 
}

func SerializeJson[T serializable](data T) ([]byte, error){
  serialized_data, err := json.Marshal(data)
  if (err != nil) {
    return make([]byte, 0), err
  }
  return serialized_data, nil
}

func DeserializeToJson(serialized []byte, data *serializable)  error {
  err := json.Unmarshal(serialized, data)
  return err
}

