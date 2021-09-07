package server

import ("net/http" 
        "encoding/json"
	    "ggiz.statistic.user/cli"
		"ggiz.statistic.user/util"
		// "time"
		"os"
        "log"
	    "fmt"
		GgizAMQP "ggiz.statistic.user/ggizamqp"    
	)


const(
	MAP_ESTRUTURA_PRIMARIO int = 1
	MAP_ESTRUTURA_SECUNDARIO int = 2
)

var ggizAmqp *GgizAMQP.GGIZ_AMQP_CREDENTIAL
var requestinput int = 0
var procFinDay bool = false
var dataCurrent string = ""
var inputStatistic map[string][]cli.CliData = make(map[string][]cli.CliData)
var inputStatisticSec map[string][]cli.CliData = make(map[string][]cli.CliData)

// GGIZ_PUBLISH_ST_CLI 
type GGIZ_PUBLISH_ST_CLI struct {
	DATAPROC string `json: "dataProc"`
	LISTDATACLI []cli.CliData `json: "listDataCli"`
}

func createPublishQueue(ggizPublish GGIZ_PUBLISH_ST_CLI) {
	
	session, err := ggizAmqp.ConfigConnAMQP()
	
	if err != nil {
		fmt.Println(util.GetDateTime(), session.GGIZ_MESSAGE_ERROR)
		fmt.Println(util.GetDateTime(),` Detalhes do Erro : `, err)
		panic(err)
	}

	b, err := json.Marshal(ggizPublish)

	if err != nil {
      panic(err)
	}

	session.GGIZ_DATA_PUBLISH = b
	ok, err := ggizAmqp.Senddata(session)

	if err != nil {
		fmt.Println(util.GetDateTime(),` Ocorreu um erro na publicação da mensagem no Broker`)
		fmt.Println(util.GetDateTime(),` Detalhes do Erro : `, err)
		panic(err)
	}
	
	if ok {
		fmt.Println(util.GetDateTime(),`Mensagem publicada com sucesso no Broker com sucesso`)
	}

}

func finDay(w http.ResponseWriter, r *http.Request) {
	fmt.Println(util.GetDateTime(),` Finalizando o dia `, dataCurrent)

	if (r.Method == "POST") {
		fmt.Println(util.GetDateTime(),` - Metodo Post`)
		dataCurrentAnt := dataCurrent
		dataCurrent = util.GetDate()

		if (dataCurrentAnt == dataCurrent) {

			fmt.Println(util.GetDateTime(), ` Não é possivel enviar dados para fila. Existe um problema com os atributos de data.`)
			fmt.Println(util.GetDateTime(),`  Data Corrente: `, dataCurrent)
			fmt.Println(util.GetDateTime(),`  Data Anterior: `, dataCurrentAnt)
			fmt.Println(util.GetDateTime(),`  Data Corrente precisa ser maior que a Data Anterior`)

		} else {
			if (len(inputStatistic) > 0  && inputStatistic[dataCurrentAnt] != nil) {
				procFinDay = true
				ggiz_publish_st_cli := GGIZ_PUBLISH_ST_CLI{
					DATAPROC: dataCurrentAnt,
					LISTDATACLI: inputStatistic[dataCurrentAnt], 
				}
	
				createPublishQueue(ggiz_publish_st_cli)
				procFinDay = false
				fmt.Println(util.GetDateTime(),` Enviando conteúdo >>>>>>>>>>>>>>> AMQP`)
				fmt.Println(`Conteudo informado na fila com sucesso....`)
				delete(inputStatistic , dataCurrentAnt)
				inputStatistic[dataCurrent] = inputStatisticSec[dataCurrent]
				delete(inputStatisticSec, dataCurrent)
				fmt.Println(`                          `)
			} else {
				fmt.Println(util.GetDateTime(),` - Não existem dados para enviar para a processamento na fila`)
			}
		}
		
	}

}

