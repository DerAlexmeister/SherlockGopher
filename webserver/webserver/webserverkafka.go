package webserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"encoding/json"
	"strconv"

	sherlockneo "github.com/DerAlexx/SherlockGopher/sherlockneo"
	sherlockkafka "github.com/DerAlexx/SherlockGopher/sherlockkafka"
	"github.com/segmentio/kafka-go"
)

const (
	topicurl         = "sendurltocrawler"
	brokerAddress = "localhost:9092"
)


type KafkaWriter struct {
	writer kafka.Writer
}

func NewKafkaWriter(brokAddress string, topic string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.Writer{
			Addr:  kafka.TCP(brokAddress),
			Topic: topic,
		},
	}
}

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(brokAddress string, topic string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokAddress},
			Topic:   topic,
		}),
	}
}


/*
SherlockWebserver will be the webserver being the man in the middle between the frontend and the backend.
type SherlockWebserver struct {
	Driver     neo4j.Driver
}

New will return a new instance of the SherlockWebserver.
func New() *SherlockWebserver {
	ldriver, err := sherlockneo.GetNewDatabaseConnection()
	if err == nil {
		return &SherlockWebserver{
			Driver: ldriver,	
		}
	}
	return &SherlockWebserver{}
}


Helloping Ping will return just for testing purposes a pong. Like PING PONG.
func (server *SherlockWebserver) Helloping(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]string{
		"Message": "Yes i am here!",
	})
}


DropGraphTable should drop the neo4j table.
func (server *SherlockWebserver) DropGraphTable(context *gin.Context) {
	session, err := sherlockneo.GetSession(server.Driver)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Message": "A Problem occurred while trying to connect to the Database",
		})
	} else {
		_, err := sherlockneo.DropTable(session)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"Message": "A Problem occurred while trying to drop the Database",
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"Message": "Dropped the table.",
			})
		}
	}
}*/

/*
ReceiveURL will handle the requested url which should be crawled.
*/
func (server *SherlockWebserver) SendUrlToCrawler(ginctx *gin.Context, ctx context.Context,) {
	var url = sherlockkafka.NewKafkaRequestedURL()
	err := ginctx.BindJSON(url)

	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while reveiving Requested Url",
		})
	}
	if govalidator.IsURL(url.URL) {

		r := NewKafkaReader(brokerAddress, topicurl)
		w := NewKafkaWriter(brokerAddress, topicurl)

		tmp := &sherlockkafka.KafkaRequestedURL{
			URL: url.URL,
		}

		res1B, _ := json.Marshal(tmp)

		err := w.writer.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(0)),
			// create an arbitrary message payload for the value
			Value: res1B,
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		msg, err := r.reader.ReadMessage(ctx)
		tmpakk := &sherlockkafka.KafkaAkkRequestedURL{}
		err = json.Unmarshal(msg.Value, &tmpakk)
		if err != nil {
			panic("parsing json failed" + err.Error())
		}
		
		if err == nil && tmpakk.Status {
			ginctx.JSON(http.StatusOK, gin.H{
				"Status": "Fine",
			})
		} else {
			ginctx.JSON(http.StatusOK, gin.H{
				"Message": "The webserver cannot submit your URL to the Crawler, couldnt reach the crawler service",
			})
		}
	} else {
		ginctx.JSON(http.StatusBadRequest, gin.H{
			"Status": "Url was empty or malformed",
		})
	}
}


/*
ReceiveMetadata will get the status of the tasks inside the analyser and crawler queue.

func (server *SherlockWebserver) AskForMetadata(context *gin.Context) {

	responseAnalyser, errAnalyser := server.receiveStatusTaskQAnalyser(context)
	responseCrawler, errCrawler := server.receiveStatusTaskQCrawler(context)

	switch {
	case errCrawler != nil && errAnalyser != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, analyser and crawler services are unavailable",
		})
	case errCrawler != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, crawler service is unavailable",
		})
	case errAnalyser != nil:
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Couldnt get Metadata, analyser service is unavailable",
		})
	default:
		metaArray := fillMetaArray(responseCrawler, responseAnalyser)
		context.JSON(http.StatusOK, metaArray.metamap)
	}
}
*/

/*
ChangeState is used to change the state of the analyser/crawler service.
Can send the new status to one of them or both at once.
Targeted Service and the new status is transmitted via the post request.
func (server *SherlockWebserver) ChangeState(context *gin.Context) {
	var status = NewRequestedStatus()
	err := context.BindJSON(status)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"Status": "Error while reveiving Requested Status",
		})
	}

	statusToAnalyser, statusToCrawler, statuserr := setRequestedStatus(status)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": statuserr,
		})
	}

	switch status.Target {
	case "Crawler":
		server.sendStateToCrawler(statusToCrawler, context)
	case "Analyser":
		server.sendStateToAnalyser(statusToAnalyser, context)
	case "All":
		server.sendStateToCrawlerAndAnalyser(statusToCrawler, statusToAnalyser, context)
	default:
		context.JSON(http.StatusBadRequest, gin.H{
			"Status": "Unknown Target, expected Crawler, Analyser or All",
		})
	}
}
*/

/*
GetServiceStatus returns the status of the crawler/analyser service. (checks if they are running).

func (server *SherlockWebserver) GetServiceStatus(context *gin.Context) {

	responseAnalyser, errAnalyser := server.getStateFromAnalyser(context)
	responseCrawler, errCrawler := server.getStateFromCrawler(context)

	switch {
	case errAnalyser != nil && errCrawler != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
			"Crawler":  "Unknown",
		})
	case errAnalyser != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": "Unknown",
			"Crawler":  responseCrawler.GetState().String(),
		})
	case errCrawler != nil:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": responseAnalyser.GetState().String(),
			"Crawler":  "Unknown",
		})
	default:
		context.JSON(http.StatusOK, gin.H{
			"Analyser": responseAnalyser.GetState().String(),
			"Crawler":  responseCrawler.GetState().String(),
		})
	}*/
}



