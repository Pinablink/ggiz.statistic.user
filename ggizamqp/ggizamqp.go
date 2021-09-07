package ggizamqp

import ("os"
        "fmt"
		"github.com/streadway/amqp"
		"ggiz.statistic.user/util"
	)

// 
type GGIZ_SESSION_AMQP struct {
	GGIZ_CONN_AMQP *amqp.Connection
	GGIZ_CHANNEL_AMQP *amqp.Channel
	GGIZ_QUEUE_AMQP amqp.Queue
	GGIZ_DATA_PUBLISH []byte
	GGIZ_MESSAGE_ERROR string
	
}
	
// Possui os dados de conectividade com o Broker RabbitMq
type GGIZ_AMQP_CREDENTIAL struct {
	GGIZ_AMQP_URL  string
    GGIZ_AMQP_USER string 
    GGIZ_AMQP_PASSWORD string 
	GGIZ_AMQP_HOST string
	GGIZ_AMQP_JOB_QUEUE string
}

// Carrega os recursos de acesso ao broker
func (c *GGIZ_AMQP_CREDENTIAL) Loadcredential() {
     c.GGIZ_AMQP_URL = os.Getenv("GGIZ_AMQP_URL")
	 c.GGIZ_AMQP_USER = os.Getenv("GGIZ_AMQP_USER")
	 c.GGIZ_AMQP_PASSWORD = os.Getenv("GGIZ_AMQP_PASSWORD")
	 c.GGIZ_AMQP_HOST = os.Getenv("GGIZ_AMQP_HOST")
	 c.GGIZ_AMQP_JOB_QUEUE = os.Getenv("GGIZ_AMQP_JOB_QUEUE")
}

// Expoe no console as credencias para conectar ao AMQP
func (c *GGIZ_AMQP_CREDENTIAL) OutputCredential() {
	fmt.Println(`Credencial de Acesso ao Recurso de Fila`)
	fmt.Println(util.GetDateTime(), `AMQP_URL :`, c.GGIZ_AMQP_URL)
	fmt.Println(util.GetDateTime(), `AMQP_HOST :`, c.GGIZ_AMQP_HOST)
	fmt.Println(util.GetDateTime(), `AMQP_USER :`, c.GGIZ_AMQP_USER)
	fmt.Println(util.GetDateTime(), `AMQP_PASSWORD : `, c.GGIZ_AMQP_PASSWORD)
	fmt.Println(util.GetDateTime(), `AMQP_JOB_QUEUE`, c.GGIZ_AMQP_JOB_QUEUE)
}

// Configura a conexão com o Broker
func (c *GGIZ_AMQP_CREDENTIAL) ConfigConnAMQP() (*GGIZ_SESSION_AMQP, error) {
	
	fmt.Println(``)
	fmt.Println(util.GetDateTime(),` Procedimento de Conexão com o cloudamqp`)
	
	ref := &GGIZ_SESSION_AMQP{}

	fmt.Println(util.GetDateTime(),` Abrindo conexão com o Broker`)
	conn, err := amqp.Dial(c.GGIZ_AMQP_URL)

	if err != nil {
		ref.GGIZ_MESSAGE_ERROR = "Problemas no procedimento de abertura de conexão"
		return nil, err
	}

	ref.GGIZ_CONN_AMQP = conn

	fmt.Println(util.GetDateTime(), ` Conectado ao Broker com sucesso`)

	fmt.Println(util.GetDateTime(),` Abrindo Canal`)
	channel, err := conn.Channel()

	if err != nil {
		ref.GGIZ_MESSAGE_ERROR = "Problemas no procedimento de abertura de canal"
		return ref, err
	}

	ref.GGIZ_CHANNEL_AMQP = channel

	fmt.Println(util.GetDateTime(), ` Canal aberto com sucesso`)

	fmt.Println(util.GetDateTime(),` Configurando Fila`)
	
	queue, err := channel.QueueDeclare (
		"input_statistic", 
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil)

    if err != nil {
		ref.GGIZ_MESSAGE_ERROR = "Problemas na configuração da fila"
		return ref, err
	}

	fmt.Println(util.GetDate(),` -> `, queue)

	ref.GGIZ_QUEUE_AMQP = queue
	 
	return ref, nil
}

// Envia os Dados informados para futuro processamento na fila
func (c *GGIZ_AMQP_CREDENTIAL) Senddata(refSession *GGIZ_SESSION_AMQP) (bool, error) {
	
	defer refSession.GGIZ_CONN_AMQP.Close()
	defer refSession.GGIZ_CHANNEL_AMQP.Close()

	err := refSession.GGIZ_CHANNEL_AMQP.Publish(
		"GGIZ_INPUT_STATISTIC",
		"statistic_ip",
		false,
		false,
		amqp.Publishing {
			ContentType: "application/json",
			Body: refSession.GGIZ_DATA_PUBLISH,
		},
	)

	if err != nil {
		return false, err
	}
	
	return true, nil
}
