# NOAA

A package for fetching and parsing NOAA data with Go (golang).

## Climate Data Online

cdo handles Climate Data Online (CDO) data from the National Center
for Environmental Information.  Use this interface to get historical
weather data (not predictions).  Note that this interface requires
you to get a free Web Service Token from NOAA.  As of this writing,
each token is limited to five requests per second, and 100 requests
per day.

The CDO tests (cdo_test.go) expect your token to be in the
NOAA_TOKEN environment variable.

More info at http://www.ncdc.noaa.gov/cdo-web/webservices/v2

## National Digital Forecast Database

ndfd handles National Digital Forecast Database (NDFD) data from the
National Weather Service (NWS) in Digital Weather Markup Language (DWML)
format.  Use this interface to get weather prediction data.  This API
does not require a NOAA token.

More info at http://graphical.weather.gov/xml/rest.php

## Installation

To install it, run:

    go get github.com/gershwinlabs/noaa

Canoncial usage can be seen in the cdo and ndfd tests.
