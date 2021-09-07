package cli

import "fmt"

// Possui os dados obtidos no webcli
type CliData struct {
	IP string `json: "ip"`
	USW int  `json: "uSw"`
	USH int  `json: "uSh"`
	SUCESS bool `json: "sucess"`
	DATE string `json: "date"`
}

// Teste de Chamada
func Chamada() {
	fmt.Println(`Teste de Chamada`)
}

// Exibi os dados vindo do WebCli
func (cliDataRef CliData) OutPut() {
	fmt.Println(`---------------------------------------`)
	fmt.Println(`         Dados informados`)
	fmt.Println(`---------------------------------------`)
	fmt.Println(`ip`, cliDataRef.IP)
	fmt.Println(`sw`, cliDataRef.USW)
	fmt.Println(`sh`, cliDataRef.USH)
	fmt.Println(`sucess`, cliDataRef.SUCESS)
	fmt.Println(`date`, cliDataRef.DATE)
	fmt.Println(`---------------------------------------`)
}

// 
func (cliDaraRef CliData) OutPutWarning(strIp string) {
	fmt.Println(`---------------------------------------`)
	fmt.Println(`O IP informado já esta em nossa estrutura `)
	fmt.Println(`de visitação diária.`)
	fmt.Println(`IP informado: `, strIp)
	fmt.Println(`---------------------------------------`)
}
