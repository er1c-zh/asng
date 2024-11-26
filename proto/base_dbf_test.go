package proto

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBaseDBF(t *testing.T) {
	d, err := os.ReadFile("base.dbf")
	if err != nil {
		t.Error(err)
		return
	}
	DBF, err := ParseBaseDBF(d)
	if err != nil {
		t.Error(err)
		return
	}
	j, _ := json.MarshalIndent(DBF.Data[1], "", "  ")
	t.Log(string(j))
}
