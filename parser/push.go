package parser

func push() {
	pj := &PushJson{}
	err = json.Unmarshal([]byte(strData), &pj)
}

type Student struct {
	title    string
	domain   string
	loginUrl string
	identity string
	writer   string
	fill     map[string]interface{}
}
