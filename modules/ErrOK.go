package modules

import (
	"log"
)

func Critical(err error) {
	if err != nil {
		log.Println("치명적인 오류 발생: ", err)
		CriticalAlert(err)
		log.Fatal(err)
	}
}
func ErrOK(err error) error {
	if err != nil {
		log.Println("ErrOK: ", err)
		return err
	} else {
		return nil
	}
}
