package productManager

import (
	pb "TapMars/admin_gateway/pkg/productManager/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"net/http"
)

type ProductManager struct {
	conn   *grpc.ClientConn
	client pb.ProductManagerClient
}

func NewProductManager(conn *grpc.ClientConn) *ProductManager {
	client := pb.NewProductManagerClient(conn)

	return &ProductManager{
		conn:   conn,
		client: client,
	}
}

func (p *ProductManager) Close() error {
	return p.conn.Close()
}

func (p *ProductManager) AddRoutes(rg *gin.RouterGroup) {
	p.addBusinessRoutes(rg)
	p.addItemRoutes(rg)
}

//func (p *ProductManager) addCommonRoutes(rg *gin.RouterGroup) {
//	rg.GET("/health", p.HealthCheck)
//}

func (p *ProductManager) addBusinessRoutes(rg *gin.RouterGroup) {
	businesses := rg.Group("/businesses")
	{
		businesses.POST("", p.CreateBusiness)
		businesses.GET("", p.QueryBusinesses)
		businesses.GET("/:id", p.GetBusiness)
		businesses.PATCH("/:id", p.UpdateBusiness)
		businesses.DELETE("/:id", p.DeleteBusiness)
		businesses.GET("/:id/items", p.QueryBusinessItems)
		businesses.POST("/:id/items", p.CreateItem)
	}
}

func (p *ProductManager) addItemRoutes(rg *gin.RouterGroup) {
	items := rg.Group("/items")
	{
		items.GET("/:id", p.GetItem)
		items.DELETE("/:id", p.DeleteItem)
	}
}

func (p *ProductManager) HealthCheck(c *gin.Context) {
	state := p.conn.GetState()

	c.JSON(http.StatusOK, gin.H{
		"state": state.String(),
	})
}

func (p *ProductManager) CreateBusiness(c *gin.Context) {
	var request *pb.BusinessProfile
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := p.client.CreateBusiness(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (p *ProductManager) QueryBusinesses(c *gin.Context) {
	setStreamHeaders(c)
	query, err := BuildBusinessesQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stream, err := p.client.QueryBusinesses(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for {
		business, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpcError(c, err)
			return
		}
		c.JSON(http.StatusOK, business)
	}
}

func (p *ProductManager) GetBusiness(c *gin.Context) {
	request := &pb.Id{Id: c.Param("id")}
	response, err := p.client.GetBusiness(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (p *ProductManager) UpdateBusiness(c *gin.Context) {
	var request *pb.BusinessProfile
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.Id = c.Param("id")
	response, err := p.client.UpdateBusiness(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (p *ProductManager) DeleteBusiness(c *gin.Context) {
	request := &pb.Id{Id: c.Param("id")}
	response, err := p.client.DeleteBusiness(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (p *ProductManager) GetItem(c *gin.Context) {
	request := &pb.Id{Id: c.Param("id")}
	response, err := p.client.GetItem(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (p *ProductManager) CreateItem(c *gin.Context) {
	var request *pb.ItemProfile
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.BusinessId = c.Param("id")
	response, err := p.client.CreateItem(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (p *ProductManager) DeleteItem(c *gin.Context) {
	request := &pb.Id{Id: c.Param("id")}
	response, err := p.client.DeleteItem(c, request)
	if err != nil {
		grpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (p *ProductManager) QueryBusinessItems(c *gin.Context) {
	setStreamHeaders(c)
	query := BuildBusinessItemsQuery(c)
	stream, err := p.client.QueryBusinessItems(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			grpcError(c, err)
			return
		}
		//c.Stream()
		//c.SSEvent()
		c.JSON(http.StatusOK, item)
	}
}
