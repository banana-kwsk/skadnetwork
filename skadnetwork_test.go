package skadnetwork_test

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"

	"github.com/google/uuid"

	"github.com/banana-kwsk/skadnetwork"
)

const (
	pem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPAYHdpbrKcTKi6qrRBB/TYN4w33jXAL0j9JMOqu5oIZoAoGCCqGSM49
AwEHoUQDQgAEBdF30K5pLjixuXnqiCNN/AgUK3DexfWqLzNOn2cZt0t9lMR8Y/Dl
MgSZN35Bv8gyUXt7xOK+hP8tDoOD2ir7bw==
-----END EC PRIVATE KEY-----
`

	v2_2 = `{
  "version" : "2.2",
  "ad-network-id" : "com.example",
  "campaign-id" : 42,
  "transaction-id" : "6aafb7a5-0170-41b5-bbe4-fe71dedf1e28",
  "app-id" : 525463029,
  "attribution-signature" : "MEYCIQDTuQ1Z4Tpy9D3aEKbxLl5J5iKiTumcqZikuY/AOD2U7QIhAJAaiAv89AoquHXJffcieEQXdWHpcV8ZgbKN0EwV9/sY",
  "redownload": true,
  "source-app-id": 1234567891,
  "fidelity-type": 1,
  "conversion-value": 20
}`

	v3_0__win = `{ 
  "version": "3.0", 
  "ad-network-id": "example123.skadnetwork", 
  "campaign-id": 42, 
  "transaction-id": "6aafb7a5-0170-41b5-bbe4-fe71dedf1e28", 
  "app-id": 525463029, 
  "attribution-signature": "MEYCIQD5eq3AUlamORiGovqFiHWI4RZT/PrM3VEiXUrsC+M51wIhAPMANZA9c07raZJ64gVaXhB9+9yZj/X6DcNxONdccQij", 
  "redownload": true, 
  "source-app-id": 1234567891, 
  "fidelity-type": 1, 
  "conversion-value": 20,
  "did-win": true
}`

	v3_0__lose = `{ 
  "version": "3.0",
  "ad-network-id": "example123.skadnetwork",
  "campaign-id": 42,
  "transaction-id": "f9ac267a-a889-44ce-b5f7-0166d11461f0",
  "app-id": 525463029,
  "attribution-signature": "MEUCIQDDetUtkyc/MiQvVJ5I6HIO1E7l598572Wljot2Onzd4wIgVJLzVcyAV+TXksGNoa0DTMXEPgNPeHCmD4fw1ABXX0g=",
  "redownload": true,
  "fidelity-type": 1,
  "did-win": false
}`

	v4_0__fine = `{
  "version": "4.0",
  "ad-network-id": "com.example",
  "source-identifier": "5239",
  "app-id": 525463029,
  "transaction-id": "6aafb7a5-0170-41b5-bbe4-fe71dedf1e30",
  "redownload": false,
  "source-domain": "example.com",
  "fidelity-type": 1,
  "did-win": true,
  "conversion-value": 63,
  "postback-sequence-index": 0,
  "attribution-signature": "MEUCIGRmSMrqedNu6uaHyhVcifs118R5z/AB6cvRaKrRRHWRAiEAv96ne3dKQ5kJpbsfk4eYiePmrZUU6sQmo+7zfP/1Bxo="
}`
	
	v4_0__coarse = `{
  "version": "4.0",
  "ad-network-id": "com.example",
  "source-identifier": "39",
  "app-id": 525463029,
  "transaction-id": "6aafb7a5-0170-41b5-bbe4-fe71dedf1e31",
  "redownload": false,
  "source-domain": "example.com", 
  "fidelity-type": 1, 
  "did-win": true,
  "coarse-conversion-value": "high",
  "postback-sequence-index": 0,
  "attribution-signature": "MEUCIQD4rX6eh38qEhuUKHdap345UbmlzA7KEZ1bhWZuYM8MJwIgMnyiiZe6heabDkGwOaKBYrUXQhKtF3P/ERHqkR/XpuA="
}`
)

func ref[T any](t T) *T { return &t }

func TestSignAndVerify(t *testing.T) {
	s, err := skadnetwork.NewSigner(pem)
	assert.NilError(t, err)

	nonce := uuid.MustParse("68483ef6-0ada-40df-ab6b-3d19a66330fa")
	timestamp, _ := time.Parse(time.RFC3339, "2022-05-06T10:00:00Z")

	for _, c := range []struct {
		in *skadnetwork.Params
	}{
		{
			&skadnetwork.Params{
				AdNetworkID:      "example123.skadnetwork",
				CampaignID:       42,
				ItunesItemID:     525463029,
				Nonce:            nonce,
				SourceAppStoreID: 1234567891,
				Timestamp:        timestamp,
			},
		},
		{
			&skadnetwork.Params{
				AdNetworkID:      "example123.skadnetwork",
				CampaignID:       42,
				ItunesItemID:     525463029,
				Nonce:            nonce,
				SourceAppStoreID: 1234567891,
				Timestamp:        timestamp,
				FidelityType:     skadnetwork.SKRenderedAds,
			},
		},
		{
			&skadnetwork.Params{
				AdNetworkID:      "example123.skadnetwork",
				CampaignID:       42,
				ItunesItemID:     525463029,
				Nonce:            nonce,
				SourceAppStoreID: 1234567891,
				Timestamp:        timestamp,
				FidelityType:     skadnetwork.SKRenderedAds,
			},
		},
	} {
		sig, err := s.Sign(c.in)
		assert.NilError(t, err)

		got, err := s.Verify(c.in, sig)
		assert.NilError(t, err)
		assert.Equal(t, got, true)
	}
}

func TestMarshalJSON(t *testing.T) {
	for _, c := range []struct {
		in   string
		want *skadnetwork.Postback
	}{
		{
			v2_2,
			&skadnetwork.Postback{
				Version:              "2.2",
				AdNetworkID:          "com.example",
				CampaignID:           42,
				TransactionID:        "6aafb7a5-0170-41b5-bbe4-fe71dedf1e28",
				AppID:                525463029,
				AttributionSignature: "MEYCIQDTuQ1Z4Tpy9D3aEKbxLl5J5iKiTumcqZikuY/AOD2U7QIhAJAaiAv89AoquHXJffcieEQXdWHpcV8ZgbKN0EwV9/sY",
				Redownload:           ref(true),
				SourceAppID:          ref[int64](1234567891),
				FidelityType:         ref(skadnetwork.SKRenderedAds),
				ConversionValue:      ref[uint8](20),
			},
		},
		{
			v3_0__win,
			&skadnetwork.Postback{
				Version:              "3.0",
				AdNetworkID:          "example123.skadnetwork",
				CampaignID:           42,
				TransactionID:        "6aafb7a5-0170-41b5-bbe4-fe71dedf1e28",
				AppID:                525463029,
				AttributionSignature: "MEYCIQD5eq3AUlamORiGovqFiHWI4RZT/PrM3VEiXUrsC+M51wIhAPMANZA9c07raZJ64gVaXhB9+9yZj/X6DcNxONdccQij",
				Redownload:           ref(true),
				SourceAppID:          ref[int64](1234567891),
				FidelityType:         ref(skadnetwork.SKRenderedAds),
				ConversionValue:      ref[uint8](20),
				DidWin:               ref(true),
			},
		},
		{
			v3_0__lose,
			&skadnetwork.Postback{
				Version:              "3.0",
				AdNetworkID:          "example123.skadnetwork",
				CampaignID:           42,
				TransactionID:        "f9ac267a-a889-44ce-b5f7-0166d11461f0",
				AppID:                525463029,
				AttributionSignature: "MEUCIQDDetUtkyc/MiQvVJ5I6HIO1E7l598572Wljot2Onzd4wIgVJLzVcyAV+TXksGNoa0DTMXEPgNPeHCmD4fw1ABXX0g=",
				Redownload:           ref(true),
				FidelityType:         ref(skadnetwork.SKRenderedAds),
				DidWin:               ref(false),
			},
		},
		{
			v4_0__fine,
			&skadnetwork.Postback{
				Version:              "4.0",
				AdNetworkID:          "com.example",
				SourceIdentifier:		"5239",
				AppID:                525463029,
				TransactionID:        "6aafb7a5-0170-41b5-bbe4-fe71dedf1e30",
				Redownload:           ref(false),
				SourceDomain:	ref[string]("example.com"),
				FidelityType:         ref(skadnetwork.SKRenderedAds),
				DidWin:               ref(true),
				ConversionValue:	ref[uint8](63),
				PostbackSequenceIndex: ref[int64](0),
				AttributionSignature: "MEUCIGRmSMrqedNu6uaHyhVcifs118R5z/AB6cvRaKrRRHWRAiEAv96ne3dKQ5kJpbsfk4eYiePmrZUU6sQmo+7zfP/1Bxo=",
			},
		},
		{
			v4_0__coarse,
			&skadnetwork.Postback{
				Version:              "4.0",
				AdNetworkID:          "com.example",
				SourceIdentifier:		"39",
				AppID:                525463029,
				TransactionID:        "6aafb7a5-0170-41b5-bbe4-fe71dedf1e31",
				Redownload:           ref(false),
				SourceDomain:	ref[string]("example.com"),
				FidelityType:         ref(skadnetwork.SKRenderedAds),
				DidWin:               ref(true),
				CoarseConversionValue:	ref[string]("high"),
				PostbackSequenceIndex: ref[int64](0),
				AttributionSignature: "MEUCIQD4rX6eh38qEhuUKHdap345UbmlzA7KEZ1bhWZuYM8MJwIgMnyiiZe6heabDkGwOaKBYrUXQhKtF3P/ERHqkR/XpuA=",
			},
		},
	} {
		var got skadnetwork.Postback
		err := json.Unmarshal([]byte(c.in), &got)
		assert.NilError(t, err)

		_, err = json.Marshal(&got)
		assert.NilError(t, err)

		assert.Check(t, is.DeepEqual(&got, c.want))
	}
}

func TestVerifyPostback(t *testing.T) {
	for _, c := range []struct {
		in   string
		want bool
	}{
		{v2_2, true},
		{v3_0__win, true},
		{v3_0__lose, true},
		{v4_0__fine, true},
		{v4_0__coarse, true},
	} {
		var p skadnetwork.Postback
		err := json.Unmarshal([]byte(c.in), &p)
		assert.NilError(t, err)

		got, err := skadnetwork.Verify(p)
		assert.NilError(t, err)

		assert.Equal(t, got, c.want)
	}
}
