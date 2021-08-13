package productServer

import (
	"XsollaSchoolBE/DB"
	"XsollaSchoolBE/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

const baseUrl = "http://localhost:8080/api/v1/products"

var testProducts = []models.InputProduct{
	{"TEST1231", "Prod1", "Type1", 10},
	{"TEST1232", "Prod2", "Type3", 12},
	{"TEST1233", "Prod3", "Type1", 3},
	{"TEST1234", "Prod4", "Type2", 312},
	{"TEST1235", "Prod5", "Type4", 222},
	{"TEST1236", "Prod6", "Type1", 13},
	{"TEST1237", "Prod7", "Type5", 35},
	{"TEST1238", "Prod8", "Type5", 345},
	{"TEST1239", "Prod9", "Type1", 353},
	{"TEST12310", "Prod10", "Type2", 1},
}

func TestMain(m *testing.M) {
	srv, err := Run(":8080", "testDB.db")
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("The end of server tests")
	if err := os.Remove("testDB.db"); err != nil {
		log.Println("Warning: ", err.Error())
	}
}

func TestCorrectPost(t *testing.T) {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	} else if resp.StatusCode != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(fmt.Sprintf("\nBad status code: %d\nResponse body: %s\n", resp.StatusCode, bodyData))
	}
	for _, prod := range testProducts {
		jsonProduct, _ := json.Marshal(prod)
		resp, err := http.Post(
			baseUrl,
			"application/json",
			bytes.NewBuffer(jsonProduct))
		if err != nil {
			t.Error(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			bodyStr, _ := ioutil.ReadAll(resp.Body)
			t.Error(fmt.Sprintf("\nBad status code: %d\nResponse body: %s\n", resp.StatusCode, bodyStr))
		}
	}
}

func TestIncorrectPost(t *testing.T) {
	jsonProduct, _ := json.Marshal(map[string]string{"Cost": "TEST12310"})
	resp, err := http.Post(
		baseUrl,
		"application/json",
		bytes.NewBuffer(jsonProduct),
	)
	if err == nil {
		respData := make(map[string]interface{})
		if resp.StatusCode != http.StatusBadRequest {
			t.Error("not 400 code for incorrect post request")
		} else if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			t.Error("response format error, response ", fmt.Sprint(resp.Body))
		} else if errMsg, ok := respData["error"]; !ok {
			t.Error("no error filed in the response, response: ", respData)
		} else if !strings.HasPrefix(errMsg.(string), "json format error") {
			t.Error("wrong error message: ", errMsg)
		}
		resp.Body.Close()
	} else {
		t.Error(err)
	}

	jsonProduct, _ = json.Marshal(testProducts[0])
	resp, err = http.Post(
		baseUrl,
		"application/json",
		bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Error(err)
	} else {
		defer resp.Body.Close()
		respData := make(map[string]interface{})
		if resp.StatusCode != http.StatusBadRequest {
			t.Error("not 400 code for incorrect post request")
		} else if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			t.Error("response format error, response ", fmt.Sprint(resp.Body))
		} else if errMsg, ok := respData["error"]; !ok {
			t.Error("no error filed in the response, response: ", respData)
		} else if errMsg != DB.ProductAlreadyExistsError.Error() {
			t.Error("wrong error message: ", errMsg)
		}
	}
}

func TestGetAll(t *testing.T) {
	resp, err := http.Get(baseUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(fmt.Sprintf("\nBad status code: %d\nResponse body: %s\nResponse body:", resp.StatusCode, bodyData))
	}
	var receivedProducts []models.Product
	err = json.NewDecoder(resp.Body).Decode(&receivedProducts)
	if err != nil {
		t.Fatal(err)
	} else if len(testProducts) != len(receivedProducts) {
		t.Fatal("Wrong length of received products")
	} else {
		for i, prod := range receivedProducts {
			if testProducts[i] != prod.InputProduct {
				t.Fatalf("SKU mismatch:\n%v,\n%v", prod.InputProduct, testProducts[i])
			}
		}
	}
}

func TestGetGroupOfProducts(t *testing.T) {
	resp, err := http.Get(baseUrl + "?groupSize=3&groupNum=2")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(fmt.Sprintf("\nBad status code: %d\nResponse body: %s\nResponse body:", resp.StatusCode, bodyData))
	}
	var receivedProducts []models.Product
	err = json.NewDecoder(resp.Body).Decode(&receivedProducts)
	if err != nil {
		t.Fatal(err)
	} else if len(receivedProducts) != 3 {
		t.Fatal("Wrong length of received products: ", len(testProducts))
	} else {
		for i, prod := range receivedProducts {
			if testProducts[3+i] != prod.InputProduct {
				t.Fatalf("SKU mismatch:\n%v,\n%v", prod.InputProduct, testProducts[3+i])
			}
		}
	}
}

func TestCorrectGet(t *testing.T) {
	for _, url := range []string{
		baseUrl + "/" + testProducts[0].SKU,
		baseUrl + "?sku=" + testProducts[0].SKU,
		baseUrl + "?id=1",
	} {
		prodPtr, err, _ := getProductFromURL(url)
		if err != nil {
			t.Error(err)
		} else if prodPtr.InputProduct != testProducts[0] {
			t.Errorf("received product is not equal the sent one:\nreceived product: %v\nsent product: %v", prodPtr, testProducts[0])
		}
	}
}

func TestIncorrectGet(t *testing.T) {
	for _, url := range []string{
		baseUrl + "/WRONG",
		baseUrl + "?sku=WRONG",
		baseUrl + "?id=999999",
	} {
		_, _, code := getProductFromURL(url)
		if code != http.StatusNotFound {
			t.Errorf("not 404 code for nonexistent product, code: %d", code)
		}
	}

	_, _, code := getProductFromURL(baseUrl + "?id=WRONG")
	if code != http.StatusBadRequest {
		t.Errorf("not 400 code for nondigital product id: %d", code)
	}
}

