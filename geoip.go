package orca

import (
	"encoding/json"
	"errors"
	"net/http"
)

const maxmindHostname = "https://geoip.maxmind.com/geoip/v2.1/city/"
const mmContentType = "application/vnd.maxmind.com-city+json; charset=UTF-8; version=2.1"
const mmErrorType = "application/vnd.maxmind.com-error+json; charset=UTF-8; version=2.1"

// MaxMindClient is a utility for making geoip requests to MaxMind services.
type MaxMindClient struct {
	userid  string       // MaxMind user id
	license string       // MaxMind license key
	client  *http.Client // http client for making requests
}

// NewMaxMindClient initializes the client struct with required info and also
// initializes the http client to make requests with HTTP Authentication.
func NewMaxMindClient(userid string, license string) *MaxMindClient {
	return &MaxMindClient{
		userid:  userid,
		license: license,
		client:  new(http.Client),
	}
}

// NewRequest creates an authenticated request to look up the given IP address
// NOTE that you can also use the string 'me' which performs auto lookup, if
// an empty string is passed to this function, then 'me' is used.
func (mm *MaxMindClient) NewRequest(ipaddr string) (*http.Request, error) {

	// Check the credentials to ensure they're set.
	if mm.userid == "" || mm.license == "" {
		return nil, errors.New("Cannot make GeoIP requests without MaxMind user id and license key!")
	}

	// Empty string converts to the autolookup string
	// NOTE we could have made this an ExternalIP lookup, but this is usually
	// behind a NAT and therefore produces poor results.
	if ipaddr == "" {
		ipaddr = "me"
	}

	// Initialize the endpoint and create the GET request.
	endpoint := maxmindHostname + ipaddr
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add the authentication headers to the request
	request.SetBasicAuth(mm.userid, mm.license)

	// Add the Accept header to the request
	request.Header.Set("Accept", mmContentType)

	return request, nil
}

// GeoIPLookup fetches the raw JSON data from the MaxMind API for the given
// IP Address. Note that 'me' is a special input to autolookup the IP address
// from the request, and empty strings are converted to 'me'. This function
// returns the parsed JSON data as a map object.
func (mm *MaxMindClient) GeoIPLookup(ipaddr string) (map[string]interface{}, error) {

	// Construct the request from the IP Adddress
	req, err := mm.NewRequest(ipaddr)
	if err != nil {
		return nil, err
	}

	// Perform the GET request against the API
	resp, err := mm.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close the Body of the response when we're done with it.
	defer resp.Body.Close()

	// Decode the JSON from the response body
	decoder := json.NewDecoder(resp.Body)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {

		// Possibly this is because there is no JSON due to a status code
		if resp.StatusCode != 200 {
			return nil, errors.New(resp.Status)
		}

		// Otherwise just return the decoding error
		return nil, err
	}

	// Handle errors that have JSON bodies associated with them.
	if resp.StatusCode != 200 {
		if val, ok := data["error"]; ok {
			return nil, errors.New(val.(string))
		}

		return nil, errors.New(resp.Status)
	}

	return data, nil
}

// Nested lookup searches a map for a value. (Helper utility function for
// parsing the data returned by the MaxMind API)
func nestedLookup(data map[string]interface{}, keys []string) interface{} {
	if val, ok := data[keys[0]]; ok {
		if len(keys)-1 > 0 {
			sub, ok := val.(map[string]interface{})
			if ok {
				return nestedLookup(sub, keys[1:])
			}

		} else {
			return val
		}
	}

	return nil
}

// Helper function that converts nested lookup to a string.
func nestedStringLookup(data map[string]interface{}, keys []string) string {
	val := nestedLookup(data, keys)
	if con, ok := val.(string); ok {
		return con
	}
	return ""
}

// Helper function that converts nested lookup to a float64.
func nestedFloat64Lookup(data map[string]interface{}, keys []string) float64 {
	val := nestedLookup(data, keys)
	if con, ok := val.(float64); ok {
		return con
	}
	return 0.0
}

// GetCurrentLocation returns the current location using the special 'me' and
// returns a Location struct ready to be saved to the database.
func (mm *MaxMindClient) GetCurrentLocation() (*Location, error) {

	// Perform the GeoIP lookup
	data, err := mm.GeoIPLookup("")
	if err != nil {
		return nil, err
	}

	// Begin parsing the data
	loc := new(Location)
	loc.IPAddr = nestedStringLookup(data, []string{"traits", "ip_address"})
	loc.Latitude = nestedFloat64Lookup(data, []string{"location", "latitude"})
	loc.Longitude = nestedFloat64Lookup(data, []string{"location", "longitude"})
	loc.City = nestedStringLookup(data, []string{"city", "names", "en"})
	loc.PostCode = nestedStringLookup(data, []string{"postal", "code"})
	loc.Country = nestedStringLookup(data, []string{"country", "names", "en"})
	loc.Organization = nestedStringLookup(data, []string{"traits", "organization"})
	loc.Domain = nestedStringLookup(data, []string{"traits", "domain"})

	return loc, nil
}
