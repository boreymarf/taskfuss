# ServiceApi

All URIs are relative to *http://localhost:4000/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**pingGet**](#pingget) | **GET** /ping | Server health check|

# **pingGet**
> DtoPongResponse pingGet()

Returns \"pong\" if the server is running

### Example

```typescript
import {
    ServiceApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ServiceApi(configuration);

const { status, data } = await apiInstance.pingGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**DtoPongResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Server is running |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

