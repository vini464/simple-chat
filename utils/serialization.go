package utils

import (
  "encoding/json"
)

func SerializeJson(data Message) ([]byte, error){
  serialized_data, err := json.Marshal(data)
  if (err != nil) {
    return make([]byte, 0), err
  }
  return serialized_data, nil
}

func DeserializeToJson(serialized []byte, data *Message)  error {
  err := json.Unmarshal(serialized, data)
  return err
}