func TestCorrectDelete(t *testing.T) {
	requestingURLs := []string{
		baseUrl + "/" + testProducts[0].SKU,
		baseUrl + "?sku=" + testProducts[1].SKU,
		baseUrl + "?id=3",
	}
	client := &http.Client{}

	for _, url := range requestingURLs {
		if _, err, _ := getProductFromURL(url); err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp, err := client.Do(request); err != nil {
			t.Error(err)
			continue
		} else {
			if resp.StatusCode != http.StatusNoContent {
				t.Error("Bad status code: ", resp.StatusCode)
				continue
			}
			resp.Body.Close()
		}

		if _, _, code := getProductFromURL(url); code != http.StatusNotFound {
			t.Error("deleted product has been found")
		}
	}
}

func TestIncorrectDelete(t *testing.T) {
	requestingURLs := []string{
		baseUrl + "/WRONG",
		baseUrl + "?sku=WRONG",
		baseUrl + "?id=999999",
	}
	client := &http.Client{}

	for _, url := range requestingURLs {
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Error(err)
			continue
		}
		if resp, err := client.Do(request); err != nil {
			t.Error(err)
		} else {
			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("not 400 code for nondigital product id: %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	}

	request, err := http.NewRequest(http.MethodDelete, baseUrl+"?id=WRONG", nil)
	if err != nil {
		t.Error(err)
	} else if resp, err := client.Do(request); err != nil {
		t.Error(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("not 404 code for nonexistent product: %d", resp.StatusCode)
		}
	}
}

func TestCorrectUpdate(t *testing.T) {
	requestingURLs := []string{
		baseUrl + "/" + testProducts[3].SKU,
		baseUrl + "?sku=" + testProducts[4].SKU,
		baseUrl + "?id=6",
	}
	checkURLs := []string{
		baseUrl + "/NewSKU0",
		baseUrl + "?sku=NewSKU1",
		baseUrl + "?id=6",
	}
	client := &http.Client{}
	for i, url := range requestingURLs {
		if _, err, _ := getProductFromURL(url); err != nil {
			t.Error(err)
			continue
		}

		newProduct := models.InputProduct{"NewSKU" + strconv.Itoa(i), "NewName", "NewType", uint(100 + i)}
		jsonProduct, _ := json.Marshal(newProduct)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonProduct))
		if err != nil {
			t.Error(err)
			continue
		}
		if resp, err := client.Do(request); err != nil {
			t.Error(err)
		} else {
			if resp.StatusCode != http.StatusOK {
				t.Error("Bad status code: ", resp.StatusCode)
			}
			resp.Body.Close()
		}

		if prod, err, code := getProductFromURL(checkURLs[i]); err != nil {
			t.Error(err)
		} else if code != http.StatusOK {
			t.Error("Bad status code: ", code)
		} else if prod.InputProduct != newProduct {
			t.Errorf("Updated product mismatch:\nSend product: %v\nReceived product: %v", newProduct, prod.InputProduct)
		}
	}
}

func TestIncorrectUpdate(t *testing.T) {
	requestingURLs := []string{
		baseUrl + "/WRONG",
		baseUrl + "?sku=WRONG",
		baseUrl + "?id=9999",
	}
	client := &http.Client{}
	for i, url := range requestingURLs {
		newProduct := models.InputProduct{"NewSKU" + strconv.Itoa(i), "NewName", "NewType", uint(100 + i)}
		jsonProduct, _ := json.Marshal(newProduct)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonProduct))
		if err != nil {
			t.Error(err)
			continue
		}
		if resp, err := client.Do(request); err != nil {
			t.Error(err)
		} else {
			if resp.StatusCode != http.StatusNotFound {
				t.Error("not 404 code for nonexistent product: ", resp.StatusCode)
			}
			resp.Body.Close()
		}
	}
	testProductJson, _ := json.Marshal(testProducts[7])

	request, err := http.NewRequest(http.MethodPut, baseUrl+"?id=WRONG", bytes.NewBuffer(testProductJson))
	if err != nil {
		t.Error(err)
	} else if resp, err := client.Do(request); err != nil {
		t.Error(err)
	} else {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("not 404 code for nonexistent product: %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
	request, err = http.NewRequest(http.MethodPut, baseUrl+"/"+testProducts[8].SKU, bytes.NewBuffer(testProductJson))
	if err != nil {
		t.Error(err)
	} else if resp, err := client.Do(request); err != nil {
		t.Error(err)
	} else {
		respData := make(map[string]interface{})
		if resp.StatusCode != http.StatusBadRequest {
			t.Error("not 400 code for incorrect put request: ", resp.StatusCode)
		} else if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			t.Error("response format error, response ", fmt.Sprint(resp.Body))
		} else if errMsg, ok := respData["error"]; !ok {
			t.Error("no error filed in the response, response: ", respData)
		} else if errMsg != DB.ProductAlreadyExistsError.Error() {
			t.Error("wrong error message: ", errMsg)
		}
		resp.Body.Close()
	}
}

func getProductFromReader(reader io.Reader) (*models.Product, error) {
	prods := []*models.Product{}
	err := json.NewDecoder(reader).Decode(&prods)
	if err != nil {
		return nil, err
	} else {
		return prods[0], nil
	}
}

func getProductFromURL(url string) (product *models.Product, err error, code int) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err, 0
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	product = nil
	err = nil
	if code != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("\nBad status code: %d\nResponse body: %s\n", resp.StatusCode, bodyData)
	} else {
		product, err = getProductFromReader(resp.Body)
	}
	return
}
