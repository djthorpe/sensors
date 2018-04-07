/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved

   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// RPC Services
	_ "github.com/djthorpe/sensors/cmd/mihome_gateway/service"
)

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("service/mihome:grpc")

	// Set the RPCServiceRecord for server discovery
	config.Service = "mihome"

	// Run the server and register all the services
	// Note the CommandLoop needs to go last as it blocks on Receive() until
	// Cancel is called from the CommandCancel task
	os.Exit(gopi.RPCServerTool(config))
}
