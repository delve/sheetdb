package sheetdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/takuoki/gsheets"
)

var modelSets = map[string][]model{}

type model struct {
	name      string
	sheetName string
	loadFunc  func(data *gsheets.Sheet) error
}

// RegisterModel registers the model specified as an argument into a model set.
// This function is usually called from generated code.
func RegisterModel(modelSetName, modelName, sheetName string, loadFunc func(data *gsheets.Sheet) error) {
	m := model{
		name:      modelName,
		sheetName: sheetName,
		loadFunc:  loadFunc,
	}
	if s, ok := modelSets[modelSetName]; ok {
		modelSets[modelSetName] = append(s, m)
	} else {
		modelSets[modelSetName] = []model{m}
	}
}

// Client is a client of this package.
// Create a new client with the `New` function.
type Client struct {
	gsClient      *gsheets.Client
	spreadsheetID string
	modelSetName  string
}

// New creates and returns a new client.
func New(ctx context.Context, credentials, token, spreadsheetID string, opts ...ClientOption) (*Client, error) {
	gsClient, err := gsheets.New(ctx, credentials, token, gsheets.ClientWritable())
	if err != nil {
		return nil, fmt.Errorf("Unable to create gsheets client: %v", err)
	}
	client := &Client{
		gsClient:      gsClient,
		spreadsheetID: spreadsheetID,
		modelSetName:  "default",
	}
	for _, opt := range opts {
		client = opt(client)
	}
	return client, nil
}

// LoadData loads data from a spreadsheet into cache.
// This function calls the load functions of models registered in advance
// into the model set in order.
func (c *Client) LoadData(ctx context.Context) error {
	if c.gsClient == nil {
		return errors.New("The client has not been created correctly")
	}
	for _, m := range modelSets[c.modelSetName] {
		data, err := c.gsClient.GetSheet(ctx, c.spreadsheetID, m.sheetName)
		if err != nil {
			return err
		}
		logger.Infof("Loading from '%s' model", m.name)
		err = m.loadFunc(data)
		if err != nil {
			return fmt.Errorf("Unable to load '%s' data: %v", m.name, err)
		}
	}
	return nil
}

// AsyncUpdate applies updates to s spreadsheet asynchronously.
// This function is usually called from generated code.
func (c *Client) AsyncUpdate(data []gsheets.UpdateValue) error {
	if c.gsClient == nil {
		return errors.New("The client has not been created correctly")
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("Data could not be reflected on the sheet because an error occurred (err=%v, data=%+v)", e, data)
			}
		}()
		if err := c.gsClient.BatchUpdate(context.Background(), c.spreadsheetID, data...); err != nil {
			panic(fmt.Sprintf("Unable to update spreadsheet: %v", err))
		}
	}()
	return nil
}