func processRequest (w http.ResponseWriter, r *http.Request) {
	fmt.Println(util.GetDateTime(),` - Solicitacao recebida`)
	
	if (r.Method == "POST") {
	  var bIn bool = false
      var strIp string

	  fmt.Println(util.GetDateTime(),` - Metodo Post`)
	  ggizWebCli :=  cli.CliData{}
	  json.NewDecoder(r.Body).Decode(&ggizWebCli)
	  ggizWebCli.DATE = util.GetDate()
	  
	  if (!procFinDay) {
		fmt.Println(util.GetDateTime(),` - Adicionando dados na Estrutura Primária`)
		bIn, strIp = inputDataInMap(ggizWebCli, MAP_ESTRUTURA_PRIMARIO)
	  } else {
		fmt.Println(util.GetDateTime(),` - Adicionando dados na Estrutura Secundária`)
	    bIn, strIp = inputDataInMap(ggizWebCli, MAP_ESTRUTURA_SECUNDARIO)
	  }
	
	  outPut(ggizWebCli, bIn, strIp)
	}
}

func outPut(ggizWebCli cli.CliData, refBool bool, refStr string) {

	if refBool {
		ggizWebCli.OutPut()
	} else {
		ggizWebCli.OutPutWarning(refStr)
	}

}

func inputDataInMap(refCli cli.CliData, ordem int) (bool, string) {

	var refArray []cli.CliData
	var strIpTarget string = refCli.IP
    // Não gostei da solução. Depois com a evolução, penso em algo melhor
	if ordem == MAP_ESTRUTURA_PRIMARIO {
		
		refArray = inputStatistic[dataCurrent]
		il :=  len(refArray)
	   
		if il > 0 {
			var ipExist bool = false

			for _, dataStruc := range refArray {
				ipExist =  (dataStruc.IP == strIpTarget)
				
				if ipExist {
					break;
				}
			}

			if !ipExist {
				inputStatistic[dataCurrent] = append((inputStatistic[dataCurrent]), refCli)
				return true, ""
			}

			return false, strIpTarget
			
		} else {
			inputStatistic[dataCurrent] = append((inputStatistic[dataCurrent]), refCli)
			return true, ""
		}
		
	
	} else if ordem == MAP_ESTRUTURA_SECUNDARIO {

		refArray = inputStatisticSec[dataCurrent]
		il :=  len(refArray)
	   
		if il > 0 {
			var ipExist bool = false

			for _, dataStruc := range refArray {
				ipExist =  (dataStruc.IP == strIpTarget)
				
				if ipExist {
					break;
				}
			}

			if !ipExist {
				inputStatisticSec[dataCurrent] = append((inputStatisticSec[dataCurrent]), refCli)
				return true, ""
			}

			return false, strIpTarget
			
		} else {
			inputStatisticSec[dataCurrent] = append((inputStatisticSec[dataCurrent]), refCli)
			return true, ""
		}

	
	}

	return false, ""
}

// Inicializa o serviço e aguarda requisições
func Request() {
	var porta string = os.Getenv("GGIZ_STATISTIC_PORT")
	dataCurrent = util.GetDate()
    fmt.Println(util.GetDateTime(), ` - Inicializando Servidor.`)
	fmt.Println(util.GetDateTime(), ` - Ouvindo na porta : `, porta)
	fmt.Println(`                                                  `)
	fmt.Println(util.GetDateTime(), `- Serviços Cadastrados`)
	fmt.Println(util.GetDateTime(), `- ggiz_statistic`)
	fmt.Println(util.GetDateTime(), `- ggiz_finish_day`)
	fmt.Println(`                                                  `)
	fmt.Println(`-----------------------------------------------------------------------`)
    fmt.Println(util.GetDateTime(),` - Carregar Dados de Conexão com a Fila`)

	ggizAmqp = &GgizAMQP.GGIZ_AMQP_CREDENTIAL{}
	ggizAmqp.Loadcredential()
	ggizAmqp.OutputCredential()
	fmt.Println(`-----------------------------------------------------------------------`)

	http.HandleFunc("/ggiz_statistic", processRequest)
	http.HandleFunc("/ggiz_finish_day", finDay)
	log.Fatal(http.ListenAndServe(porta, nil))	 
}
