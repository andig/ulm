package tariff

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/tariff/entsoe"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
	"github.com/evcc-io/evcc/util/transport"
)

type Entsoe struct {
	*request.Helper
	*embed
	mux     sync.Mutex
	log     *util.Logger
	token   string
	domain  string
	data    api.Rates
	updated time.Time
}

var _ api.Tariff = (*Entsoe)(nil)

func init() {
	registry.Add("entsoe", NewEntsoeFromConfig)
}

func NewEntsoeFromConfig(other map[string]interface{}) (api.Tariff, error) {
	var cc struct {
		embed         `mapstructure:",squash"`
		Securitytoken string
		Domain        string
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.Securitytoken == "" {
		return nil, errors.New("securitytoken must be defined")
	}

	domain, err := entsoe.Area(entsoe.BZN, strings.ToUpper(cc.Domain))
	if err != nil {
		return nil, err
	}

	log := util.NewLogger("entsoe").Redact(cc.Securitytoken)

	t := &Entsoe{
		log:    log,
		Helper: request.NewHelper(log),
		embed:  &cc.embed,
		token:  cc.Securitytoken,
		domain: domain,
	}

	// Wrap the client with a decorator that adds the security token to each request.
	t.Client.Transport = &transport.Decorator{
		Base: t.Client.Transport,
		Decorator: transport.DecorateQuery(map[string]string{
			"securityToken": cc.Securitytoken,
		}),
	}

	done := make(chan error)
	go t.run(done)
	err = <-done

	return t, err
}

func (t *Entsoe) run(done chan error) {
	var once sync.Once

	bo := newBackoff()

	// Data updated by ESO every half hour, but we only need data every hour to stay current.
	for ; true; <-time.Tick(time.Hour) {
		var tr entsoe.PublicationMarketDocument

		if err := backoff.Retry(func() error {
			// Request the next 24 hours of data.
			data, err := t.DoBody(entsoe.DayAheadPricesRequest(t.domain, time.Hour*24))

			// Consider whether errors.As would be more appropriate if this needs to start dealing with wrapped errors.
			if se, ok := err.(request.StatusError); ok {
				if se.HasStatus(http.StatusBadRequest) {
					return backoff.Permanent(se)
				}

				return se
			}

			var doc entsoe.Document
			if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&doc); err != nil {
				return backoff.Permanent(err)
			}

			switch doc.XMLName.Local {
			case entsoe.AcknowledgementMarketDocumentName:
				var doc entsoe.AcknowledgementMarketDocument
				if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&doc); err != nil {
					return backoff.Permanent(err)
				}

				return backoff.Permanent(errors.New(doc.Reason.Text))

			case entsoe.PublicationMarketDocumentName:
				if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&tr); err != nil {
					return backoff.Permanent(err)
				}

				if tr.Type != string(entsoe.ProcessTypeDayAhead) {
					return backoff.Permanent(errors.New("invalid document type: " + tr.Type))
				}

				return nil

			default:
				return backoff.Permanent(errors.New("invalid document name: " + doc.XMLName.Local))
			}
		}, bo); err != nil {
			once.Do(func() { done <- err })

			t.log.ERROR.Println(err)
			continue
		}

		if len(tr.TimeSeries) == 0 {
			once.Do(func() { done <- entsoe.ErrInvalidData })
			t.log.ERROR.Println(entsoe.ErrInvalidData)
			continue
		}

		// extract desired series
		tsdata, err := entsoe.GetTsPriceData(tr.TimeSeries, entsoe.ResolutionHour)
		if err != nil {
			once.Do(func() { done <- err })
			t.log.ERROR.Println(err)
			continue
		}

		once.Do(func() { close(done) })

		t.mux.Lock()
		t.updated = time.Now()

		t.data = make(api.Rates, 0, len(tsdata))
		for _, r := range tsdata {
			ar := api.Rate{
				Start: r.ValidityStart,
				End:   r.ValidityEnd,
				Price: t.totalPrice(r.Value),
			}
			t.data = append(t.data, ar)
		}

		t.mux.Unlock()
	}
}

// Rates implements the api.Tariff interface
func (t *Entsoe) Rates() (api.Rates, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	return slices.Clone(t.data), outdatedError(t.updated, time.Hour)
}

// Type implements the api.Tariff interface
func (t *Entsoe) Type() api.TariffType {
	return api.TariffTypePriceDynamic
}
