package productManager

import (
	pb "TapMars/productManager/pkg/proto"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net/http"
)

type ProductManager struct {
	conn   *grpc.ClientConn
	client pb.ProductManagerClient
}

func NewProductManager(serverAddress *string, opts []grpc.DialOption) (*ProductManager, error) {
	conn, err := grpc.Dial(*serverAddress, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewProductManagerClient(conn)

	return &ProductManager{
		conn:   conn,
		client: client,
	}, nil

}

func (pm *ProductManager) Close() {
	_ = pm.conn.Close()
}

func (pm *ProductManager) RegisterHandlers(router *mux.Router) {
	businessesRouter := router.PathPrefix("/businesses").Subrouter()
	itemsRouter := businessesRouter.PathPrefix("/{id}/items").Subrouter()

	router.Methods("GET").Path("/health").HandlerFunc(stateCheck)

	businessesRouter.Methods("GET").Path("/{id}").
		HandlerFunc(getBusiness)
	businessesRouter.Methods("POST").Path("").
		HandlerFunc(createBusiness)
	businessesRouter.Methods("PUT").Path("/{id}").
		HandlerFunc(updateBusiness)
	businessesRouter.Methods("DELETE").Path("/{id}").
		HandlerFunc(deleteBusiness)
	businessesRouter.Methods("GET").Path("").
		Queries("filterDistance", "{filter-distance}", "orderBy", "{order-by}").
		HandlerFunc(searchBusinesses)

	itemsRouter.Methods("GET").Path("/{id}").
		HandlerFunc(getItem)
	itemsRouter.Methods("POST").Path("").
		HandlerFunc(createItem)
	itemsRouter.Methods("DELETE").Path("/{id}").
		HandlerFunc(deleteItem)
	itemsRouter.Methods("GET").Path("").
		Queries("dayOfWeek", "{day-of-week}").
		HandlerFunc(searchItems)
}

func (pm *ProductManager) stateCheck(w http.ResponseWriter, r *http.Request) {
	pm.conn.GetState()
	fmt.Println("Not implemented")
}

func getBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func createBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func updateBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func deleteBusiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func searchBusinesses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func getItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func createItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}

func searchItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not implemented")
}
