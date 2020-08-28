package main

import (
	"fmt"
	"github.com/djumanoff/amqp"
	"log"
	products "github.com/kirigaikabuto/products-lib"
)

var cfg = amqp.Config{
	Host: "localhost",
	VirtualHost: "",
	User: "",
	Password: "",
	Port: 5672,
	LogLevel: 5,
}

var srvCfg = amqp.ServerConfig{
	ResponseX: "response",
	RequestX: "request",
}

func main(){
	fmt.Println("Start")

	sess := amqp.NewSession(cfg)

	if err := sess.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer sess.Close()

	srv, err := sess.Server(srvCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	postgreConf := products.Config{
		Host:             "localhost",
		User:             "kirito",
		Password:         "passanya",
		Port:             5432,
		Database:         "crm",
		ConnectionString: "",
		Params:           "sslmode=disable",
	}
	productStore, err := products.NewPostgreStore(postgreConf)
	if err != nil {
		log.Fatal(err)
	}
	productService := products.NewProductService(productStore)
	productsAmqpEndpoints := products.NewAMQPEndpointFactory(productService)
	srv.Endpoint("products_lib.get",productsAmqpEndpoints.GetProductByIdAMQPEndpoint())
    srv.Endpoint("products_lib.create", productsAmqpEndpoints.CreateProductAMQPEndpoint())
	srv.Endpoint("products_lib.list", productsAmqpEndpoints.ListProductsAMQPEndpoint())
	srv.Endpoint("products_lib.delete", productsAmqpEndpoints.DeleteProductAMQPEndpoint())
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("End")
}