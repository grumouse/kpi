package kpi

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"time"

	"net/http"

	"github.com/gorilla/schema"
	"github.com/grumouse/kpi/dripper"
)

const (
	kpiURL   = "https://development.kpi-drive.ru/_api/facts/save_fact" // url
	kpiToken = "48ab34464a5573519725deb5865cc74c"                      // токен
)

type KPIRequest struct {
	PeriodStart       time.Time `schema:"period_start"`
	PeriodEnd         time.Time `schema:"period_end"`
	PeriodKey         string    `schema:"period_key"`
	IndicatorToMoID   int       `schema:"indicator_to_mo_id"`
	IndicatorToFactID int       `schema:"indicator_to_mo_fact_id"`
	Value             int       `schema:"value"`
	FactTime          time.Time `schema:"fact_time"`
	IsPlan            int       `schema:"is_plan"`
	AuthUserID        int       `schema:"auth_user_id"`
	Comment           string    `schema:"comment"`
}

type Client struct {
	queue *dripper.Dripper
}

func NewClient() *Client {
	return &Client{queue: dripper.NewDripper(1024, 1)}
}

var encoder = func() *schema.Encoder {
	enc := schema.NewEncoder()

	enc.RegisterEncoder(time.Time{}, func(rv reflect.Value) string {
		return rv.Interface().(time.Time).Format("2006-01-02")
	})

	return enc
}()

func (r *Client) Do(data *KPIRequest) {

	//  создаём payload
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	form := url.Values{}
	err := encoder.Encode(data, form)
	if err != nil {
		panic("encode err:" + err.Error())
	}
	for k, v := range form {
		_ = writer.WriteField(k, v[0])
	}

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		panic("encode error")
	}

	// функция запроса
	f := func() bool {
		client := &http.Client{}
		req, err := http.NewRequest("POST", kpiURL, payload)

		if err != nil {
			fmt.Println(err)
			return false
		}
		req.Header.Add("Authorization", "Bearer "+kpiToken)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return false
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return false
		}

		fmt.Println("OK:", string(body))
		return true
	}

	r.queue.Do(f)
}

func (r *Client) Wait() { r.queue.Wait() }
