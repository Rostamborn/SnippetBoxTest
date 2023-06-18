package forms

type errors map[string][]string


func (e errors) Add(field, message string) {
    e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
    mes := e[field]
    if len(mes) == 0 {
        return ""
    }
    return mes[0]
}
