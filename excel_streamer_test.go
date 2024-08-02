package excelstreamer

import (
	"encoding/json"
	"fmt"
	"github.com/theo1893/excelstreamer/consts"
	"github.com/theo1893/excelstreamer/data_source"
	"github.com/theo1893/excelstreamer/format_print"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// Step1: Create a new excelstreamer
	s, err := NewExcelStreamWriter("", "conf/test_excel.yml", "test_excel_dst.xlsx")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	// Step2: Create an input channel, add it into the streamer above, and trigger it
	inputCh := make(chan interface{})
	s.AddInChan(inputCh)
	s.Trigger()

	// Optional: Register custom format func
	format_print.RegisterFormat("TIME", func(data string, args ...string) string {
		timestamp, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return consts.StringMinus
		}

		if len(args) != 1 {
			return consts.StringMinus
		}

		return time.Unix(timestamp, 0).Format(time.TimeOnly)
	})

	epoch := 10

	// Step3: Put data into the input channel, and listen to the status of streamer
	for i := 0; i < epoch; i++ {
		mockData := MockDataGenerate()
		select {
		// the streamer has been aborted for some reason
		case <-s.Aborted():
			close(inputCh)
			t.Log(s.Error())
			return
		case inputCh <- mockData:
		}
	}

	// * this is just for comparison
	//mockDataCollection := make([]string, 0)
	//for i := 0; i < epoch; i++ {
	//	mockDataCollection = append(mockDataCollection, MockDataGenerate())
	//}
	//for _, mockData := range mockDataCollection {
	//	select {
	//	case <-s.Aborted():
	//		close(inputCh)
	//		return
	//	case inputCh <- mockData:
	//	}
	//}

	// Step4: Close the input channel to stop the streamer
	close(inputCh)

	select {
	case <-s.Aborted():
		t.Log(s.Error())
		return
	}
}

func TestAnchor(t *testing.T) {
	locator := NewSheetLocator()
	data := MockProfile{
		Address: MockAddress{
			Country: "China",
			City:    "Shenzhen",
			Street:  "",
		},
		FirstName: "Theo",
	}
	anchor := "address.country=hina||address.city=Shenzhn||first_name=Theo"

	ds := make(map[string]interface{})
	b, _ := json.Marshal(data)
	d := json.NewDecoder(strings.NewReader(string(b)))
	d.UseNumber()
	d.Decode(&ds)

	t.Log(locator.Locate(data_source.MapDataSource(ds), anchor))
}

type MockProfile struct {
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	Alias           string      `json:"alias"`
	Address         MockAddress `json:"address"`
	Family          MockFamily  `json:"family"`
	College         string      `json:"college"`
	Company         string      `json:"company"`
	Salary          float64     `json:"salary"`
	Debt            float64     `json:"debt"`
	Age             int         `json:"age"`
	Height          int         `json:"height"`
	Weight          int         `json:"weight"`
	BornIn          uint32      `json:"born_in"`
	LengthOfService int         `json:"length_of_service"`
}

type MockFamily struct {
	Father   string `json:"father"`
	Mother   string `json:"mother"`
	Daughter string `json:"daughter"`
	Son      string `json:"son"`
}

type MockAddress struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func MockDataGenerate() string {
	FirstNameSet := []string{"Alex", "Theo", "Amily", "Brook", "Bruce", "Carl", "Ward"}
	LastNameSet := []string{"Smith", "Jones", "Williams", "Brown", "Taylor", "Davis", "Wilson", "Evans", "Johnson"}
	countrySet := []string{"China", "America"}
	citySet := []string{"CityA", "CityB", "CityC", "CityD", "CityE"}
	streetSet := []string{"StreetA", "StreetB", "StreetC", "StreetD"}

	randNum := rand.Int()
	mockData := MockProfile{
		FirstName: FirstNameSet[randNum%len(FirstNameSet)],
		LastName:  LastNameSet[randNum%len(LastNameSet)],
		Alias:     GetRandomString(4),
		Family: MockFamily{
			Father:   FirstNameSet[randNum%len(FirstNameSet)],
			Mother:   FirstNameSet[randNum%len(FirstNameSet)],
			Daughter: FirstNameSet[randNum%len(FirstNameSet)],
			Son:      FirstNameSet[randNum%len(FirstNameSet)],
		},
		College: GetRandomString(4),
		Company: GetRandomString(4),
		Salary:  rand.Float64(),
		Debt:    rand.Float64(),
		Height:  rand.Int(),
		Weight:  rand.Int(),
		Address: MockAddress{
			Country: countrySet[randNum%len(countrySet)],
			City:    citySet[randNum%len(citySet)],
			Street:  streetSet[randNum%len(streetSet)],
		},
		BornIn:          1722481906,
		LengthOfService: rand.Int(),
		Age:             rand.Int(),
	}

	b, _ := json.Marshal(mockData)
	return string(b)
}
