package productManager_test

import (
	pb "TapMars/admin_gateway/pkg/productManager/proto"
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/genproto/googleapis/type/timeofday"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"net/http/httptest"
)

//goland:noinspection SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection
var _ = Describe("Proxy", func() {
	var (
		businessId      string
		itemId          string
		businessProfile *pb.BusinessProfile
		itemProfile     *pb.ItemProfile
	)
	//Describe("Get Health", func() {
	//
	//	It("Should pass", func() {
	//		req, err := http.NewRequest("GET", "/health", nil)
	//		Expect(err).ToNot(HaveOccurred())
	//
	//		w := httptest.NewRecorder()
	//		router.ServeHTTP(w, req)
	//
	//		Expect(w.Code).To(Equal(http.StatusOK))
	//		var data map[string]interface{}
	//		err = json.NewDecoder(w.Body).Decode(&data)
	//		Expect(err).ToNot(HaveOccurred())
	//		//Expect(data["state"]).To(Equal("READY"))
	//		Expect(data["state"]).ToNot(BeNil())
	//	})
	//})
	Describe("Business and Item lifecycle", func() {

		It("Should Create Business", func() {
			businessProfile = &pb.BusinessProfile{
				Name: "Proxy",
				Address: &pb.Address{
					Lines:      []string{"Proxy rd", "#501"},
					City:       "Not Here",
					State:      "Kansas",
					Zip:        "12345",
					RegionCode: "US",
				},
				LatLng: &latlng.LatLng{
					Latitude:  39.1233472,
					Longitude: -94.818466,
				},
			}
			body, err := json.Marshal(businessProfile)
			Expect(err).ToNot(HaveOccurred())

			req, err := http.NewRequest("POST", "/businesses", bytes.NewBuffer(body))
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
			var id *pb.Id
			err = json.NewDecoder(w.Body).Decode(&id)
			Expect(err).ToNot(HaveOccurred())
			businessId = id.Id
			Expect(businessId).ToNot(BeZero())
		})
		It("Should Get Business", func() {
			path := fmt.Sprintf("/businesses/%s", businessId)
			req, err := http.NewRequest("GET", path, nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))

			var business *pb.Business
			err = json.NewDecoder(w.Body).Decode(&business)
			Expect(err).ToNot(HaveOccurred())

			Expect(business.Id).To(Equal(businessId))
			Expect(business.Name).To(Equal(businessProfile.Name))
		})
		It("Should Create Item", func() {
			//goland:noinspection SpellCheckingInspection
			itemProfile = &pb.ItemProfile{
				Name: "Margs",
				HappyHourPeriod: &pb.HappyHourPeriod{
					DayOfWeek: 1,
					Start: &timeofday.TimeOfDay{
						Hours:   10,
						Minutes: 30,
					},
					End: &timeofday.TimeOfDay{
						Hours:   17,
						Minutes: 0,
					},
				},
				Details: &pb.ItemDetails{
					Description: "$1 off",
					IsDrink:     true,
					IsFood:      false,
					IsOther:     false,
				},
			}
			body, err := json.Marshal(itemProfile)
			Expect(err).ToNot(HaveOccurred())

			path := fmt.Sprintf("/businesses/%s/items", businessId)
			req, err := http.NewRequest("POST", path, bytes.NewBuffer(body))
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusCreated))

			var id *pb.Id
			err = json.NewDecoder(w.Body).Decode(&id)
			Expect(err).ToNot(HaveOccurred())
			itemId = id.Id
			Expect(itemId).ToNot(BeZero())
		})
		It("Should Get Item", func() {
			path := fmt.Sprintf("/items/%s", itemId)
			req, err := http.NewRequest("GET", path, nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))

			var item *pb.Item
			err = json.NewDecoder(w.Body).Decode(&item)
			Expect(err).ToNot(HaveOccurred())
			Expect(item.Id).To(Equal(itemId))
			Expect(item.Name).To(Equal(item.Name))
		})
		It("Should Update Business", func() {
			businessProfile = &pb.BusinessProfile{
				Name: "Proxy Update",
				Address: &pb.Address{
					Lines:      []string{"Proxy st", "#117"},
					City:       "Dallas",
					State:      "Texas",
					Zip:        "77899",
					RegionCode: "US",
				},
				LatLng: &latlng.LatLng{
					Latitude:  32.8208751,
					Longitude: -96.8716357,
				},
			}
			body, err := json.Marshal(businessProfile)
			Expect(err).ToNot(HaveOccurred())

			path := fmt.Sprintf("/businesses/%s", businessId)
			req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(body))
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))

			var itemsAffected *pb.ItemsAffected
			err = json.NewDecoder(w.Body).Decode(&itemsAffected)
			Expect(err).ToNot(HaveOccurred())
			Expect(itemsAffected.Count).To(Equal(int32(1)))
		})
		Context("Search Businesses", func() {
			AssertSuccessfulBusinessSearch := func(req *http.Request) {
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				dec := json.NewDecoder(w.Body)
				for dec.More() {
					var business *pb.Business
					// decode an array value (Message)
					err := dec.Decode(&business)
					Expect(err).ToNot(HaveOccurred())
					Expect(business).ToNot(BeZero())
				}
			}
			It("With No Query", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Name", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				//goland:noinspection SpellCheckingInspection,SpellCheckingInspection
				query.Add("name", "Prox")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Sort", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("sort", "2")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Location", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("latitude", "32.8208751")
				query.Add("longitude", "-94.818466")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Filter Distance", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("filter-distance", "1")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Location and Filter Distance", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("latitude", "32.8208751")
				query.Add("longitude", "-94.818466")
				query.Add("filter-distance", "1")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With Everything", func() {

				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("name", "Prox")
				query.Add("sort", "0")
				query.Add("latitude", "32.8208751")
				query.Add("longitude", "-94.818466")
				query.Add("filter-distance", "1")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessSearch(req)
			})
			It("With an empty list", func() {
				path := fmt.Sprintf("/businesses")
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("latitude", "0")
				query.Add("longitude", "0")
				query.Add("filter-distance", "3")
				req.URL.RawQuery = query.Encode()

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				dec := json.NewDecoder(w.Body)
				for dec.More() {
					Fail("This loop should not enter")
				}

			})
		})
		Context("Search Business Items", func() {
			AssertSuccessfulBusinessItemsSearch := func(req *http.Request) {
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var body map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&body)
				Expect(err).ToNot(HaveOccurred())
				itemList, found := body["items"]
				if found {
					var items []*pb.Item
					for _, b := range itemList.([]interface{}) {
						var item pb.Item
						data, err := json.Marshal(b)
						Expect(err).ToNot(HaveOccurred())
						err = protojson.Unmarshal(data, &item)
						items = append(items, &item)
					}
					Expect(items).ToNot(BeNil())
				}
			}
			It("With No Query", func() {

				path := fmt.Sprintf("/businesses/%s/items", businessId)
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())

				AssertSuccessfulBusinessItemsSearch(req)
			})
			It("With Day of Week", func() {

				path := fmt.Sprintf("/businesses/%s/items", businessId)
				req, err := http.NewRequest("GET", path, nil)
				Expect(err).ToNot(HaveOccurred())
				query := req.URL.Query()
				query.Add("day-of-week", "1")
				req.URL.RawQuery = query.Encode()

				AssertSuccessfulBusinessItemsSearch(req)
			})
		})
		It("Should Delete Item", func() {
			path := fmt.Sprintf("/items/%s", itemId)
			req, err := http.NewRequest("DELETE", path, nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
		It("Should Delete Business", func() {
			path := fmt.Sprintf("/businesses/%s", businessId)
			req, err := http.NewRequest("DELETE", path, nil)
			Expect(err).ToNot(HaveOccurred())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
