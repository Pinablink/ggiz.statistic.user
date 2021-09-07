package main

import ("fmt" 
        GgizServerStatistic "ggiz.statistic.user/server"
)

func main () {
	fmt.Println(`Servidor de GgizStatistic inicializado`)
	GgizServerStatistic.Request()
}
