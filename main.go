// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strconv"

	//"firebase.google.com/go/db"
	//"firebase.google.com/go"
)

var bot *linebot.Client
// KitchenSink app
// type KitchenSink struct {
// 	bot         *linebot.Client
// 	appBaseURL  string
// 	downloadDir string
// }

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	log.Println(addr)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {


	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}


	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				//quota, err := bot.GetMessageQuota().Do()
				//if err != nil {
				//	log.Println("Quota err:", err)
				//}
				log.Println("Ini Text nya : " + message.Text)
				detail, _:= detectKtp(w,r,event.Source.UserID)
				//result, err := detectIntent(w,r,message.Text,event.Source.UserID)
				log.Println("Ini result detectKtp : IDLine :" + detail.LineID + "Ktp" + detail.Ktp)
				//detail1, err:= updateNoKTP(w,r,event.Source.UserID,message.Text)
				//log.Println("Ini result update :" + detail1.LineID + detail1.Ktp)
				//log.Println("Ini error detect intent : ",err)
				//log.Println("Ini result detect intent : " + result.Answer)
				log.Println("userId", event.Source.UserID)
				//log.Println("intent:", result.Intent)
				if detail.Ktp == "" {
					str2 := message.Text
					i2, err := strconv.ParseInt(str2, 10, 64)
					if err == nil {
						log.Println(i2)
						registerNewUser(w, r, event.Source.UserID, message.Text)
						return
					} else {
						log.Println("string error", i2)
						log.Println("Belum terdaftar")
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Anda belum terdaftar, silahkan masukkan nomor KTP Anda")).Do(); err != nil {
							log.Print(err)
						}
						return
					}

					// Terima kasih anda telah terdaftar silahkan ajukan pertanyaan anda atau ketik 'menu' untuk mengetahui layanan kami
					//}
					//else if detail.Ktp == "" {
					//	log.Println("No KTP tidak ada")
					//	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Tolong masukkan nomor ktp Anda")).Do(); err != nil {
					//		log.Print(err)
					//	}
					//	updateNoKTP(w,r,event.Source.UserID,message.Text)
					//	return
					//		}
				}else if detail.Ktp != "" {
					result, err := detectIntent(w,r,message.Text,event.Source.UserID)
							//if result.Intent == "CLOSINGS"{
							//	//log.Println("Run 1st")
							//	//bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s",result.Answer))).Do()
							//	////time.Sleep(2*time.Second)
							//	////log.Println("Run 2nd")
							//	//err := handleText(message, event.ReplyToken)
							//	//log.Println("Check Error : ",err)
							//	//log.Println("Reply Token : ", event.ReplyToken)
							//	//carousel := handleText(message,event.ReplyToken)
							//	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s",result.Answer)),linebot.NewTemplateMessage("Carousel alt text", carouselBuilder(message,event.ReplyToken))).Do(); err != nil {
							//		log.Print(err)
							//	}
							//}
							//else
							if message.Text == "Menu" || message.Text == "menu" {
								//carouselBuilder(message, event.ReplyToken)
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", carouselBuilder(message,event.ReplyToken))).Do(); err != nil {
							log.Print(err)
						}
					} else if message.Text == "transaksi" {
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", transactionCarousel(message,event.ReplyToken))).Do(); err != nil {
									log.Print(err)
								}
					} else if message.Text != "Menu" || message.Text != "menu" {
					//carouselBuilder(message, event.ReplyToken)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s", result.Answer))).Do(); err != nil {
							log.Print(err)
					}
				}
			}
				//if detail.Ktp == "" {
				//	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Tolong masukkan nomor ktp Anda"))
				//} else if detail.Ktp != "" {
					

				//}

			// case *linebot.ImageMessage:
			// 	if err := handleText(message, event.ReplyToken); err != nil {
			// 		log.Print(err)
			// 	}
			}
		} else if event.Type == linebot.EventTypeFollow {
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Selamat datang di chatbot SUSAN. Silahkan ajukan pertanyaan Anda")).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

func registerNewUser(w http.ResponseWriter, r *http.Request, userLineId string,ktp string) (UserDetail, error) {
	log.Println("masuk registerNewUser")
	var detail UserDetail

	// if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
	// 	return RuleBasedModel{},nil
	// }

	reqBody := UserDetail{
		LineID : userLineId,
		Ktp:ktp,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/post/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return UserDetail{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return UserDetail{},err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			return UserDetail{},err
		} else {
			log.Println("INI RESULT LINE ID dan KTP dari register : ",detail)
			events, _ := bot.ParseRequest(r)
			for _, event := range events {
				if event.Type == linebot.EventTypeMessage {
					switch err := event.Message.(type) {
					case *linebot.TextMessage:
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Terima kasih anda telah terdaftar"))
						log.Println(err)
					}
				}
			}
			return detail,nil
		}
	}
}

func updateNoKTP(w http.ResponseWriter, r *http.Request, userLineId string,ktp string) (UserDetail, error) {
	log.Println("masuk updateNoKTP")
	var detail UserDetail


	// if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
	// 	return RuleBasedModel{},nil
	// }

	reqBody := UserDetail{
		LineID : userLineId,
		Ktp:ktp,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/update/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return UserDetail{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return UserDetail{},err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			return UserDetail{},err
		} else {
			log.Println("INI RESULT LINE ID dan KTP dari update: ",detail)
			return detail,nil
		}
	}
}


func detectKtp(w http.ResponseWriter, r *http.Request, text string) (UserDetail, error) {
	log.Println("masuk detectKtp")
	var detail UserDetail

	// if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
	// 	return RuleBasedModel{},nil
	// }

	reqBody := KtpRequestModel{
		UserLineId : text,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return UserDetail{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return UserDetail{},err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			return UserDetail{},err
		} else {
			log.Println("INI RESULT KTP : ",detail)
			return detail,nil
		}
	}
}


func detectIntent(w http.ResponseWriter, r *http.Request, text string, lineId string) (RuleBasedModel,error) {
	log.Println("masuk detectIntent")
	var result RuleBasedModel
	

	// if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
	// 	return RuleBasedModel{},nil
	// }

	reqBody := RequestModel{
		Sentence : text,
		UserLineId:lineId,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://susan-service.herokuapp.com/listener/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return RuleBasedModel{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	
		return RuleBasedModel{},err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return RuleBasedModel{},err
		} else {
			log.Println("INI RESULT : ",result)
			return result,nil
		}
	}
}

