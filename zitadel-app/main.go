package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/zitadel/oidc/pkg/oidc"

	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
)

var (
	issuer = flag.String("issuer", "http://localhost:8080", "issuer of your ZITADEL instance (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	api    = flag.String("api", "localhost:8080", "gRPC endpoint of your ZITADEL instance (in the form: <instance>.zitadel.cloud:443 or <yourdomain>:443)")
	orgID  = flag.String("orgID", "256404670501685762", "orgID used in the example call")
)

func main() {
	// Define the flag
	logLevel := flag.String("ZITADEL_LOG_LEVEL", "info", "set the logging level")

	// Parse the flags
	flag.Parse()

	// Use the flag value
	fmt.Printf("Log Level is set to: %s\n", *logLevel)
	flag.PrintDefaults()

	//create a client for the admin api providing:
	//- issuer (e.g. https://acme-dtfhdg.zitadel.cloud)
	//- api (e.g. acme-dtfhdg.zitadel.cloud:443)
	//- scopes (including the ZITADEL project ID),
	//- a JWT Profile source token (e.g. path to your key json), if not provided, the file will be read from the path set in env var ZITADEL_KEY_PATH
	//zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath("key.json")),
	client, err := admin.NewClient(
		*issuer,
		*api,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}, // , zitadel.ScopeProjectID(*orgID)
		//zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(keyPath)),
	)
	if err != nil {
		log.Fatalln("could not create client", err)
	}
	defer func() {
		err := client.Connection.Close()
		if err != nil {
			log.Println("could not close grpc connection", err)
		}
	}()

	fmt.Println("ordID:", *orgID)

	ctx := context.Background()
	//call ZITADEL and print the name and creation date of the requested organisation
	//the call was successful if no error occurred
	resp, err := client.GetOrgByID(ctx, &pb.GetOrgByIDRequest{
		Id: *orgID,
	})
	if err != nil {
		log.Fatalln("call failed: ", err)
	}
	log.Printf("%s was created on %s", resp.Org.Name, resp.Org.Details.CreationDate.AsTime())
}
