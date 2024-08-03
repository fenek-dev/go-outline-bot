package main

import (
	"fmt"

	"github.com/fenek-dev/go-outline-bot/pkg/outline_client"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
)

func main() {
	// This is a placeholder for the main function
	cert := "03245423789118C805AF7C10AFE2E7F8AD257DEFF99DE264634B4DF380333DFB"
	apiUrl := "https://5.42.73.65:50279/gPHqLn1TGr8zo3fjqyw0Dw"

	_, err := outline_client.NewOutlineVPN(apiUrl, cert)
	if err != nil {
		panic(err)
	}

	paymentClient := payment_service.NewClient(&payment_service.Options{
		BaseUrl: "https://api.payment-service.com",
	}, nil, nil)

	fmt.Println(paymentClient)

}
