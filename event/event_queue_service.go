package event

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gingerxman/eel/config"
//	"github.com/gingerxman/eel/log"
//	"strings"
//)
//
//const _MNS_QUEUE_RECEIVE_MESSAGE_TIMEOUT = 10
//const _MNS_QUEUE_MESSAGE_VISIBILITY_TIMEOUT = 30
//const _LOG_INTERVAL = 300
//
//type mnsQueueConf struct{
//	endpoint string
//	accessId string
//	accessKey string
//	queue string
//}
//
//var queueConf *mnsQueueConf
//
//type MessageHandler interface {
//	Handle(map[string]interface{}) error
//}
//
//var event2handler = make(map[string]MessageHandler)
//
//func handleMessage(queue mns.AliMNSQueue, resp *mns.MessageReceiveResponse) {
//	data := make(map[string]interface{})
//	messageData := make(map[string]interface{})
//	event := "__default__"
//
//	defer func(){
//		if err := recover(); err!=nil{
//			//errMsg := fmt.Sprintf("handle event '%v' panic: %v", event, err)
//			//beego.PushErrorWithExtraDataToSentry(errMsg, map[string]interface{}{
//			//	"message": data["Message"],
//			//}, nil)
//			log.Logger.Error(err)
//		}
//	}()
//
//	err := json.Unmarshal([]byte(resp.MessageBody), &data)
//	if err != nil {
//		log.Logger.Error(err)
//	}
//
//	err = json.Unmarshal([]byte(data["Message"].(string)), &messageData)
//	if err != nil {
//		log.Logger.Error(err)
//	}
//
//	event = messageData["_event_name"].(string)
//	if ret, err := queue.ChangeMessageVisibility(resp.ReceiptHandle, _MNS_QUEUE_MESSAGE_VISIBILITY_TIMEOUT); err != nil {
//		log.Logger.Error(err)
//	} else {
//		//handle event
//		canDeleteMessage := false
//		if handler, ok := event2handler[event]; ok {
//			err := handler.Handle(messageData)
//			if err != nil {
//				log.Logger.Error(err)
//			} else {
//			}
//
//			canDeleteMessage = true
//		} else {
//			log.Logger.Error(fmt.Sprintf("[mns_queue_service] no handler for event '%s'", event))
//			canDeleteMessage = true
//		}
//
//		//eel.Logger.Debug("delete it now: ", ret.ReceiptHandle)
//		if canDeleteMessage {
//			log.Logger.Debug("[mns_queue_service] delete message now: ", ret.ReceiptHandle)
//			if err := queue.DeleteMessage(ret.ReceiptHandle); err != nil {
//				log.Logger.Error(err)
//			}
//		}
//	}
//}
//
//func RegisterEventHandler(event string, handler MessageHandler) {
//	event2handler[event] = handler
//}
//
//type EventQueueService struct {
//}
//
//func NewEventQueueService() *EventQueueService {
//	service := new(EventQueueService)
//	return service
//}
//
//func (this *EventQueueService) Listen() {
//	defer func(){
//		if err := recover(); err!=nil{
//			log.Logger.Error(err)
//		}
//	}()
//
//	client := mns.NewAliMNSClient(queueConf.endpoint, queueConf.accessId, queueConf.accessKey)
//	queue := mns.NewMNSQueue(queueConf.queue, client)
//
//	messageCount := 0
//	fetchCount := 0
//
//	respChan := make(chan mns.MessageReceiveResponse)
//	errChan := make(chan error)
//	go func() {
//		defer func(){
//			if err := recover(); err!=nil{
//				log.Logger.Error(err)
//			}
//		}()
//
//		for {
//			select {
//			case resp := <-respChan:
//				{
//					messageCount += 1
//					go handleMessage(queue, &resp)
//				}
//			case err := <-errChan:
//				{
//					if strings.Contains(err.Error(), "code: MessageNotExist") {
//						log.Logger.Debug("no message, continue receive...")
//					} else {
//						log.Logger.Error(err)
//					}
//				}
//			}
//		}
//	}()
//
//	for {
//		fetchCount += 1
//		if fetchCount % _LOG_INTERVAL == 0 {
//			log.Logger.Warn(fmt.Sprintf("[mns_queue_service] receive for %d times, %d messages", fetchCount, messageCount))
//		}
//		queue.ReceiveMessage(respChan, errChan, _MNS_QUEUE_RECEIVE_MESSAGE_TIMEOUT)
//	}
//}
//
//func init() {
//	queueConf = new(mnsQueueConf)
//	queueConf.accessId = config.ServiceConfig.String("aliyun::MNS_ACCESS_ID")
//	queueConf.accessKey = config.ServiceConfig.String("aliyun::MNS_ACCESS_KEY")
//	queueConf.endpoint = config.ServiceConfig.String("aliyun::MNS_ENDPOINT")
//	queueConf.queue = config.ServiceConfig.String("aliyun::MNS_QUEUE")
//}
