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
	"strings"
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
				detail, err:= detectKtp(w,r,event.Source.UserID)
				log.Println("Ini result detectKtp : IDLine :" + detail.LineID + " | Ktp" + detail.Ktp)
				log.Println("Ini error detect KTP : ",err)
				if detail.LineID != "" {
					if detail.Ktp != "" {
						result, err := detectIntent(w,r,message.Text,event.Source.UserID)
						log.Println("detect intent running")
						if result.Intent == "CLOSINGS"{
							//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s", result.Answer))).Do(); err != nil {
							//	log.Print(err)
							//}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(result.Answer, carouselBuilder(message,event.ReplyToken))).Do(); err != nil {
								log.Print(err)
							}
						}
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
						if strings.ToLower(message.Text) == "menu" {
							//carouselBuilder(message, event.ReplyToken)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", carouselBuilder(message,event.ReplyToken))).Do(); err != nil {
								log.Print(err)
							}
							//} else if strings.ToLower(message.Text) == "transaksi" {
							//			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", transactionCarousel(message,event.ReplyToken))).Do(); err != nil {
							//				log.Print(err)
							//			}
						} else if strings.ToLower(message.Text) != "menu" {
							//carouselBuilder(message, event.ReplyToken)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s", result.Answer))).Do(); err != nil {
								log.Print(err)
							}
						}
					} else if detail.Ktp == "" || detail.Ktp == "null"{
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Anda belum terdaftar di BPJS.")).Do(); err != nil {
							log.Print(err)
						}
					}
				} else if detail.LineID == "" || detail.LineID == "null" {
					//detectKtp(w,r,event.Source.UserID)
					//if detail.LineID == "" || detail.LineID == "null" {
					//	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Id Line Anda belum terdaftar di sistem kami. Untuk memulai proses otentikasi, silahkan masukkan nomor KTP Anda.")).Do()
					//} else if detail.LineID != "" || detail.LineID != "null" {
						str2 := message.Text
						i1, err := strconv.ParseInt(str2, 10, 64)
						if err == nil {
							log.Println(i1)
							// panggil function utk cek available ktp
							avail, err := detectAvailKtp(w,r, message.Text)
							if avail.Ktp == "" {
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Pengguna dengan nomor KTP" + message.Text + " belum terdaftar di BPJS. Harap masukkan nomor KTP yang terdaftar.")).Do(); err != nil {
									log.Print(err)
								}
							} else  if avail.Ktp != "" {
								//panggil fungsi updateKTP
								updateLineUser(w,r, event.Source.UserID, message.Text)
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Id Line Anda berhasil tercatat untuk nomor KTP " + avail.Ktp + " silahkan ajukan pertanyaan Anda.")).Do(); err != nil {
									log.Print(err)
								}
								return
							}

						} else {
							log.Println("string error", i1)
							log.Println("input bukan ktp")
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Id Line Anda belum terdaftar di sistem kami. Untuk memulai proses otentikasi, silahkan masukkan nomor KTP Anda.")).Do(); err != nil {
								log.Print(err)
							}
							return
						}
					}


			}
		} else if event.Type == linebot.EventTypeFollow {
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Selamat datang di chatbot SUSAN. Silahkan mengajukan pertanyaan Anda")).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

func registerNewUser(w http.ResponseWriter, r *http.Request, userLineId string,ktp string) (UserDetail, error) {
	log.Println("masuk registerNewUser")
	var detail UserDetail

	reqBody := UserDetail{
		LineID : userLineId,
		Ktp:ktp,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/post/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Println("http request error")
		return UserDetail{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error response", resp)
	} else if err == nil {
		log.Println("no error", UserDetail{})
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			return UserDetail{}, err
		} else {
			return UserDetail{},err
		}
	}
	return UserDetail{}, err
}



func detectKtp(w http.ResponseWriter, r *http.Request, text string) (UserDetail, error) {
	log.Println("masuk detectKtp")
	var detail UserDetail

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

func detectAvailKtp(w http.ResponseWriter, r *http.Request, text string) (AvailKtpResponseModel, error) {
	log.Println("masuk detectKtp")
	var avail AvailKtpResponseModel

	reqBody := AvailKtpRequestModel{
		NoKTP : text,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/avail/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		return AvailKtpResponseModel{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return AvailKtpResponseModel{},err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&avail); err != nil {
			return AvailKtpResponseModel{},err
		} else {
			log.Println("INI AVAIL KTP : ",avail)
			return avail,nil
		}
	}
}


func detectIntent(w http.ResponseWriter, r *http.Request, text string, lineId string) (RuleBasedModel,error) {
	log.Println("masuk detectIntent")
	var result RuleBasedModel
	

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

func updateLineUser(w http.ResponseWriter, r *http.Request, userLineId string,ktp string) (UserDetail, error) {
	log.Println("masuk registerNewUser")
	var detail UserDetail

	reqBody := UserDetail{
		LineID : userLineId,
		Ktp:ktp,
	}

	reqBytes,err := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://susan-service.herokuapp.com/ktp/update/"), bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Println("http request error")
		return UserDetail{}, err
	}
	req.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error response", resp)
	} else if err == nil {
		log.Println("no error", UserDetail{})
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			return UserDetail{}, err
		} else {
			return UserDetail{},err
		}
	}
	return UserDetail{}, err
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

type AvailKtpRequestModel struct {
	NoKTP string `json:"no_ktp"`
}
type AvailKtpResponseModel struct {
	Ktp string `json:"ktp"`
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
		template := linebot.NewCarouselTemplate(
			linebot.NewCarouselColumn(
				"https://i.ibb.co/ggN2QJ4/Profile.jpg", "Profil", "Berisi berbagai macam informasi mengenai profil pelanggan",
				linebot.NewPostbackAction("Profil", "Profil", "Di provinsi mana bpjs saya terdaftar?", ""),
			),
			linebot.NewCarouselColumn(
				"https://i.ibb.co/G32j10f/Transaksi.jpg", "Transaksi", "Berisi berbagai macam informasi mengenai transaksi pelanggan",
				linebot.NewPostbackAction("Transaksi", "Transaksi", "Berapa biaya bpjs saya?", ""),
			),
			linebot.NewCarouselColumn(
				"https://i.ibb.co/svJSyy7/Riwayat.jpg", "Riwayat", "Berisi berbagai macam informasi mengenai riwayat pelanggan",
				linebot.NewPostbackAction("Riwayat", "Riwayat", "Apa jenis segmen saya?", ""),
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

//func transactionCarousel(message *linebot.TextMessage, replyToken string) *linebot.CarouselTemplate {
//	log.Println("masuk transactionCarousel")
//	template := linebot.NewCarouselTemplate(
//		linebot.NewCarouselColumn(
//			"", "Biaya", "keterangan biaya",
//			linebot.NewPostbackAction("Biaya", "Biaya", "Berapa biaya bpjs saya?", ""),
//		),
//		linebot.NewCarouselColumn(
//			"", "Tagihan", "Keterangan tagihan",
//			linebot.NewPostbackAction("Tagihan", "Tagihan", "Berapa tagihan bpjs saya?", ""),
//		),
//		linebot.NewCarouselColumn(
//			"", "Iuran", "Keterangan iuran",
//			linebot.NewPostbackAction("Iuran", "Iuran", "Berapa iuran bpjs saya", ""),
//		),
//	)
//	return template
//}