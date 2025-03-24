// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license.

package gospeech

import (
	"errors"
	"fmt"
	"strconv"
)

// SpeechConfig contains configuration for speech recognition services
type SpeechConfig struct {
	properties *PropertyCollection
}

// NewSpeechConfig creates a new empty speech configuration
func NewSpeechConfig() *SpeechConfig {
	config := &SpeechConfig{
		properties: NewPropertyCollection(),
	}
	config.properties.SetPropertyByName("GOSPEECH-SDK-PROGRAMMING-LANGUAGE", "Go")
	return config
}

// FromSubscription creates a speech config from subscription information
func FromSubscription(subscriptionKey, region string) (*SpeechConfig, error) {
	if subscriptionKey == "" {
		return nil, errors.New("subscription key cannot be empty")
	}
	if region == "" {
		return nil, errors.New("region cannot be empty")
	}

	config := NewSpeechConfig()
	config.properties.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	config.properties.SetProperty(SpeechServiceConnectionRegion, region)

	return config, nil
}

// FromEndpoint creates a speech config from an endpoint
func FromEndpoint(endpoint, subscriptionKey string) (*SpeechConfig, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint cannot be empty")
	}

	config := NewSpeechConfig()
	config.properties.SetProperty(SpeechServiceConnectionEndpoint, endpoint)

	if subscriptionKey != "" {
		config.properties.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	}

	return config, nil
}

// FromHost creates a speech config from a host address
func FromHost(host, subscriptionKey string) (*SpeechConfig, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}

	config := NewSpeechConfig()
	config.properties.SetProperty(SpeechServiceConnectionHost, host)

	if subscriptionKey != "" {
		config.properties.SetProperty(SpeechServiceConnectionKey, subscriptionKey)
	}

	return config, nil
}

// FromAuthorizationToken creates a speech config from an authorization token
func FromAuthorizationToken(authToken, region string) (*SpeechConfig, error) {
	if authToken == "" {
		return nil, errors.New("authorization token cannot be empty")
	}
	if region == "" {
		return nil, errors.New("region cannot be empty")
	}

	config := NewSpeechConfig()
	config.properties.SetProperty(SpeechServiceAuthorizationToken, authToken)
	config.properties.SetProperty(SpeechServiceConnectionRegion, region)

	return config, nil
}

// GetProperty gets a property value
func (c *SpeechConfig) GetProperty(propertyID PropertyID) string {
	return c.properties.GetProperty(propertyID)
}

// GetPropertyByName gets a property value by name
func (c *SpeechConfig) GetPropertyByName(propertyName string) string {
	return c.properties.GetPropertyByName(propertyName)
}

// SetProperty sets a property value
func (c *SpeechConfig) SetProperty(propertyID PropertyID, value string) {
	c.properties.SetProperty(propertyID, value)
}

// SetPropertyByName sets a property value by name
func (c *SpeechConfig) SetPropertyByName(propertyName string, value string) {
	c.properties.SetPropertyByName(propertyName, value)
}

// SetProperties sets multiple properties by ID
func (c *SpeechConfig) SetProperties(properties map[PropertyID]string) {
	for id, value := range properties {
		c.SetProperty(id, value)
	}
}

// SetPropertiesByName sets multiple properties by name
func (c *SpeechConfig) SetPropertiesByName(properties map[string]string) {
	for name, value := range properties {
		c.SetPropertyByName(name, value)
	}
}

// SetOutputFormat sets the output format
func (c *SpeechConfig) SetOutputFormat(format OutputFormat) {
	c.SetProperty(PropertyID("OutputFormat"), strconv.Itoa(int(format)))
}

// GetOutputFormat gets the output format
func (c *SpeechConfig) GetOutputFormat() OutputFormat {
	val := c.GetProperty(PropertyID("OutputFormat"))
	if val == "" {
		return OutputFormatSimple
	}

	format, err := strconv.Atoi(val)
	if err != nil {
		return OutputFormatSimple
	}

	return OutputFormat(format)
}

// SetSpeechRecognitionLanguage sets the speech recognition language
func (c *SpeechConfig) SetSpeechRecognitionLanguage(language string) {
	c.SetProperty(SpeechServiceConnectionRecoLanguage, language)
}

// GetSpeechRecognitionLanguage gets the speech recognition language
func (c *SpeechConfig) GetSpeechRecognitionLanguage() string {
	return c.GetProperty(SpeechServiceConnectionRecoLanguage)
}

// SetEndpointId sets the endpoint ID
func (c *SpeechConfig) SetEndpointID(endpointID string) {
	c.SetProperty(SpeechServiceConnectionEndpointID, endpointID)
}

// GetEndpointId gets the endpoint ID
func (c *SpeechConfig) GetEndpointID() string {
	return c.GetProperty(SpeechServiceConnectionEndpointID)
}

// SetAuthorizationToken sets the authorization token
func (c *SpeechConfig) SetAuthorizationToken(token string) {
	c.SetProperty(SpeechServiceAuthorizationToken, token)
}

// GetAuthorizationToken gets the authorization token
func (c *SpeechConfig) GetAuthorizationToken() string {
	return c.GetProperty(SpeechServiceAuthorizationToken)
}

// GetSubscriptionKey gets the subscription key
func (c *SpeechConfig) GetSubscriptionKey() string {
	return c.GetProperty(SpeechServiceConnectionKey)
}

// GetRegion gets the region
func (c *SpeechConfig) GetRegion() string {
	return c.GetProperty(SpeechServiceConnectionRegion)
}

// SetProxy sets proxy information
func (c *SpeechConfig) SetProxy(hostname string, port int, username, password string) error {
	if hostname == "" {
		return errors.New("hostname cannot be empty")
	}
	if port <= 0 {
		return errors.New("port must be greater than 0")
	}

	c.SetProperty(SpeechServiceConnectionProxyHostName, hostname)
	c.SetProperty(SpeechServiceConnectionProxyPort, strconv.Itoa(port))

	if username != "" {
		c.SetProperty(SpeechServiceConnectionProxyUserName, username)
	}

	if password != "" {
		c.SetProperty(SpeechServiceConnectionProxyPassword, password)
	}

	return nil
}

// SetServiceProperty sets a property that will be passed to the service
func (c *SpeechConfig) SetServiceProperty(name, value string, channel ServicePropertyChannel) {
	c.SetPropertyByName(fmt.Sprintf("ServiceProperty:%s:%d", name, channel), value)
}
