package productManager_test

import (
	"TapMars/admin_gateway/pkg/config"
	"TapMars/admin_gateway/pkg/productManager"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"testing"
)

var ctx context.Context
var router *gin.Engine
var pmProxy *productManager.ProductManager
var conn *grpc.ClientConn

//var opts []grpc.DialOption

func TestProductManagerTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProxyTest Suite")
}

var _ = BeforeSuite(func() {
	_ = os.Setenv("PORT", "4010")
	_ = os.Setenv("HOST", "dev.admin-gateway.tapmars.com")
	//_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/nicky/go/src/TapMars/admin_gateway/private/service_account/admin-gateway.json")

	ctx = context.Background()

	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	_, _, err := config.GetEnvironmentVariables()
	Expect(err).NotTo(HaveOccurred())
	//address := host + ":" + port
	//

	//// Create an identity token.
	//// With a global TokenSource tokens would be reused and auto-refreshed at need.
	//// A given TokenSource is specific to the audience.
	//audience := "https://product-manager-lfxeqzzqba-uc.a.run.app"
	//tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	//if err != nil {
	//	Fail(fmt.Sprintf("idtoken.NewTokenSource: %v", err))
	//}
	//token, err := tokenSource.Token()
	//if err != nil {
	//	Fail(fmt.Sprintf("TokenSource.Token: %v", err))
	//}
	//
	////// Add token to gRPC Request.
	//ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	//timeout := time.Second * 10
	router = gin.Default()
	//router.Use(productManager.RequestTimeoutWrapper(timeout))

	conn, err = NewConn("https://product-manager-lfxeqzzqba-uc.a.run.app:4010", false)
	Expect(err).NotTo(HaveOccurred())
	state := conn.GetState()
	print(state)
	pmProxy = productManager.NewProductManager(conn)
	pmRouter := router.Group("")
	pmProxy.AddRoutes(pmRouter)
})

var _ = AfterSuite(func() {
	//cancel()
	err := pmProxy.Close()
	Expect(err).NotTo(HaveOccurred())
})

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, isInsecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	//opts = append(opts, option.WithCredentialsFile("/Users/nicky/go/src/TapMars/admin_gateway/private/service_account/admin-gateway.json"))
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if isInsecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(host, opts...)
}
