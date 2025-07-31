# ProfileApi

All URIs are relative to *http://localhost:4000/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**profileGet**](#profileget) | **GET** /profile | Get user profile|

# **profileGet**
> ProfileGet200Response profileGet()

Retrieves the authenticated user\'s profile information

### Example

```typescript
import {
    ProfileApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ProfileApi(configuration);

let authorization: string; //Bearer token (default to undefined)

const { status, data } = await apiInstance.profileGet(
    authorization
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **authorization** | [**string**] | Bearer token | defaults to undefined|


### Return type

**ProfileGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Profile retrieved successfully |  -  |
|**401** | Unauthorized (code: UNAUTHORIZED) |  -  |
|**404** | Profile not found (code: PROFILE_NOT_FOUND) |  -  |
|**500** | Internal server error (code: INTERNAL_ERROR) |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

