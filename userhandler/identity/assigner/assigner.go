package assigner

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adhityasan/ekyc-api/config"
	"github.com/adhityasan/ekyc-api/userhandler/identity"
)

// DukcapilSimulatorAssigner Assign data from dukcapil into local
func DukcapilSimulatorAssigner(nik string, identity *identity.Identity) error {
	requstBody, err := json.Marshal(map[string]string{
		"NIK": nik,
	})

	req, err := http.NewRequest("POST", config.Of.Dukcapil.Endpoint, bytes.NewBuffer(requstBody))
	if err != nil {
		log.Fatalln(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Fail generate fake identity. Bad request")
	}

	var decoded struct {
		Content []interface{}
	}

	errDecode := json.NewDecoder(resp.Body).Decode(&decoded)
	if errDecode != nil {
		log.Println("errdecode", errDecode)
	}

	thecontents, errMarshalContent := json.Marshal(decoded.Content[0])
	if errMarshalContent != nil {
		log.Println(errMarshalContent)
	}

	json.Unmarshal(thecontents, &identity)

	return nil
}
