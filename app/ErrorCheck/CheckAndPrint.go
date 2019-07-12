package ErrorCheck

import "log"

func Check(err error, failMessage string, successMessage string){
	if err != nil {
		log.Fatalf("%s: %s", failMessage, err)
	} else {
		log.Println(successMessage)
	}
}
