package productManager

import (
	pb "TapMars/admin_gateway/pkg/productManager/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strconv"
)

func readBusinessProfile(c *gin.Context) (*pb.BusinessProfile, error) {
	var businessProfile *pb.BusinessProfile
	if err := c.ShouldBindJSON(&businessProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return nil, err
	//}
	//
	//businessProfile := &pb.BusinessProfile{}
	//err = protojson.Unmarshal(body, businessProfile)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return nil, err
	//}
	return businessProfile, nil
}

func BuildBusinessesQuery(c *gin.Context) (*pb.BusinessesQuery, error) {
	latitude := c.DefaultQuery("latitude", "0.0")
	longitude := c.DefaultQuery("longitude", "0.0")
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, err
	}
	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing longitude"})
		return nil, err
	}

	var filterDistance pb.FilterDistance
	switch c.Query("filter-distance") {
	case "0":
		filterDistance = pb.FilterDistance_None
	case "1":
		filterDistance = pb.FilterDistance_One
	case "2":
		filterDistance = pb.FilterDistance_Five
	case "3":
		filterDistance = pb.FilterDistance_Twenty
	default:
		filterDistance = pb.FilterDistance_None
	}

	var sort pb.Sort
	switch c.Query("sort") {
	case "0":
		sort = pb.Sort_NameRelevance
	case "1":
		sort = pb.Sort_Favorites
	case "2":
		sort = pb.Sort_Updated
	case "3":
		sort = pb.Sort_Created
	default:
		sort = pb.Sort_NameRelevance
	}

	return &pb.BusinessesQuery{
		Name: c.Query("name"),
		LatLng: &latlng.LatLng{
			Latitude:  lat,
			Longitude: lng,
		},
		FilterDistance: filterDistance,
		Sort:           sort,
	}, nil
}

func BuildBusinessItemsQuery(c *gin.Context) *pb.BusinessItemsQuery {
	var dayOfWeek dayofweek.DayOfWeek
	switch c.Query("day-of-week") {
	case "0":
		dayOfWeek = dayofweek.DayOfWeek_DAY_OF_WEEK_UNSPECIFIED
	case "1":
		dayOfWeek = dayofweek.DayOfWeek_MONDAY
	case "2":
		dayOfWeek = dayofweek.DayOfWeek_TUESDAY
	case "3":
		dayOfWeek = dayofweek.DayOfWeek_WEDNESDAY
	case "4":
		dayOfWeek = dayofweek.DayOfWeek_THURSDAY
	case "5":
		dayOfWeek = dayofweek.DayOfWeek_FRIDAY
	case "6":
		dayOfWeek = dayofweek.DayOfWeek_SATURDAY
	case "7":
		dayOfWeek = dayofweek.DayOfWeek_SUNDAY
	default:
		dayOfWeek = dayofweek.DayOfWeek_DAY_OF_WEEK_UNSPECIFIED
	}

	return &pb.BusinessItemsQuery{
		BusinessId: c.Param("id"),
		DayOfWeek:  dayOfWeek,
	}
}

func readItemProfile(w http.ResponseWriter, r *http.Request) (*pb.ItemProfile, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading json body object", http.StatusBadRequest)
		return &pb.ItemProfile{}, err
	}

	itemProfile := &pb.ItemProfile{}
	err = protojson.Unmarshal(body, itemProfile)
	if err != nil {
		http.Error(w, "Error unmarshalling request", http.StatusBadRequest)
		return &pb.ItemProfile{}, err
	}
	return itemProfile, nil
}

func grpcErrorResponse(w http.ResponseWriter, err error) {
	stat, ok := status.FromError(err)
	if ok {
		switch stat.Code() {
		case codes.Internal:
			http.Error(w, stat.Message(), http.StatusInternalServerError)
		case codes.NotFound:
			http.Error(w, stat.Message(), http.StatusNotFound)
		case codes.DeadlineExceeded:
			http.Error(w, stat.Message(), http.StatusGatewayTimeout)
		case codes.Unimplemented:
			http.Error(w, stat.Message(), http.StatusInternalServerError)
		default:
			message := fmt.Sprintf("Code: %v - Message: %s", stat.Code(), stat.Message())
			http.Error(w, message, http.StatusBadGateway)
		}
	} else {
		message := fmt.Sprintf("Not able to parse error: %v", err)
		http.Error(w, message, http.StatusBadGateway)
	}
}

func respond(w http.ResponseWriter, expectedHttpStatus int, protoMessage proto.Message, err error) {

	if err != nil {
		grpcErrorResponse(w, err)
		return
	}

	jsonResponse, err := protojson.Marshal(protoMessage)
	if err != nil {
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(expectedHttpStatus)
	_, _ = w.Write(jsonResponse)
}

func grpcError(c *gin.Context, err error) {
	stat, ok := status.FromError(err)
	if ok {
		switch stat.Code() {
		case codes.Internal:
			c.JSON(http.StatusBadRequest, gin.H{"error": stat.Message()})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": stat.Message()})
		case codes.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": stat.Message()})
		case codes.Unimplemented:
			c.JSON(http.StatusInternalServerError, gin.H{"error": stat.Message()})
		default:
			message := fmt.Sprintf("Code: %v - Message: %s", stat.Code(), stat.Message())
			c.JSON(http.StatusBadGateway, gin.H{"error": message})
		}
	} else {
		message := fmt.Sprintf("Not able to parse error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": message})
	}
}
