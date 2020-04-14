package main

import (
	"fmt"

	"github.com/micro/go-micro"
	sherlockneo "github.com/ob-algdatii-20ss/SherlockGopher/analyser/ormneo4j"
	proto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto/analyser"
	crawlerproto "github.com/ob-algdatii-20ss/SherlockGopher/sherlockcrawler/proto/crawlertoanalyser"

	//streamproto "github.com/ob-algdatii-20ss/SherlockGopher/analyser/proto/filestreamproto"
	sherlockanalyser "github.com/ob-algdatii-20ss/SherlockGopher/analyser/sherlockanalyser"
)

const (
	serviceName      = "analyser-service"
	streamingService = "filestransfer-service-analyser"
)

func main() {

	// Analyserservice
	service := micro.NewService(
		micro.Name(serviceName),
	)
	service.Init()

	AnalyserService := sherlockanalyser.NewAnalyserServiceHandler()

	if driver, err := sherlockneo.GetNewDatabaseConnection(); err == nil {
		if session, sessionerror := sherlockneo.GetSession(&driver); sessionerror == nil {
			AnalyserService.InjectDependency(&sherlockanalyser.AnalyserDependency{
				Crawler: func() crawlerproto.AnalyserInterfaceService {
					return crawlerproto.NewAnalyserInterfaceService("crawler-service", service.Client()) // TODO: FIX BY DERALEXX
				}, Neo4J: &session,
			})
		}
		fmt.Println("Could not get a session to talk to the neo4j db. Service will shutdown.")
		//os.Exit(3)
	} else {
		fmt.Println("Could not reach the neo4j DB. Is the DB up?")
	}

	err := proto.RegisterAnalyserHandler(service.Server(), AnalyserService) // TODO: FIX BY DERALEXX

	if err != nil {
		fmt.Println(err)
	} else if lerr := service.Run(); lerr != nil {
		fmt.Println(lerr)
	} else {
		fmt.Printf("Service %s started as intended... ", serviceName)
	}
	// FileTransferService.

	streamingservice := micro.NewService(
		micro.Name(streamingService),
	)

	streamingservice.Init()
	/*
		newService := streamreceiver.NewServerGRPC()


		err1 := streamproto.RegisterReceiverHandler(service.Server(), newService)
		if err1 == nil {
			if err := service.Run(); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err1)
		}

	*/
}