type RuleBasedModel struct {
	Answer string `json:"answer"`
	Intent string `json:"intent"`
}

type RequestModel struct {
	Sentence string `json:"sentence"`
	UserLineId string `json:"userLineId"`
}

type KtpRequestModel struct {
	UserLineId string `json:"userLineId"`
}

type UserDetail struct {
	LineID string `json:"userLineId"`
	Ktp string `json:"ktp"`
}

func carouselBuilder(message *linebot.TextMessage, replyToken string) *linebot.CarouselTemplate {
	log.Println("masuk carouselBuilder")
	// switch message.Text {
	// case "carousel":
		log.Println("iki carousel")
		imageURL1 := "https://i.ibb.co/ggN2QJ4/Profile.jpg"
		imageURL2 := "https://i.ibb.co/G32j10f/Transaksi.jpg"
		imageURL3 := "https://i.ibb.co/svJSyy7/Riwayat.jpg"
		template := linebot.NewCarouselTemplate(
			linebot.NewCarouselColumn(
				imageURL1, "Profil", "Berisi berbagai macam informasi mengenai profil pelanggan",
				linebot.NewPostbackAction("profil", "profil", "profil", ""),
			),
			linebot.NewCarouselColumn(
				imageURL2, "Transaksi", "Berisi berbagai macam informasi mengenai transaksi pelanggan",
				linebot.NewPostbackAction("transaksi", "transaksi", "transaksi", ""),
			),
			linebot.NewCarouselColumn(
				imageURL3, "Riwayat", "Berisi berbagai macam informasi mengenai riwayat pelanggan",
				linebot.NewPostbackAction("riwayat", "riwayat", "riwayat", ""),
			),
		)
		//if _, err := bot.ReplyMessage(
		//	replyToken,
		//	linebot.NewTemplateMessage("Carousel alt text", template),
		//).Do(); err != nil {
		//	return err
		//}

		return template
	// default:
	// 	log.Printf("Echo message to %s: %s", replyToken, message.Text)
	// 	if _, err := bot.ReplyMessage(
	// 		replyToken,
	// 		linebot.NewTextMessage(message.Text),
	// 	).Do(); err != nil {
	// 		return err
	// 	}
	// }
	
}

func transactionCarousel(message *linebot.TextMessage, replyToken string) *linebot.CarouselTemplate {
	log.Println("masuk transactionCarousel")
	template := linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			"", "Biaya", "keterangan biaya",
			linebot.NewPostbackAction("Biaya", "Biaya", "Berapa biaya bpjs saya?", ""),
		),
		linebot.NewCarouselColumn(
			"", "Tagihan", "Keterangan tagihan",
			linebot.NewPostbackAction("Tagihan", "Tagihan", "Berapa tagihan bpjs saya?", ""),
		),
		linebot.NewCarouselColumn(
			"", "Iuran", "Keterangan iuran",
			linebot.NewPostbackAction("Iuran", "Iuran", "Berapa iuran bpjs saya", ""),
		),
	)
	return template
}