# WISdoM OSS — External API Service
<p>
<img src="https://img.shields.io/github/go-mod/go-version/wisdom-oss/service-external-apis?style=for-the-badge" alt="Go Lang Version"/>
<a href="openapi.yaml">
<img src="https://img.shields.io/badge/Schema%20Version-3.0.0-6BA539?style=for-the-badge&logo=OpenAPI%20Initiative" alt="Open
API Schema Version"/></a>
</p>

## Overview
This microservice allows calling external APIs via the API gateway and therefore
without causing CORS issues. 

It also supports setting Authorization information on the request and additional
headers which might be used by the api. It also is authenticated like every 
other microservie in the WISdoM architecture.