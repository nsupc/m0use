package ns

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type region struct {
	Name    string `xml:"id,attr"`
	Nations []string
}

func (r *region) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			r.Name = attr.Value
		}
	}

	type Aux struct {
		Nations string `xml:"NATIONS"`
	}

	var aux Aux
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}

	r.Nations = strings.Split(aux.Nations, ":")

	return nil
}

// retrieve list of all nations in region
func (c *Client) GetNations(regionName string) ([]string, error) {
	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?region=%s&q=nations", regionName)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := region{}

	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return r.Nations, nil
}

func (c *Client) GetRecruitmentEligibleNations(inRegion string, fromRegion string) ([]string, error) {
	nations, err := c.GetNations(inRegion)
	if err != nil {
		return nil, err
	}

	result := []string{}

	for _, nationName := range nations {
		slog.Debug("checking recruitment eligibility status", slog.String("nation", nationName))
		status, err := c.IsRecruitmentEligible(nationName, fromRegion)
		if err != nil {
			slog.Warn("unable to get recruitment eligibility status", slog.String("nation", nationName), slog.Any("error", err))
			continue
		}

		if status.CanRecruit {
			result = append(result, nationName)
		}
	}

	return result, nil
}
